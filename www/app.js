const API_BASE = "https://orhalevanah-production.up.railway.app";

const MONTHS_ES = [
  "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
  "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"
];

const BIBLICAL_MONTHS = [
  "Aviv",
  "Segundo mes",
  "Tercer mes",
  "Cuarto mes",
  "Quinto mes",
  "Sexto mes",
  "Séptimo mes",
  "Octavo mes",
  "Noveno mes",
  "Décimo mes",
  "Undécimo mes",
  "Duodécimo mes"
];

const state = {
  currentView: "inicio",
  currentMonthDate: new Date(),
  selectedDate: null,
  todayPayload: null,
  todayFeastsPayload: null,
  monthPayload: null,
  selectedDatePayload: null
};

async function fetchJson(url) {
  const response = await fetch(url, { cache: "no-store" });
  const text = await response.text();

  if (!response.ok) {
    throw new Error(`HTTP ${response.status} - ${text}`);
  }

  try {
    return JSON.parse(text);
  } catch {
    throw new Error(`Respuesta no JSON: ${text}`);
  }
}

function safeText(value, fallback = "-") {
  if (value === undefined || value === null || value === "") return fallback;
  return String(value);
}

function pad(n) {
  return String(n).padStart(2, "0");
}

function formatDateToYMD(date) {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

function isSameYMD(a, b) {
  return a === b;
}

function parseYmdToLocalDate(ymd) {
  if (!ymd || typeof ymd !== "string") return null;
  const [year, month, day] = ymd.split("-").map(Number);
  if (!year || !month || !day) return null;
  return new Date(year, month - 1, day);
}

function addDays(date, days) {
  const copy = new Date(date.getFullYear(), date.getMonth(), date.getDate());
  copy.setDate(copy.getDate() + days);
  return copy;
}

function getBiblicalMonthName(monthNumber) {
  if (!monthNumber || monthNumber < 1 || monthNumber > 12) return null;
  return BIBLICAL_MONTHS[monthNumber - 1];
}

function inferBiblicalMonthNumber(data) {
  if (!data) return null;

  if (data.biblical_month) return data.biblical_month;

  const civilDate = parseYmdToLocalDate(data.civil_date);
  const monthStart = parseYmdToLocalDate(data.month_start);
  const nextMonthStart = parseYmdToLocalDate(data.next_month_start);

  if (!civilDate || !monthStart) return null;

  if (nextMonthStart && civilDate >= monthStart && civilDate < nextMonthStart) {
    return 1;
  }

  if (nextMonthStart && civilDate >= nextMonthStart) {
    return 2;
  }

  return 1;
}

function inferBiblicalDay(data) {
  if (!data) return null;

  if (data.biblical_day) return data.biblical_day;

  const civilDate = parseYmdToLocalDate(data.civil_date);
  const monthStart = parseYmdToLocalDate(data.month_start);
  if (!civilDate || !monthStart) return null;

  const diffMs = civilDate.getTime() - monthStart.getTime();
  const diffDays = Math.floor(diffMs / 86400000);

  return diffDays >= 0 ? diffDays + 1 : null;
}

function formatBiblicalDate(data) {
  if (!data) return "-";

  const day = inferBiblicalDay(data);
  const monthNumber = inferBiblicalMonthNumber(data);
  const monthName = data.month_name || getBiblicalMonthName(monthNumber);

  if (day && monthName) {
    return `${day} de ${monthName}`;
  }

  if (day) {
    return `Día ${day}`;
  }

  if (data.biblical_date) {
    return data.biblical_date;
  }

  return "-";
}

function formatFeastDate(item) {
  if (!item) return "-";
  return item.gregorian_date || item.gregorian_start_date || "-";
}

function setHeaderSubtitle(view) {
  const subtitle = document.getElementById("headerSubtitle");
  if (!subtitle) return;

  const map = {
    inicio: "Calendario Bíblico",
    hoy: "Hoy Bíblico",
    calendario: "Vista mensual",
    fiestas: "Fiestas Bíblicas",
    ajustes: "Configuración"
  };

  subtitle.textContent = map[view] || "Calendario Bíblico";
}

function setView(view) {
  state.currentView = view;

  document.querySelectorAll(".view").forEach((section) => {
    section.classList.remove("active");
  });

  document.querySelectorAll(".nav-btn").forEach((btn) => {
    btn.classList.remove("active");
  });

  const targetView = document.getElementById(`view-${view}`);
  if (targetView) targetView.classList.add("active");

  const activeBtn = document.querySelector(`.nav-btn[data-view="${view}"]`);
  if (activeBtn) activeBtn.classList.add("active");

  setHeaderSubtitle(view);
}

function renderFeastList(containerId, items, emptyMessage) {
  const container = document.getElementById(containerId);
  if (!container) return;

  if (!Array.isArray(items) || items.length === 0) {
    container.className = "list-stack empty-state";
    container.innerHTML = emptyMessage;
    return;
  }

  container.className = "list-stack";
  container.innerHTML = items.map((item) => `
    <div class="list-card">
      <div class="list-card-title">${safeText(item.name, "Evento")}</div>
      ${
        item.gregorian_date || item.gregorian_start_date
          ? `
            <div class="list-card-meta">
              <span>Fecha: ${formatFeastDate(item)}</span>
            </div>
          `
          : ""
      }
      ${
        item.description
          ? `<div class="list-card-text">${item.description}</div>`
          : ""
      }
    </div>
  `).join("");
}

async function getTodayBiblical() {
  return fetchJson(`${API_BASE}/api/biblical/jerusalem/today`);
}

async function getMonthBiblical() {
  return fetchJson(`${API_BASE}/api/biblical/jerusalem/month`);
}

async function getFeastsJerusalem() {
  return fetchJson(`${API_BASE}/api/feasts/jerusalem`);
}

async function getBiblicalByDate(date) {
  return fetchJson(`${API_BASE}/api/biblical/jerusalem/date?date=${encodeURIComponent(date)}`);
}

function getAfterSunsetText(value) {
  if (value === true) return "Sí, después del atardecer";
  if (value === false) return "No, antes del atardecer";
  return "-";
}

function fillMonthSummaryUI(monthData) {
  const biblicalText = formatBiblicalDate(monthData);

  const bindings = {
    inicioBiblicalDayNumber: biblicalText,
    inicioMonthStart: monthData?.month_start,
    inicioNextMonthStart: monthData?.next_month_start,
    inicioAfterSunset: getAfterSunsetText(monthData?.after_sunset),
    hoyBiblicalDayNumber: biblicalText,
    hoyMonthStart: monthData?.month_start,
    hoyNextMonthStart: monthData?.next_month_start,
    hoyAfterSunset: getAfterSunsetText(monthData?.after_sunset)
  };

  Object.entries(bindings).forEach(([id, value]) => {
    const el = document.getElementById(id);
    if (el) el.textContent = safeText(value, "-");
  });

  const monthStatus = document.getElementById("monthStatusBadge");
  if (monthStatus) {
    monthStatus.textContent = biblicalText !== "-" ? biblicalText : "Calendario bíblico";
    monthStatus.className = "pill";
  }

  const nextStatus = document.getElementById("nextMonthStatus");
  if (nextStatus) {
    nextStatus.textContent = monthData?.next_month_start
      ? `Próxima cabeza de mes estimada: ${monthData.next_month_start}`
      : "Próxima cabeza de mes no disponible";
  }
}

function fillTodayUI(todayData, monthData) {
  const merged = { ...(monthData || {}), ...(todayData || {}) };

  const inicioCivilDate = document.getElementById("inicioCivilDate");
  const inicioBiblicalDate = document.getElementById("inicioBiblicalDate");
  const inicioJerusalemTime = document.getElementById("inicioJerusalemTime");
  const inicioSunset = document.getElementById("inicioSunset");
  const inicioDayNote = document.getElementById("inicioDayNote");

  const hoyCivilDate = document.getElementById("hoyCivilDate");
  const hoyBiblicalDate = document.getElementById("hoyBiblicalDate");
  const hoyJerusalemTime = document.getElementById("hoyJerusalemTime");
  const hoySunsetTime = document.getElementById("hoySunsetTime");
  const hoyDayNote = document.getElementById("hoyDayNote");

  if (inicioCivilDate) inicioCivilDate.textContent = safeText(merged?.civil_date, "Sin fecha");
  if (inicioBiblicalDate) inicioBiblicalDate.textContent = formatBiblicalDate(merged);
  if (inicioJerusalemTime) inicioJerusalemTime.textContent = safeText(merged?.jerusalem_time, "No disponible");
  if (inicioSunset) inicioSunset.textContent = safeText(merged?.sunset_time, "No disponible");
  if (inicioDayNote) {
    inicioDayNote.textContent = safeText(
      merged?.day_note,
      "Información cargada correctamente."
    );
  }

  if (hoyCivilDate) hoyCivilDate.textContent = safeText(merged?.civil_date, "Sin fecha");
  if (hoyBiblicalDate) hoyBiblicalDate.textContent = formatBiblicalDate(merged);
  if (hoyJerusalemTime) hoyJerusalemTime.textContent = safeText(merged?.jerusalem_time, "No disponible");
  if (hoySunsetTime) hoySunsetTime.textContent = safeText(merged?.sunset_time, "No disponible");
  if (hoyDayNote) {
    hoyDayNote.textContent = safeText(
      merged?.day_note,
      "Sin nota disponible."
    );
  }

  fillMonthSummaryUI(merged);
}

function fillSelectedExtraUI(data) {
  const selectedBiblicalDayNumber = document.getElementById("selectedBiblicalDayNumber");
  const selectedMonthStart = document.getElementById("selectedMonthStart");
  const selectedNextMonthStart = document.getElementById("selectedNextMonthStart");
  const selectedAfterSunset = document.getElementById("selectedAfterSunset");

  if (selectedBiblicalDayNumber) {
    selectedBiblicalDayNumber.textContent = formatBiblicalDate(data);
  }
  if (selectedMonthStart) {
    selectedMonthStart.textContent = safeText(data?.month_start, "-");
  }
  if (selectedNextMonthStart) {
    selectedNextMonthStart.textContent = safeText(data?.next_month_start, "-");
  }
  if (selectedAfterSunset) {
    selectedAfterSunset.textContent = getAfterSunsetText(data?.after_sunset);
  }
}

function isShabat(dateObj) {
  return dateObj.getDay() === 6;
}

function normalizeFeastPool() {
  const fromToday = state.todayFeastsPayload || {};

  return [
    ...(Array.isArray(fromToday.current_feasts) ? fromToday.current_feasts : []),
    ...(Array.isArray(fromToday.upcoming_feasts) ? fromToday.upcoming_feasts : [])
  ];
}

function getFeastsForDate(ymd) {
  const pool = normalizeFeastPool();

  return pool.filter((item) => {
    const startYmd = item?.gregorian_start_date || item?.gregorian_date;
    const endYmd = item?.gregorian_end_date || item?.gregorian_date;

    if (startYmd && endYmd) {
      return ymd >= startYmd && ymd <= endYmd;
    }

    return startYmd === ymd || endYmd === ymd;
  });
}

function buildSelectedDayEvents(dateString) {
  const events = [];
  const selectedDateObj = parseYmdToLocalDate(dateString);

  if (selectedDateObj && isShabat(selectedDateObj)) {
    events.push({
      name: "Shabat",
      description: "Este día corresponde a Shabat."
    });
  }

  const sameDayFeasts = getFeastsForDate(dateString);

  sameDayFeasts.forEach((item) => {
    const alreadyExists = events.some((e) => e.name === item.name);
    if (!alreadyExists) {
      events.push({
        name: item.name,
        gregorian_date: item.gregorian_date || item.gregorian_start_date || dateString,
        description: item.description || `Este día corresponde a ${item.name}.`
      });
    }
  });

  return events;
}

function buildCalendarChip(label, className = "") {
  return `<span class="calendar-chip ${className}">${label}</span>`;
}

function renderCalendar() {
  const grid = document.getElementById("calendarGrid");
  const monthLabel = document.getElementById("calendarMonthLabel");
  const monthSubLabel = document.getElementById("calendarMonthSubLabel");

  if (!grid || !monthLabel || !monthSubLabel) return;

  const current = state.currentMonthDate;
  const year = current.getFullYear();
  const month = current.getMonth();

  monthLabel.textContent = `${MONTHS_ES[month]} ${year}`;
  monthSubLabel.textContent = "Calendario civil con referencia bíblica";

  grid.innerHTML = "";

  const firstDay = new Date(year, month, 1).getDay();
  const lastDate = new Date(year, month + 1, 0).getDate();
  const todayString = state.todayPayload?.civil_date || state.monthPayload?.civil_date || formatDateToYMD(new Date());
  const monthStartYmd = state.monthPayload?.month_start || null;
  const nextMonthStartYmd = state.monthPayload?.next_month_start || null;

  for (let i = 0; i < firstDay; i++) {
    const empty = document.createElement("div");
    empty.className = "calendar-cell empty";
    grid.appendChild(empty);
  }

  for (let day = 1; day <= lastDate; day++) {
    const dateObj = new Date(year, month, day);
    const ymd = formatDateToYMD(dateObj);

    const cell = document.createElement("button");
    cell.type = "button";
    cell.className = "calendar-cell";

    if (isSameYMD(ymd, todayString)) {
      cell.classList.add("today");
    }

    if (state.selectedDate && isSameYMD(ymd, state.selectedDate)) {
      cell.classList.add("selected");
    }

    if (monthStartYmd && isSameYMD(ymd, monthStartYmd)) {
      cell.classList.add("month-start");
    }

    if (nextMonthStartYmd && isSameYMD(ymd, nextMonthStartYmd)) {
      cell.classList.add("next-month-start");
    }

    if (isShabat(dateObj)) {
      cell.classList.add("shabat");
    }

    const dayFeasts = getFeastsForDate(ymd);
    if (dayFeasts.length > 0) {
      cell.classList.add("has-feast");
    }

    const chips = [];

    if (monthStartYmd && isSameYMD(ymd, monthStartYmd)) {
      chips.push(buildCalendarChip("Jodesh", "chip-jodesh"));
    }

    if (nextMonthStartYmd && isSameYMD(ymd, nextMonthStartYmd)) {
      chips.push(buildCalendarChip("Próx.", "chip-next"));
    }

    if (isSameYMD(ymd, todayString)) {
      chips.push(buildCalendarChip("Hoy", "chip-today"));
    }

    if (isShabat(dateObj)) {
      chips.push(buildCalendarChip("Shabat", "chip-shabat"));
    }

    dayFeasts.slice(0, 2).forEach((feast) => {
      chips.push(buildCalendarChip(feast.name, "chip-feast"));
    });

    cell.innerHTML = `
      <span class="calendar-day-number">${day}</span>
      <span class="calendar-day-label">${ymd}</span>
      <span class="calendar-chip-stack">${chips.join("")}</span>
    `;

    cell.addEventListener("click", async () => {
      state.selectedDate = ymd;
      renderCalendar();
      await loadSelectedDate(ymd);
    });

    grid.appendChild(cell);
  }
}

async function loadTodayData() {
  try {
    const [todayPayload, monthPayload, feastsPayload] = await Promise.all([
      getTodayBiblical(),
      getMonthBiblical(),
      getFeastsJerusalem().catch(() => null)
    ]);

    const todayData = todayPayload?.data || {};
    const monthData = monthPayload?.data || {};

    state.todayPayload = todayData;
    state.monthPayload = monthData;

    fillTodayUI(todayData, monthData);

    const todayDate = todayData.civil_date || monthData.civil_date || formatDateToYMD(new Date());

    let feastData = feastsPayload?.data || null;

    if (!feastData || (!Array.isArray(feastData.current_feasts) && !Array.isArray(feastData.upcoming_feasts))) {
      feastData = {
        current_feasts: [],
        upcoming_feasts: []
      };
    }

    state.todayFeastsPayload = feastData;

    renderFeastList(
      "hoyCurrentFeasts",
      feastData.current_feasts || [],
      "No hay fiestas activas para este día."
    );

    renderFeastList(
      "hoyUpcomingFeasts",
      feastData.upcoming_feasts || [],
      "No se encontraron próximas fiestas."
    );

    renderFeastList(
      "fiestasCurrentList",
      feastData.current_feasts || [],
      "No hay fiestas activas para este día."
    );

    renderFeastList(
      "fiestasUpcomingList",
      feastData.upcoming_feasts || [],
      "No se encontraron próximas fiestas."
    );

    state.selectedDate = todayDate;
    renderCalendar();
  } catch (error) {
    console.error("Error cargando hoy:", error);

    [
      "inicioCivilDate", "inicioBiblicalDate", "inicioJerusalemTime", "inicioSunset",
      "hoyCivilDate", "hoyBiblicalDate", "hoyJerusalemTime", "hoySunsetTime"
    ].forEach((id) => {
      const el = document.getElementById(id);
      if (el) el.textContent = "Error";
    });

    const inicioDayNote = document.getElementById("inicioDayNote");
    const hoyDayNote = document.getElementById("hoyDayNote");
    if (inicioDayNote) inicioDayNote.textContent = error.message || "Failed to fetch";
    if (hoyDayNote) hoyDayNote.textContent = error.message || "Failed to fetch";

    renderFeastList("hoyCurrentFeasts", [], "No fue posible cargar las fiestas actuales.");
    renderFeastList("hoyUpcomingFeasts", [], "No fue posible cargar las próximas fiestas.");
    renderFeastList("fiestasCurrentList", [], "No fue posible cargar las fiestas actuales.");
    renderFeastList("fiestasUpcomingList", [], "No fue posible cargar las próximas fiestas.");
  }
}

async function loadSelectedDate(dateString) {
  const selectedCivilDate = document.getElementById("selectedCivilDate");
  const selectedBiblicalDate = document.getElementById("selectedBiblicalDate");
  const selectedJerusalemTime = document.getElementById("selectedJerusalemTime");
  const selectedSunsetTime = document.getElementById("selectedSunsetTime");
  const selectedDayNote = document.getElementById("selectedDayNote");

  if (selectedCivilDate) selectedCivilDate.textContent = dateString;
  if (selectedBiblicalDate) selectedBiblicalDate.textContent = "Cargando...";
  if (selectedJerusalemTime) selectedJerusalemTime.textContent = "Cargando...";
  if (selectedSunsetTime) selectedSunsetTime.textContent = "Cargando...";
  if (selectedDayNote) selectedDayNote.textContent = "Consultando detalle del día...";

  fillSelectedExtraUI(null);
  renderFeastList("selectedCurrentFeasts", [], "Consultando si es fiesta o Shabat...");

  try {
    const payload = await getBiblicalByDate(dateString);
    const data = payload?.data || {};

    state.selectedDatePayload = data;

    if (selectedCivilDate) selectedCivilDate.textContent = safeText(data.civil_date, dateString);
    if (selectedBiblicalDate) selectedBiblicalDate.textContent = formatBiblicalDate(data);
    if (selectedJerusalemTime) selectedJerusalemTime.textContent = safeText(data.jerusalem_time, "No disponible");
    if (selectedSunsetTime) selectedSunsetTime.textContent = safeText(data.sunset_time, "No disponible");
    if (selectedDayNote) {
      selectedDayNote.textContent = safeText(
        data.day_note,
        "El día bíblico comienza al atardecer en Jerusalén y termina al siguiente atardecer."
      );
    }

    fillSelectedExtraUI(data);

    const dayEvents = buildSelectedDayEvents(dateString);

    renderFeastList(
      "selectedCurrentFeasts",
      dayEvents,
      "No hay fiesta ni Shabat para esta fecha."
    );

    renderCalendar();
  } catch (error) {
    console.error("Error cargando fecha seleccionada:", error);

    if (selectedBiblicalDate) selectedBiblicalDate.textContent = "No disponible";
    if (selectedJerusalemTime) selectedJerusalemTime.textContent = "No disponible";
    if (selectedSunsetTime) selectedSunsetTime.textContent = "No disponible";
    if (selectedDayNote) {
      selectedDayNote.textContent = `Error al consultar la fecha: ${error.message || "No se pudo cargar el detalle."}`;
    }

    fillSelectedExtraUI(null);

    const fallbackEvents = buildSelectedDayEvents(dateString);
    renderFeastList(
      "selectedCurrentFeasts",
      fallbackEvents,
      "No fue posible cargar los eventos de esta fecha."
    );
  }
}

function bindEvents() {
  document.querySelectorAll(".nav-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const view = btn.getAttribute("data-view");
      if (view) setView(view);
    });
  });

  document.querySelectorAll("[data-go-view]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const view = btn.getAttribute("data-go-view");
      if (view) setView(view);
    });
  });

  const prevMonthBtn = document.getElementById("prevMonthBtn");
  const nextMonthBtn = document.getElementById("nextMonthBtn");
  const refreshDataBtn = document.getElementById("refreshDataBtn");
  const goTodayBtn = document.getElementById("goTodayBtn");

  prevMonthBtn?.addEventListener("click", () => {
    state.currentMonthDate = new Date(
      state.currentMonthDate.getFullYear(),
      state.currentMonthDate.getMonth() - 1,
      1
    );
    renderCalendar();
  });

  nextMonthBtn?.addEventListener("click", () => {
    state.currentMonthDate = new Date(
      state.currentMonthDate.getFullYear(),
      state.currentMonthDate.getMonth() + 1,
      1
    );
    renderCalendar();
  });

  refreshDataBtn?.addEventListener("click", async () => {
    await loadTodayData();
    renderCalendar();
  });

  goTodayBtn?.addEventListener("click", async () => {
    const now = new Date();
    state.currentMonthDate = new Date(now.getFullYear(), now.getMonth(), 1);
    state.selectedDate = state.todayPayload?.civil_date || state.monthPayload?.civil_date || formatDateToYMD(now);
    renderCalendar();
    await loadSelectedDate(state.selectedDate);
    setView("calendario");
  });
}

document.addEventListener("DOMContentLoaded", async () => {
  bindEvents();
  setView("inicio");
  renderCalendar();
  await loadTodayData();

  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("./sw.js").catch((error) => {
      console.warn("No se pudo registrar el Service Worker:", error);
    });
  }
});

window.setView = setView;