/* Shared by the web portal and Capacitor; copied by scripts/sync-day-panels.cjs. */
globalThis.DayPanels = (() => {
  const pending = new Map();
  async function json(url) {
    if (!pending.has(url)) {
      pending.set(url, fetch(url, { cache: "no-store" }).then(async response => {
        if (!response.ok) throw new Error(`Error ${response.status} al consultar información`);
        const payload = await response.json();
        if (!payload.success) throw new Error("Respuesta inválida del servidor");
        return payload;
      }).finally(() => pending.delete(url)));
    }
    return pending.get(url);
  }
  function node(tag, text) {
    const el = document.createElement(tag);
    el.textContent = text;
    return el;
  }
  function dateLabel(date) {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(date || "")) return "Fecha no disponible";
    return new Intl.DateTimeFormat("es", {
      weekday: "long", day: "numeric", month: "long", year: "numeric", timeZone: "UTC"
    }).format(new Date(`${date}T12:00:00Z`));
  }
  function feastProgress(feast, date) {
    const start = feast.gregorian_start_date || feast.gregorian_date;
    const end = feast.gregorian_end_date;
    if (![start, end, date].every(value => /^\d{4}-\d{2}-\d{2}$/.test(value || ""))) return null;
    const dayNumber = value => Date.parse(`${value}T00:00:00Z`) / 86400000;
    const total = dayNumber(end) - dayNumber(start);
    const day = dayNumber(date) - dayNumber(start) + 1;
    return total > 1 && day >= 1 && day <= total ? { day, total, label: `Día ${day} de ${total}` } : null;
  }
  function dayEvents(el, date, feasts, omer, shabbat) {
    if (!el) return;
    el.hidden = false;
    el.replaceChildren(node("h3", dateLabel(date)));
    for (const feast of feasts) {
      const progress = feastProgress(feast, date);
      el.append(node("p", `${feast.name}${progress ? ` · ${progress.label}` : ""}`));
    }
    if (omer) el.append(node("p", `Día ${omer} del Omer`));
    if (shabbat) el.append(node("p", shabbat));
    if (!feasts.length && !omer && !shabbat) el.append(node("p", "No hay fiestas registradas para esta fecha."));
    el.append(node("small", "El calendario indica el día bíblico que comienza al atardecer."));
  }
  function today(el, data, feasts = {}, localDate = data?.civil_date) {
    if (!el) return;
    // Requests straddling Jerusalem midnight must not mix two civil dates.
    if (feasts?.civil_date && feasts.civil_date !== data.civil_date) feasts = null;
    el.replaceChildren(node("strong", "HOY · FECHA LOCAL"), node("h2", dateLabel(localDate)));
    const indicators = node("div", "");
    indicators.className = "day-indicators";
    const labels = [
      `Día bíblico ${data.biblical_day ?? "—"}`,
      `Atardecer ${data.sunset_time || "—"}`,
      data.after_sunset ? "Después del atardecer" : "Antes del atardecer"
    ];
    const phaseLabels = {
      "New Moon": "Luna nueva", "Waxing Crescent": "Creciente",
      "First Quarter": "Cuarto creciente", "Waxing Gibbous": "Gibosa creciente",
      "Full Moon": "Luna llena", "Waning Gibbous": "Gibosa menguante",
      "Last Quarter": "Cuarto menguante", "Waning Crescent": "Menguante"
    };
    if (data.moon?.phase_name) labels.push(phaseLabels[data.moon.phase_name] || data.moon.phase_name);
    if (feasts?.omer?.is_omer_counting) labels.push(`Omer: ${feasts.omer.omer_day_today}`);
    for (const label of labels) indicators.append(node("span", label));
    el.append(indicators);
    let feastDate = data.civil_date;
    if (!data.after_sunset && /^\d{4}-\d{2}-\d{2}$/.test(feastDate || "")) {
      const previous = new Date(`${feastDate}T00:00:00Z`);
      previous.setUTCDate(previous.getUTCDate() - 1);
      feastDate = previous.toISOString().slice(0, 10);
    }
    const names = (feasts?.current_feasts || []).map(feast => {
      const progress = feastProgress(feast, feastDate);
      return `${feast.name}${progress ? ` · ${progress.label}` : ""}`;
    });
    el.append(node("p", !feasts ? "No se pudieron consultar las fiestas de hoy." : names.length ? names.join(" · ") : "No hay fiestas registradas para hoy."));
    const localTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "zona horaria del dispositivo";
    el.append(node("small", `Fecha local · ${localTimezone} · Calendario y hora de referencia: Jerusalén · ${data.jerusalem_time || ""}`));
  }
  async function sources(el, date, base = "", force = false) {
    if (!el) return;
    // A new selection invalidates the previous request, including errors.
    const token = {};
    el.requestToken = token;
    el.hidden = true;
    el.replaceChildren();
    if (!/^\d{4}-\d{2}-\d{2}$/.test(date || "")) return;
    if (force) {
      el.hidden = false;
      el.append(node("p", "Consultando fuentes…"));
    }
    try {
      const payload = await json(`${base}/api/biblical/jerusalem/sources?date=${encodeURIComponent(date)}`);
      if (el.requestToken !== token) return;
      if (payload.data?.date !== date) throw new Error("La respuesta no corresponde a la fecha seleccionada");
      if (!payload.data.is_month_start) { el.hidden = true; return; }
      el.hidden = false;
      el.replaceChildren(node("h3", `Cabeza de mes · Fuentes · ${date}`));
      const list = payload.data.sources || [];
      const incomplete = ["pending", "partial"].includes(payload.data.search_state);
      if (!list.length) el.append(node("p", incomplete ? "La búsqueda de fuentes está pendiente o incompleta; todavía no se puede confirmar que no existan fuentes." : "No existen fuentes registradas para este período"));
      else if (incomplete) el.append(node("p", "Fuentes disponibles; la actualización externa está pendiente o incompleta."));
      for (const source of list) {
        const item = node("p", "");
        let link;
        try {
          const url = new URL(source.url);
          if (["http:", "https:"].includes(url.protocol) && !url.username && !url.password) {
            link = node("a", source.name || "Consultar fuente");
            link.href = url.href;
            link.target = "_blank";
            link.rel = "noopener noreferrer";
          }
        } catch { /* Invalid links are rendered as plain text. */ }
        item.append(link || node("strong", source.name || "Fuente registrada"));
        const labels = {
          confirmed: "Confirmada · Avistamiento reportado",
          declared: "Declaración de inicio · Sin avistamiento",
          related: "Reporte / noticia relacionada · No confirma esta fecha",
          estimated: "Pronóstico, pendiente de confirmación",
          pending: "Pendiente de confirmación"
        };
        item.append(node("span", ` · ${labels[source.status] || "Mención indirecta · Pendiente de confirmación"}`));
        if (source.title) item.append(node("span", ` · ${source.title}`));
        if (source.date) item.append(node("small", ` · Fecha de la evidencia: ${source.date}`));
        el.append(item);
      }
      el.scrollIntoView?.({ block: "nearest" });
    } catch (error) {
      if (el.requestToken !== token) return;
      el.hidden = false;
      el.replaceChildren(node("h3", `Fuentes · ${date}`), node("p", "No se pudieron consultar las fuentes. Vuelve a seleccionar la fecha para reintentar."));
      el.scrollIntoView?.({ block: "nearest" });
      console.error("Consulta de fuentes:", error);
    }
  }
  return { json, dateLabel, feastProgress, dayEvents, today, sources };
})();
