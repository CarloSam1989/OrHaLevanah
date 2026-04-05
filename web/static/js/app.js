async function fetchJson(url) {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`Error ${response.status}: no se pudo cargar ${url}`);
  }
  return response.json();
}

const calendarState = {
  currentDate: new Date(),
  selectedDate: null,
  feastData: null,
  biblicalMonthInfo: null,
  jerusalemClockInterval: null,
};

const BIBLICAL_VERSE_LIBRARY = {
  shabbat: [
    {
      reference: "Éxodo 20:8-11",
      text: "Acuérdate del día de reposo para santificarlo...",
      note: "Mandamiento del Shabat."
    },
    {
      reference: "Isaías 58:13-14",
      text: "Si retrajeres del día de reposo tu pie...",
      note: "La bendición ligada al deleite en el día apartado."
    }
  ],
  months: {
    1: [
      {
        reference: "Éxodo 12:2",
        text: "Este mes os será principio de los meses...",
        note: "Aviv como comienzo del año bíblico."
      },
      {
        reference: "Éxodo 13:4",
        text: "Vosotros salís hoy en el mes de Abib.",
        note: "Identificación del primer mes."
      },
      {
        reference: "Deuteronomio 16:1",
        text: "Guardarás el mes de Abib, y harás pascua a YHWH tu Elohim.",
        note: "Relación entre Aviv y Pesaj."
      }
    ],
    7: [
      {
        reference: "Levítico 23:24",
        text: "En el mes séptimo, al primero del mes, tendréis día de reposo...",
        note: "Mes de convocaciones solemnes."
      }
    ]
  },
  feasts: {
    "Pesaj": [
      {
        reference: "Éxodo 12:5-14",
        text: "Y tomarán de la sangre, y la pondrán en los dos postes...",
        note: "Institución de Pesaj."
      },
      {
        reference: "1 Corintios 5:7",
        text: "Porque nuestra pascua, que es el Mesías, ya fue sacrificada por nosotros.",
        note: "Relación mesiánica de Pesaj."
      }
    ],
    "HaMatzot": [
      {
        reference: "Éxodo 12:15",
        text: "Siete días comeréis panes sin levadura...",
        note: "Mandato de los panes sin levadura."
      }
    ],
    "Bikurim": [
      {
        reference: "Levítico 23:10-11",
        text: "Traeréis al sacerdote una gavilla por primicia...",
        note: "Primicias delante de YHWH."
      }
    ],
    "Shavuot": [
      {
        reference: "Levítico 23:15-16",
        text: "Y contaréis desde el día que sigue al día de reposo...",
        note: "Conteo hacia Shavuot."
      }
    ],
    "Yom Teruah": [
      {
        reference: "Levítico 23:24",
        text: "Tendréis una conmemoración al son de trompetas...",
        note: "Memorial de aclamación."
      }
    ],
    "Yom Kippur": [
      {
        reference: "Levítico 23:27",
        text: "A los diez días de este mes séptimo será el día de expiación...",
        note: "Día de aflicción y expiación."
      }
    ],
    "Sukkot": [
      {
        reference: "Levítico 23:34",
        text: "A los quince días de este mes séptimo será la fiesta solemne de los tabernáculos...",
        note: "Inicio de Sukkot."
      }
    ],
    "Shemini Atzeret": [
      {
        reference: "Levítico 23:36",
        text: "El octavo día tendréis santa convocación...",
        note: "Octavo día de asamblea."
      }
    ]
  },
  omer: {
    general: [
      {
        reference: "Levítico 23:15-16",
        text: "Y contaréis desde el día que sigue al día de reposo, desde el día en que ofrecisteis la gavilla de la ofrenda mecida; siete semanas cumplidas serán.",
        note: "Base bíblica del conteo del Omer."
      },
      {
        reference: "Deuteronomio 16:9",
        text: "Siete semanas contarás; desde que comenzare a meterse la hoz en las mieses comenzarás a contar las siete semanas.",
        note: "Confirmación del conteo hacia Shavuot."
      },
      {
        reference: "Salmo 90:12",
        text: "Enséñanos de tal modo a contar nuestros días, que traigamos al corazón sabiduría.",
        note: "Aplicación espiritual del conteo."
      }
    ],
    day1: [
      {
        reference: "Levítico 23:10-11",
        text: "Cuando hayáis entrado en la tierra que yo os doy, y seguéis su mies, traeréis al sacerdote una gavilla por primicia de los primeros frutos de vuestra siega.",
        note: "Inicio del conteo desde Bikurim."
      },
      {
        reference: "Levítico 23:14",
        text: "No comeréis pan, ni grano tostado, ni espiga fresca, hasta este mismo día, hasta que hayáis ofrecido la ofrenda de vuestro Elohim.",
        note: "Relación entre primicias y consagración."
      }
    ]
  },
  general: [
    {
      reference: "Salmo 119:105",
      text: "Lámpara es a mis pies tu palabra, y lumbrera a mi camino.",
      note: "Dirección diaria por medio de la Palabra."
    },
    {
      reference: "Génesis 1:14",
      text: "Haya lumbreras en la expansión de los cielos para separar el día de la noche; y sirvan de señales para las estaciones, para días y años.",
      note: "Base bíblica de los tiempos señalados."
    },
    {
      reference: "Eclesiastés 3:1",
      text: "Todo tiene su tiempo, y todo lo que se quiere debajo del cielo tiene su hora.",
      note: "El Eterno estableció tiempos para cada propósito."
    },
    {
      reference: "Salmo 90:12",
      text: "Enséñanos de tal modo a contar nuestros días, que traigamos al corazón sabiduría.",
      note: "Sabiduría para vivir conforme al tiempo del Eterno."
    }
  ]
}; 

function pickRandomItems(items = [], count = 1) {
  if (!Array.isArray(items) || items.length === 0) return [];

  const shuffled = [...items].sort(() => Math.random() - 0.5);
  return shuffled.slice(0, Math.max(1, count));
}

function getRandomContextVerses(data) {
  const pool = [...(BIBLICAL_VERSE_LIBRARY.general || [])];

  if (Number(data?.omer_day || 0) > 0) {
    pool.push(...(BIBLICAL_VERSE_LIBRARY.shabbat || []));
  }

  if (Array.isArray(data?.current_feasts) && data.current_feasts.length > 0) {
    pool.push(...(BIBLICAL_VERSE_LIBRARY.general || []));
  }

  return pickRandomItems(pool, 1);
}

function getOmerDayForDate(dateKey) {
  const omer = calendarState.feastData?.omer;
  if (!omer?.bikkurim_date || !dateKey) return null;

  const biblicalOmerStart = adjustToBiblicalStart(omer.bikkurim_date);
  if (!biblicalOmerStart) return null;

  const start = normalizeDateInput(biblicalOmerStart);
  const current = normalizeDateInput(dateKey);

  const diffMs = current.getTime() - start.getTime();
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffDays >= 0 && diffDays < 49) {
    return diffDays + 1;
  }

  return null;
}

function getFeastsForDate(dateKey) {
  if (!calendarState.feastData || !dateKey) return [];

  const feastsMap = groupFeastsByVisibleDay(calendarState.feastData);
  return feastsMap.get(dateKey) || [];
}

function buildVerseContextForDate(dateKey) {
  const monthInfo = calendarState.biblicalMonthInfo || {};
  const todayData = calendarState.feastData || {};
  const feasts = getFeastsForDate(dateKey);
  const omerDay = getOmerDayForDate(dateKey);

  const dateObj = normalizeDateInput(dateKey);
  const dayOfWeek = dateObj.getDay();
  const isFriday = dayOfWeek === 5;
  const isSaturday = dayOfWeek === 6;

  const monthStartType = getMonthStartType(dateKey);
  const isMonthStart = monthStartType !== null;

  let inferredBiblicalMonth = Number(todayData?.biblical_month || 0);
  let inferredBiblicalDay = 0;

  if (isMonthStart) {
    if (monthStartType === "current") {
      inferredBiblicalDay = 1;
    } else if (monthStartType === "next") {
      inferredBiblicalDay = 1;
      inferredBiblicalMonth = inferredBiblicalMonth > 0 ? inferredBiblicalMonth + 1 : 0;
    }
  }

  if (feasts.length > 0) {
    inferredBiblicalMonth = Number(feasts[0]?.biblical_month || inferredBiblicalMonth || 0);
    inferredBiblicalDay = Number(feasts[0]?.biblical_day || inferredBiblicalDay || 0);
  }

  const biblicalDateLabel =
    inferredBiblicalMonth > 0 && inferredBiblicalDay > 0
      ? `${inferredBiblicalDay} de ${getBiblicalMonthName(inferredBiblicalMonth)}`
      : "Selección del calendario";

  return {
    civil_date: dateKey,
    biblical_date: biblicalDateLabel,
    biblical_month: inferredBiblicalMonth,
    biblical_day: inferredBiblicalDay,
    after_sunset: Boolean(monthInfo?.after_sunset),
    current_feasts: feasts,
    is_shabbat: isSaturday || (isFriday && Boolean(monthInfo?.after_sunset)),
    is_month_start: isMonthStart,
    month_start_type: monthStartType,
    omer_day: omerDay
  };
}

function monthTitle(date) {
  return date.toLocaleDateString("es-EC", {
    month: "long",
    year: "numeric",
  });
}

function escapeHtml(value) {
  return String(value ?? "")
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#39;");
}

function renderVerseItems(items = []) {
  if (!items.length) {
    return `<div class="verse-empty">No hay versículos asociados para esta sección.</div>`;
  }

  return `
    <div class="verse-list">
      ${items.map(item => `
        <article class="verse-item">
          <span class="verse-ref">${escapeHtml(item.reference)}</span>
          <p class="verse-text">${escapeHtml(item.text)}</p>
          ${item.note ? `<div class="verse-note">${escapeHtml(item.note)}</div>` : ""}
        </article>
      `).join("")}
    </div>
  `;
}

function buildBiblicalVersesHtml(data) {
  const month = Number(data?.biblical_month || 0);
  const currentFeasts = Array.isArray(data?.current_feasts) ? data.current_feasts : [];
  const afterSunset = Boolean(data?.after_sunset);
  const omerDay = Number(data?.omer_day || 0);

  const isShabbat =
    Boolean(data?.is_shabbat) ||
    getTodayShabbatState(data?.civil_date, afterSunset).active;

  const groups = [];

  currentFeasts.forEach((feast) => {
    const feastName = feast?.name;
    const items = BIBLICAL_VERSE_LIBRARY.feasts[feastName] || [];

    if (items.length) {
      groups.push({
        title: `Versículos de ${feastName}`,
        items
      });
    }
  });

  if (omerDay > 0) {
    const omerItems =
      omerDay === 1
        ? (BIBLICAL_VERSE_LIBRARY.omer?.day1 || BIBLICAL_VERSE_LIBRARY.omer?.general || [])
        : (BIBLICAL_VERSE_LIBRARY.omer?.general || []);

    if (omerItems.length) {
      groups.push({
        title: omerDay === 1 ? "Versículos del día 1 del Omer" : `Versículos del Omer · día ${omerDay}`,
        items: omerItems
      });
    }
  }

  if (isShabbat) {
    groups.push({
      title: "Versículos para Shabat",
      items: BIBLICAL_VERSE_LIBRARY.shabbat || []
    });
  }

  if (data?.is_month_start && BIBLICAL_VERSE_LIBRARY.months[month]?.length) {
    groups.push({
      title: "Versículos del inicio del mes bíblico",
      items: BIBLICAL_VERSE_LIBRARY.months[month]
    });
  }

  if (!groups.length && BIBLICAL_VERSE_LIBRARY.months[month]?.length) {
    groups.push({
      title: "Versículos del mes bíblico",
      items: BIBLICAL_VERSE_LIBRARY.months[month]
    });
  }
  const randomVerses = getRandomContextVerses(data);

  if (randomVerses.length) {
    groups.push({
      title: "Cita adicional",
      items: randomVerses
    });
  }
  if (!groups.length) {
    groups.push({
      title: "Lectura sugerida para hoy",
      items: [
        {
          reference: "Salmo 119:105",
          text: "Lámpara es a mis pies tu palabra, y lumbrera a mi camino.",
          note: "Lectura base para cualquier día bíblico."
        },
        {
          reference: "Génesis 1:14",
          text: "Haya lumbreras en la expansión de los cielos para separar el día de la noche; y sirvan de señales para las estaciones, para días y años.",
          note: "Base del entendimiento bíblico de los tiempos señalados."
        }
      ]
    });
  }

  return groups.map(group => `
    <div class="verse-group">
      <h3 class="verse-group-title">${escapeHtml(group.title)}</h3>
      ${renderVerseItems(group.items)}
    </div>
  `).join("");
}

function renderBiblicalVersesSection(data) {
  const container = document.getElementById("biblicalVersesContent");
  const subtitle = document.getElementById("biblicalVersesSubtitle");

  if (!container) return;

  if (!data) {
    container.innerHTML = `
      <div class="verse-empty">
        No hay datos disponibles para mostrar versículos relacionados.
      </div>
    `;
    return;
  }

  const safeData = {
    biblical_month: Number(data?.biblical_month || 0),
    biblical_day: Number(data?.biblical_day || 0),
    civil_date: data?.civil_date || "",
    biblical_date: data?.biblical_date || "",
    after_sunset: Boolean(data?.after_sunset),
    current_feasts: Array.isArray(data?.current_feasts) ? data.current_feasts : [],
    is_shabbat: Boolean(data?.is_shabbat),
    is_month_start: Boolean(data?.is_month_start),
    month_start_type: data?.month_start_type || null,
    omer_day: Number(data?.omer_day || 0)
  };

  const html = buildBiblicalVersesHtml(safeData);

  container.innerHTML = html || `
    <div class="verse-empty">
      No hay versículos definidos para este día.
    </div>
  `;

  if (subtitle) {
    const parts = [
      `Fecha bíblica: ${safeData.biblical_date || "-"}`,
      `Fecha civil: ${safeData.civil_date || "-"}`
    ];

    if (safeData.omer_day > 0) {
      parts.push(`Omer: día ${safeData.omer_day}`);
    }

    subtitle.textContent = parts.join(" · ");
  }
}

function normalizeDateInput(dateStr) {
  const [year, month, day] = dateStr.split("-").map(Number);
  return new Date(year, month - 1, day);
}

function formatDateLabel(date) {
  return date.toLocaleDateString("es-EC", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

function formatLiveClock(dateObj) {
  return dateObj.toLocaleString("es-EC", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

function startJerusalemClock(initialDateTime) {
  const clockEl = document.getElementById("todayJerusalemClock");
  if (!clockEl || !initialDateTime) return;

  if (calendarState.jerusalemClockInterval) {
    clearInterval(calendarState.jerusalemClockInterval);
    calendarState.jerusalemClockInterval = null;
  }

  let liveTime = new Date(initialDateTime.replace(" ", "T"));

  if (Number.isNaN(liveTime.getTime())) {
    clockEl.textContent = initialDateTime;
    return;
  }

  clockEl.textContent = formatLiveClock(liveTime);

  calendarState.jerusalemClockInterval = setInterval(() => {
    liveTime = new Date(liveTime.getTime() + 1000);
    clockEl.textContent = formatLiveClock(liveTime);
  }, 1000);
}

function getWeekdayName(date) {
  return date.toLocaleDateString("es-EC", { weekday: "long" });
}

function capitalize(text) {
  if (!text) return "";
  return text.charAt(0).toUpperCase() + text.slice(1);
}

function getBiblicalMonthName(monthNumber) {
  const names = {
    1: "Aviv",
    2: "Segundo mes",
    3: "Tercer mes",
    4: "Cuarto mes",
    5: "Quinto mes",
    6: "Sexto mes",
    7: "Séptimo mes",
    8: "Octavo mes",
    9: "Noveno mes",
    10: "Décimo mes",
    11: "Undécimo mes",
    12: "Duodécimo mes",
  };

  return names[monthNumber] || `Mes ${monthNumber}`;
}

function getBiblicalDayWindowLabel(dateObj) {
  const startDay = new Date(dateObj);
  startDay.setDate(startDay.getDate() - 1);

  const startName = capitalize(getWeekdayName(startDay));
  const endName = capitalize(getWeekdayName(dateObj));

  return `${startName} al atardecer → ${endName} al atardecer`;
}

function formatDateKey(dateObj) {
  return [
    dateObj.getFullYear(),
    String(dateObj.getMonth() + 1).padStart(2, "0"),
    String(dateObj.getDate()).padStart(2, "0"),
  ].join("-");
}

function addDays(dateObj, days) {
  const d = new Date(dateObj);
  d.setDate(d.getDate() + days);
  return d;
}

function adjustToBiblicalStart(dateStr) {
  if (!dateStr) return null;
  return formatDateKey(addDays(normalizeDateInput(dateStr), -1));
}

function getFeastIcon(feastName) {
  const name = (feastName || "").toLowerCase();

  if (name.includes("pesaj")) return "🕊️";
  if (name.includes("levadura") || name.includes("matzot")) return "🫓";
  if (name.includes("bikkurim") || name.includes("primeros frutos")) return "🌾";
  if (name.includes("shavuot")) return "📜";
  if (name.includes("trompetas")) return "📯";
  if (name.includes("kippur")) return "🕯️";
  if (name.includes("tabern")) return "⛺";
  return "✦";
}

function getFeastTheme(feastName) {
  const name = (feastName || "").toLowerCase();

  if (name.includes("pesaj")) {
    return { cls: "theme-pesaj", short: "Pesaj" };
  }
  if (name.includes("levadura") || name.includes("matzot")) {
    return { cls: "theme-matzot", short: "Matzot" };
  }
  if (name.includes("bikkurim") || name.includes("primeros frutos")) {
    return { cls: "theme-bikkurim", short: "Bikkurim" };
  }
  if (name.includes("shavuot")) {
    return { cls: "theme-shavuot", short: "Shavuot" };
  }
  if (name.includes("trompetas")) {
    return { cls: "theme-trumpets", short: "Trompetas" };
  }
  if (name.includes("kippur")) {
    return { cls: "theme-kippur", short: "Kippur" };
  }
  if (name.includes("tabern")) {
    return { cls: "theme-sukkot", short: "Sukkot" };
  }

  return { cls: "theme-generic", short: feastName || "Fiesta" };
}

function getFeastOrder(feastName) {
  const name = (feastName || "").toLowerCase();

  if (name.includes("pesaj")) return 1;
  if (name.includes("levadura") || name.includes("matzot")) return 2;
  if (name.includes("bikkurim") || name.includes("primeros frutos")) return 3;
  if (name.includes("shavuot")) return 4;
  if (name.includes("trompetas")) return 5;
  if (name.includes("kippur")) return 6;
  if (name.includes("tabern")) return 7;

  return 99;
}

function shouldHideFeastBandLabel(feastName, spanDays) {
  const shortAlwaysVisible = ["Bikurim", "Pesaj", "Matzot", "Shavuot", "Shabat"];

  if (spanDays > 1) return false;
  if (shortAlwaysVisible.includes(feastName)) return false;

  return true;
}

function getBiblicalMonthStarts() {
  const starts = [];

  if (calendarState.biblicalMonthInfo?.month_start) {
    starts.push(calendarState.biblicalMonthInfo.month_start);
  }

  if (calendarState.biblicalMonthInfo?.next_month_start) {
    starts.push(calendarState.biblicalMonthInfo.next_month_start);
  }

  return starts;
}

function getMonthStartType(dateKey) {
  if (dateKey === calendarState.biblicalMonthInfo?.month_start) {
    return "current";
  }

  if (dateKey === calendarState.biblicalMonthInfo?.next_month_start) {
    return "next";
  }

  return null;
}

function isMonthStartDate(dateKey) {
  return getBiblicalMonthStarts().includes(dateKey);
}

async function loadBiblicalMonthInfo() {
  const payload = await fetchJson("/api/biblical/jerusalem/month");
  calendarState.biblicalMonthInfo = payload?.data || null;
}

function getTodayShabbatState(civilDate, afterSunset) {
  if (!civilDate) {
    return { active: false, upcoming: false, label: "" };
  }

  const dateObj = normalizeDateInput(civilDate);
  const dayOfWeek = dateObj.getDay();

  if ((dayOfWeek === 5 && afterSunset) || (dayOfWeek === 6 && !afterSunset)) {
    return {
      active: true,
      upcoming: false,
      label: "🕯️ Shabat activo",
    };
  }

  if (dayOfWeek === 5 && !afterSunset) {
    return {
      active: false,
      upcoming: true,
      label: "🌇 Hoy entra Shabat al atardecer",
    };
  }

  return { active: false, upcoming: false, label: "" };
}

function buildTodayBiblicalHtml(todayData, feastData, monthData) {
  const civilDate = feastData?.civil_date || todayData?.civil_date || monthData?.civil_date || "-";
  const biblicalDate = feastData?.biblical_date || todayData?.biblical_date || "-";
  const biblicalMonth = feastData?.biblical_month ?? todayData?.biblical_month ?? "-";
  const biblicalDay = feastData?.biblical_day ?? todayData?.biblical_day ?? "-";
  const jerusalemTime = monthData?.jerusalem_time || todayData?.jerusalem_time || "-";
  const sunsetTime = monthData?.sunset_time || todayData?.sunset_time || "-";
  const afterSunset = monthData?.after_sunset ?? todayData?.after_sunset ?? false;
  const currentFeasts = feastData?.current_feasts || [];
  const omer = feastData?.omer || null;
  const dayNote = feastData?.day_note || todayData?.day_note || monthData?.day_note || "";

  const shabbatState = getTodayShabbatState(civilDate, afterSunset);
  const biblicalMonthName =
    typeof biblicalMonth === "number" ? getBiblicalMonthName(biblicalMonth) : "-";
  const biblicalDisplay =
    biblicalDay !== "-" && biblicalMonthName !== "-"
      ? `${biblicalDay} de ${biblicalMonthName}`
      : "-";

  const badges = [];

  if (shabbatState.label) {
    badges.push(`<span class="today-status-badge is-shabbat">${shabbatState.label}</span>`);
  }

  if (monthData?.month_start === civilDate) {
    badges.push(`<span class="today-status-badge is-new-month">🌒 Cabeza del mes bíblico</span>`);
  }

  if (currentFeasts.length > 0) {
    currentFeasts.forEach((feast) => {
      badges.push(`
        <span class="today-status-badge is-feast">
          ${getFeastIcon(feast.name)} ${feast.name}
        </span>
      `);
    });
  }

  if (omer?.is_omer_counting && omer?.omer_day_today) {
    badges.push(`
      <span class="today-status-badge is-omer">
        🌾 Omer día ${omer.omer_day_today}
      </span>
    `);
  }

  if (badges.length === 0) {
    badges.push(`<span class="today-status-badge is-common">📖 Día común</span>`);
  }

  let summaryTitle = "Día bíblico en curso";
  let summaryText = "Hoy no hay una fiesta bíblica activa registrada para este momento.";

  if (currentFeasts.length > 0) {
    summaryTitle = currentFeasts.map((feast) => feast.name).join(" · ");
    summaryText = "Hoy hay una fiesta bíblica activa según el cálculo actual basado en Jerusalén.";
  } else if (shabbatState.active) {
    summaryTitle = "Shabat";
    summaryText = "Ahora mismo estás dentro de la ventana bíblica de Shabat.";
  } else if (shabbatState.upcoming) {
    summaryTitle = "Preparación para Shabat";
    summaryText = "Hoy todavía no ha comenzado Shabat, pero entra al atardecer.";
  }

  return `
    <div class="today-biblical-hero">
      <div class="today-biblical-main">
        <span class="today-biblical-kicker">Resumen del día</span>
        <h3>${summaryTitle}</h3>
        <p>${summaryText}</p>

        <div class="today-biblical-badges">
          ${badges.join("")}
        </div>
      </div>

      <div class="today-biblical-side">
        <div class="today-mini-card">
          <span class="today-mini-label">Fecha civil</span>
          <strong>${civilDate}</strong>
        </div>

        <div class="today-mini-card">
          <span class="today-mini-label">Fecha bíblica (equivalencia civil)</span>
          <strong>${biblicalDate}</strong>
        </div>

        <div class="today-mini-card">
          <span class="today-mini-label">Fecha bíblica</span>
          <strong>${biblicalDisplay}</strong>
        </div>

        <div class="today-mini-card">
          <span class="today-mini-label">Jerusalén</span>
          <strong id="todayJerusalemClock">${jerusalemTime}</strong>
        </div>
      </div>
    </div>

    <div class="today-biblical-grid">
      <div class="today-info-box">
        <span class="today-info-label">Atardecer en Jerusalén</span>
        <strong>${sunsetTime}</strong>
      </div>

      <div class="today-info-box">
        <span class="today-info-label">Después del atardecer</span>
        <strong>${afterSunset ? "Sí" : "No"}</strong>
      </div>

      <div class="today-info-box">
        <span class="today-info-label">Inicio del mes actual</span>
        <strong>${monthData?.month_start || "-"}</strong>
      </div>

      <div class="today-info-box">
        <span class="today-info-label">Próxima cabeza de mes</span>
        <strong>${monthData?.next_month_start || "-"}</strong>
      </div>
    </div>

    <div class="today-biblical-note">
      <strong>Nota del día:</strong>
      <span>
        La fecha civil cambia a medianoche; la fecha bíblica cambia al atardecer en Jerusalén.
        ${dayNote ? ` ${dayNote}` : ""}
      </span>
    </div>
  `;
}

async function loadTodayBiblicalPanel() {
  const panel = document.getElementById("todayBiblicalContent");
  if (!panel) return;

  panel.innerHTML = `<div class="today-biblical-loading">Cargando resumen bíblico de hoy...</div>`;

  try {
    const [todayPayload, feastPayload, monthPayload] = await Promise.all([
      fetchJson("/api/biblical/jerusalem/today"),
      fetchJson("/api/feasts/jerusalem"),
      fetchJson("/api/biblical/jerusalem/month"),
    ]);

    const todayData = todayPayload?.data || {};
    const feastData = feastPayload?.data || {};
    const monthData = monthPayload?.data || {};

    calendarState.biblicalMonthInfo = monthData;

    panel.innerHTML = buildTodayBiblicalHtml(todayData, feastData, monthData);
    startJerusalemClock(monthData?.jerusalem_time || todayData?.jerusalem_time || "");

    const shabbatState = getTodayShabbatState(
      feastData?.civil_date || todayData?.civil_date || monthData?.civil_date,
      monthData?.after_sunset ?? todayData?.after_sunset ?? false
    );

    const currentCivilDate =
      feastData?.civil_date || todayData?.civil_date || monthData?.civil_date || "";

    const verseData = {
      ...feastData,
      civil_date: currentCivilDate,
      biblical_date: feastData?.biblical_date || todayData?.biblical_date || "",
      biblical_month: feastData?.biblical_month ?? todayData?.biblical_month ?? 0,
      biblical_day: feastData?.biblical_day ?? todayData?.biblical_day ?? 0,
      after_sunset: monthData?.after_sunset ?? todayData?.after_sunset ?? false,
      is_shabbat: shabbatState.active,
      is_month_start: monthData?.month_start === currentCivilDate,
      month_start_type:
        monthData?.month_start === currentCivilDate
          ? "current"
          : monthData?.next_month_start === currentCivilDate
            ? "next"
            : null,
      omer_day: feastData?.omer?.omer_day_today || null
    };

    renderBiblicalVersesSection(verseData);

  } catch (error) {
    console.error("Error cargando el panel Hoy bíblico:", error);
    panel.innerHTML = `
      <div class="today-biblical-error">
        No se pudo cargar el resumen bíblico de hoy.
        <br>
        <small>${error.message}</small>
      </div>
    `;

    renderBiblicalVersesSection({
      biblical_month: 0,
      biblical_day: 0,
      civil_date: "",
      biblical_date: "",
      after_sunset: false,
      current_feasts: [],
      is_shabbat: false
    });
  }
}

function serializeFeastForDataset(feast) {
  return encodeURIComponent(JSON.stringify(feast));
}

function getFeastStableId(feast) {
  return [
    feast.name || "",
    feast.gregorian_start_date || "",
    feast.gregorian_end_date || "",
    feast.biblical_month || "",
    feast.biblical_day || "",
  ].join("|");
}

function isMobileCalendarView() {
  return window.innerWidth <= 768;
}

function buildMobileMarkers(feasts, isFriday, isSaturday, isBiblicalNewMonth) {
  const markers = [];

  if (feasts.length > 0) {
    const firstFeast = feasts
      .slice()
      .sort((a, b) => getFeastOrder(a.name) - getFeastOrder(b.name))[0];

    if (firstFeast) {
      const theme = getFeastTheme(firstFeast.name);
      const feastId = getFeastStableId(firstFeast);

      markers.push(`
        <span
          class="calendar-marker feast ${theme.cls} feast-hover-target"
          data-feast="${serializeFeastForDataset(firstFeast)}"
          data-feast-id="${feastId}"
          title="${firstFeast.name}"
        >
          ${getFeastIcon(firstFeast.name)}
        </span>
      `);
    }

    if (feasts.length > 1) {
      markers.push(`
        <span class="calendar-marker" title="${feasts.length - 1} evento(s) adicional(es)">
          +${Math.min(feasts.length - 1, 9)}
        </span>
      `);
    }
  }

  if (isFriday) {
    markers.push(`<span class="calendar-marker sunset-start" title="Desde este atardecer inicia el Shabat">🌇</span>`);
  }

  if (isSaturday) {
    markers.push(`<span class="calendar-marker shabbat" title="Shabat">🕯️</span>`);
  }

  if (isBiblicalNewMonth) {
    markers.push(`<span class="calendar-marker new-month" title="Cabeza del mes bíblico">🌒</span>`);
  }

  return markers.slice(0, 3).join("");
}

function groupFeastsByVisibleDay(data) {
  const map = new Map();
  const all = [
    ...(data.current_feasts || []),
    ...(data.upcoming_feasts || []),
  ];

  for (const feast of all) {
    const rawStartStr = feast.gregorian_start_date || feast.gregorian_date;
    const rawEndStr = feast.gregorian_end_date || feast.gregorian_start_date || feast.gregorian_date;

    if (!rawStartStr) continue;

    const shiftedStartStr = adjustToBiblicalStart(rawStartStr);
    const shiftedEndStr = adjustToBiblicalStart(rawEndStr);

    if (!shiftedStartStr || !shiftedEndStr) continue;

    const start = normalizeDateInput(shiftedStartStr);
    const end = normalizeDateInput(shiftedEndStr);

    let current = new Date(start);
    while (current <= end) {
      const key = formatDateKey(current);
      if (!map.has(key)) {
        map.set(key, []);
      }

      map.get(key).push({
        ...feast,
        display_start_date: shiftedStartStr,
        display_end_date: shiftedEndStr
      });

      current = addDays(current, 1);
    }
  }

  for (const [key, feasts] of map.entries()) {
    feasts.sort((a, b) => {
      const orderA = getFeastOrder(a.name);
      const orderB = getFeastOrder(b.name);
      if (orderA !== orderB) return orderA - orderB;

      const startA = a.display_start_date || a.gregorian_start_date || "";
      const startB = b.display_start_date || b.gregorian_start_date || "";
      return startA.localeCompare(startB);
    });
    map.set(key, feasts);
  }

  return map;
}

function ensureHoverCard() {
  let card = document.getElementById("feastHoverCard");
  if (!card) {
    card = document.createElement("div");
    card.id = "feastHoverCard";
    card.className = "feast-hover-card";
    document.body.appendChild(card);
  }
  return card;
}

function buildHoverCardHtml(feast) {
  const visualStart = feast.display_start_date || adjustToBiblicalStart(feast.gregorian_start_date || feast.gregorian_date) || feast.gregorian_start_date || feast.gregorian_date || "-";
  const visualEnd = feast.display_end_date || adjustToBiblicalStart(feast.gregorian_end_date || feast.gregorian_start_date || feast.gregorian_date) || feast.gregorian_end_date || feast.gregorian_start_date || feast.gregorian_date || "-";

  return `
    <div class="hover-card-title">${getFeastIcon(feast.name)} ${feast.name}</div>
    <div class="hover-card-line"><strong>Fecha bíblica:</strong> mes ${feast.biblical_month}, día ${feast.biblical_day}</div>
    <div class="hover-card-line"><strong>Visual en calendario:</strong> ${visualStart} → ${visualEnd}</div>
    <div class="hover-card-line"><strong>Fecha gregoriana API:</strong> ${feast.gregorian_start_date} → ${feast.gregorian_end_date}</div>
    <div class="hover-card-line"><strong>Horario bíblico:</strong> ${feast.biblical_start_at} → ${feast.biblical_end_at}</div>
    <div class="hover-card-line"><strong>Día:</strong> ${feast.weekday}</div>
    ${feast.description ? `<div class="hover-card-desc">${feast.description}</div>` : ""}
  `;
}

function positionHoverCard(mouseEvent) {
  const card = ensureHoverCard();
  const offsetX = 18;
  const offsetY = 18;

  let left = mouseEvent.clientX + offsetX;
  let top = mouseEvent.clientY + offsetY;

  const rect = card.getBoundingClientRect();
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  if (left + rect.width > viewportWidth - 12) {
    left = mouseEvent.clientX - rect.width - 16;
  }

  if (top + rect.height > viewportHeight - 12) {
    top = mouseEvent.clientY - rect.height - 16;
  }

  if (left < 12) left = 12;
  if (top < 12) top = 12;

  card.style.left = `${left}px`;
  card.style.top = `${top}px`;
}

function highlightFeastDays(feastId) {
  document.querySelectorAll(".calendar-day").forEach((day) => {
    day.classList.remove("is-feast-highlighted");
  });

  document.querySelectorAll(`.calendar-day[data-feast-ids*="${CSS.escape(feastId)}"]`).forEach((day) => {
    day.classList.add("is-feast-highlighted");
  });
}

function clearFeastHighlight() {
  document.querySelectorAll(".calendar-day").forEach((day) => {
    day.classList.remove("is-feast-highlighted");
  });
}

function showHoverCard(feast, mouseEvent) {
  const card = ensureHoverCard();
  card.innerHTML = buildHoverCardHtml(feast);
  card.classList.add("is-visible");
  positionHoverCard(mouseEvent);
}

function hideHoverCard() {
  const card = document.getElementById("feastHoverCard");
  if (!card) return;
  card.classList.remove("is-visible");
}

function bindFeastHoverCards() {
  const hoverTargets = document.querySelectorAll(".feast-hover-target");

  hoverTargets.forEach((target) => {
    target.addEventListener("mouseenter", (event) => {
      const encoded = target.dataset.feast;
      const feastId = target.dataset.feastId;
      if (!encoded) return;

      try {
        const feast = JSON.parse(decodeURIComponent(encoded));
        showHoverCard(feast, event);
        if (feastId) highlightFeastDays(feastId);
      } catch (err) {
        console.error("No se pudo leer data-feast:", err);
      }
    });

    target.addEventListener("mousemove", (event) => {
      positionHoverCard(event);
    });

    target.addEventListener("mouseleave", () => {
      hideHoverCard();
      clearFeastHighlight();
    });
  });
}

function bindShabbatHoverCard() {
  const hoverCard = document.getElementById("calendarHoverCard");
  if (!hoverCard) return;

  const fridayBridges = document.querySelectorAll(".shabbat-bridge");
  const saturdayMarkers = document.querySelectorAll(".calendar-marker.shabbat");

  function moveSimpleCard(e) {
    const offset = 12;

    hoverCard.style.display = "block";

    const rect = hoverCard.getBoundingClientRect();
    let x = e.clientX + offset;
    let y = e.clientY + offset;

    if (x + rect.width > window.innerWidth - 10) {
      x = e.clientX - rect.width - offset;
    }

    if (y + rect.height > window.innerHeight - 10) {
      y = e.clientY - rect.height - offset;
    }

    if (x < 10) x = 10;
    if (y < 10) y = 10;

    hoverCard.style.left = `${x}px`;
    hoverCard.style.top = `${y}px`;
  }

  function showShabbatCard(e) {
    hoverCard.textContent = "Shabat";
    hoverCard.style.display = "block";
    moveSimpleCard(e);
  }

  function hideShabbatCard() {
    hoverCard.style.display = "none";
  }

  fridayBridges.forEach((bridge) => {
    bridge.addEventListener("mouseenter", showShabbatCard);
    bridge.addEventListener("mousemove", moveSimpleCard);
    bridge.addEventListener("mouseleave", hideShabbatCard);
  });

  saturdayMarkers.forEach((marker) => {
    marker.addEventListener("mouseenter", showShabbatCard);
    marker.addEventListener("mousemove", moveSimpleCard);
    marker.addEventListener("mouseleave", hideShabbatCard);
  });
}

function renderDayDetail(dateKey, feastsMap) {
  const detail = document.getElementById("calendarDetailPanel");
  if (!detail || !dateKey) return;

  detail.classList.remove("is-hidden");
  detail.classList.add("is-visible");

  const dateObj = normalizeDateInput(dateKey);
  const feasts = feastsMap.get(dateKey) || [];
  const data = calendarState.feastData;
  const omer = data?.omer;
  const monthInfo = calendarState.biblicalMonthInfo;

  const dayOfWeek = dateObj.getDay();
  const isFriday = dayOfWeek === 5;
  const isSaturday = dayOfWeek === 6;
  const monthStartType = getMonthStartType(dateKey);
  const isBiblicalNewMonth = monthStartType !== null;

  const headerExtras = [];

  if (isFriday) {
    headerExtras.push("🌇 Desde este atardecer comienza la transición hacia Shabat");
  }

  if (isSaturday) {
    headerExtras.push(`🕯️ Shabat bíblico: ${getBiblicalDayWindowLabel(dateObj)}`);
  }

  if (isBiblicalNewMonth) {
    headerExtras.push(`🌒 Cabeza del mes bíblico: ${getBiblicalDayWindowLabel(dateObj)}`);
  }

  let currentDateInfoHtml = "";
  if (isBiblicalNewMonth) {
    currentDateInfoHtml = `
      <div class="feast-item">
        <h4>🌒 Este día marca el inicio del mes</h4>
        <p><strong>Ventana bíblica:</strong> ${getBiblicalDayWindowLabel(dateObj)}</p>
        <p><strong>Fecha gregoriana:</strong> ${dateKey}</p>
        <p>El mes bíblico comienza al atardecer de este día y continúa hasta el siguiente atardecer.</p>
      </div>
    `;
  } else if (isFriday) {
    currentDateInfoHtml = `
      <div class="feast-item">
        <h4>🌇 Inicio de Shabat al atardecer</h4>
        <p><strong>Ventana:</strong> ${capitalize(getWeekdayName(dateObj))} al atardecer → ${capitalize(getWeekdayName(addDays(dateObj, 1)))} al atardecer</p>
      </div>
    `;
  } else if (isSaturday) {
    currentDateInfoHtml = `
      <div class="feast-item">
        <h4>🕯️ Shabat bíblico</h4>
        <p><strong>Ventana:</strong> ${getBiblicalDayWindowLabel(dateObj)}</p>
      </div>
    `;
  }

  let monthInfoHtml = "";
  if (monthInfo) {
    monthInfoHtml = `
      <div class="feast-item">
        <h4>🌙 Información del mes bíblico actual</h4>
        <p><strong>Inicio del mes bíblico actual:</strong> ${monthInfo.month_start ?? "-"}</p>
        <p><strong>Inicio del próximo mes bíblico:</strong> ${monthInfo.next_month_start ?? "-"}</p>
        <p><strong>Día bíblico actual:</strong> ${monthInfo.biblical_day ?? "-"}</p>
        <p><strong>Hora en Jerusalén:</strong> ${monthInfo.jerusalem_time ?? "-"}</p>
        <p><strong>Atardecer en Jerusalén:</strong> ${monthInfo.sunset_time ?? "-"}</p>
        <p><strong>Después del atardecer:</strong> ${monthInfo.after_sunset ? "Sí" : "No"}</p>
        <p><strong>Posible día 1 hoy:</strong> ${monthInfo.is_possible_day_one ? "Sí" : "No"}</p>
        <p><strong>Posible día 1 próximo:</strong> ${monthInfo.is_possible_next_day_one ? "Sí" : "No"}</p>
        <p>${monthInfo.day_note ?? ""}</p>
      </div>
    `;
  }

  let feastHtml = "";
  if (feasts.length > 0) {
    feastHtml = feasts.map((feast) => `
      <div class="feast-item ${getFeastTheme(feast.name).cls}">
        <h4>${getFeastIcon(feast.name)} ${feast.name}</h4>
        <p><strong>Fecha bíblica:</strong> mes ${feast.biblical_month}, día ${feast.biblical_day}</p>
        <p><strong>Visual en calendario:</strong> ${feast.display_start_date || "-"} → ${feast.display_end_date || "-"}</p>
        <p><strong>Fecha gregoriana API:</strong> ${feast.gregorian_start_date} → ${feast.gregorian_end_date}</p>
        <p><strong>Horario bíblico:</strong> ${feast.biblical_start_at} → ${feast.biblical_end_at}</p>
        <p><strong>Día:</strong> ${feast.weekday}</p>
        <p>${feast.description || ""}</p>
      </div>
    `).join("");
  }

  let omerHtml = "";
  if (omer) {
    omerHtml = `
      <div class="feast-item">
        <h4>🌾 Conteo del Omer</h4>
        <p><strong>Primer Shabat en Panes sin Levadura:</strong> ${omer.first_shabbat_in_matzot ?? "-"}</p>
        <p><strong>Bikkurim:</strong> ${omer.bikkurim_date ?? "-"} (${omer.bikkurim_weekday ?? "-"})</p>
        <p><strong>Shavuot:</strong> ${omer.shavuot_date ?? "-"} (${omer.shavuot_weekday ?? "-"})</p>
        <p><strong>Día del Omer hoy:</strong> ${omer.omer_day_today ?? "-"}</p>
        <p><strong>Conteo activo:</strong> ${omer.is_omer_counting ? "Sí" : "No"}</p>
        <p>${omer.day_note ?? ""}</p>
      </div>
    `;
  }
const content = document.getElementById("calendarDetailContent");
const title = document.getElementById("calendarDetailTitle");

if (title) {
  title.textContent = `📅 ${formatDateLabel(dateObj)}`;
}

if (content) {
  content.innerHTML = `
    ${headerExtras.length ? `<p>${headerExtras.join(" · ")}</p>` : ""}
    ${currentDateInfoHtml}
    ${monthInfoHtml}
    ${feastHtml}
    ${omerHtml}
  `;
}
}

function buildWeekFeastBands(weekDates, feastsMap) {
  const allFeasts = new Map();

  weekDates.forEach((dateKey, colIndex) => {
    const feasts = feastsMap.get(dateKey) || [];
    feasts.forEach((feast) => {
      const feastId = getFeastStableId(feast);
      if (!allFeasts.has(feastId)) {
        allFeasts.set(feastId, {
          feast,
          cols: [],
        });
      }
      allFeasts.get(feastId).cols.push(colIndex);
    });
  });

  const sorted = Array.from(allFeasts.values()).sort((a, b) => {
    const orderA = getFeastOrder(a.feast.name);
    const orderB = getFeastOrder(b.feast.name);
    if (orderA !== orderB) return orderA - orderB;
    return (a.feast.gregorian_start_date || "").localeCompare(b.feast.gregorian_start_date || "");
  });

  const laneTopMap = {
    1: 42,
    2: 60,
    3: 78,
    4: 96,
    5: 42,
    6: 60,
    7: 78,
    99: 96,
  };

  return sorted.map((item) => {
    const { feast, cols } = item;
    const theme = getFeastTheme(feast.name);
    const feastId = getFeastStableId(feast);
    const minCol = Math.min(...cols);
    const maxCol = Math.max(...cols);
    const order = getFeastOrder(feast.name);
    const top = laneTopMap[order] || 64;
    const spanDays = maxCol - minCol + 1;
    const shortLabel = theme.short || feast.name || "Fiesta";
    const hideLabelClass = shouldHideFeastBandLabel(shortLabel, spanDays) ? "is-label-hidden" : "";

    const bandWidth =
      spanDays === 1
        ? "calc((100% / 7) - 12px)"
        : `calc(${spanDays} * (100% / 7) - 16px)`;

    const bandLeft =
      spanDays === 1
        ? `calc(${minCol} * (100% / 7) + 6px)`
        : `calc(${minCol} * (100% / 7) + 8px)`;

    return `
      <div
        class="week-feast-band ${theme.cls} feast-hover-target ${hideLabelClass} ${spanDays === 1 ? "is-single-day-band" : ""}"
        style="left: ${bandLeft}; width: ${bandWidth}; top: ${top}px;"
        data-feast="${serializeFeastForDataset(feast)}"
        data-feast-id="${feastId}"
        title="${feast.name}"
      >
        <span class="week-feast-band-label">${shortLabel}</span>
      </div>
    `;
  }).join("");
}

function renderCalendar(data) {
  const grid = document.getElementById("calendarGrid");
  const title = document.getElementById("calendarTitle");
  const detail = document.getElementById("calendarDetailPanel");
  if (!grid || !title) return;

  const feastsMap = groupFeastsByVisibleDay(data);
  const isMobile = isMobileCalendarView();

  const year = calendarState.year;
  const month = calendarState.month;

  const firstDay = new Date(year, month, 1);
  const startDay = firstDay.getDay();
  const daysInMonth = new Date(year, month + 1, 0).getDate();

  const todayKey = formatDateKey(new Date());

  title.textContent = firstDay.toLocaleString("es-ES", {
    month: "long",
    year: "numeric",
  });

  const cells = [];

  // Días vacíos al inicio
  for (let i = 0; i < startDay; i++) {
    cells.push({
      type: "empty",
      html: `<div class="calendar-day empty"></div>`,
    });
  }

  // Días del mes
  for (let day = 1; day <= daysInMonth; day++) {
    const dateObj = new Date(year, month, day);
    const dateKey = formatDateKey(dateObj);

    const feasts = feastsMap.get(dateKey) || [];
    const omerDayForThisDate = getOmerDayForDate(dateKey);

    const isToday = dateKey === todayKey;
    const isSelected = calendarState.selectedDate === dateKey;

    const dayOfWeek = dateObj.getDay();
    const isFriday = dayOfWeek === 5;
    const isSaturday = dayOfWeek === 6;

    const monthStartType = getMonthStartType(dateKey);
    const isBiblicalNewMonth = monthStartType !== null;

    const omerBand = omerDayForThisDate
      ? `
        <div class="omer-day-band" title="Omer día ${omerDayForThisDate}">
          ${isMobile ? `Omer ${omerDayForThisDate}` : `Omer día ${omerDayForThisDate}`}
        </div>
      `
      : "";

    const feastMarkers = !isMobile
      ? feasts
          .slice(0, 2)
          .sort((a, b) => getFeastOrder(a.name) - getFeastOrder(b.name))
          .map((f) => {
            const theme = getFeastTheme(f.name);
            const feastId = getFeastStableId(f);
            return `
              <span
                class="calendar-marker feast ${theme.cls} feast-hover-target"
                data-feast="${serializeFeastForDataset(f)}"
                data-feast-id="${feastId}"
                title="${f.name}"
              >
                ${getFeastIcon(f.name)}
              </span>
            `;
          })
          .join("")
      : buildMobileMarkers(feasts, isFriday, isSaturday, isBiblicalNewMonth);

    const fridaySunsetMarker =
      !isMobile && isFriday
        ? `<span class="calendar-marker sunset-start" title="Desde este atardecer inicia el Shabat">🌇</span>`
        : "";

    const shabbatMarker =
      !isMobile && isSaturday
        ? `<span class="calendar-marker shabbat" title="Shabat: viernes al atardecer → sábado al atardecer">🕯️</span>`
        : "";

    const newMonthMarker =
      !isMobile && isBiblicalNewMonth
        ? `<span class="calendar-marker new-month" title="Cabeza del mes bíblico">🌒</span>`
        : "";

    const shabbatBridge = isFriday
      ? `<div class="shabbat-bridge shabbat-span-2" title="Shabat">
            <span class="shabbat-bridge-label">${isMobile ? "" : "Shabat"}</span>
         </div>`
      : "";

    const feastIds = feasts.map((f) => getFeastStableId(f)).join("||");

    cells.push({
      type: "current",
      number: day,
      dateKey,
      html: `
        <div class="calendar-day 
          ${isToday ? "is-today" : ""} 
          ${feasts.length ? "has-feast" : ""} 
          ${isSelected ? "is-selected" : ""} 
          ${isFriday ? "is-friday" : ""} 
          ${isSaturday ? "is-shabbat-day" : ""} 
          ${monthStartType === "current" ? "is-new-month-day" : ""}
          ${monthStartType === "next" ? "is-next-new-month-day" : ""}"
          data-date="${dateKey}" 
          data-feast-ids="${feastIds}">

          ${shabbatBridge}
          ${omerBand}

          <div class="calendar-day-number">
            <div>${day}</div>
            ${
              omerDayForThisDate
                ? `<div class="omer-label">${isMobile ? omerDayForThisDate : `Omer: ${omerDayForThisDate}`}</div>`
                : ""
            }
          </div>

          <div class="calendar-markers">
            ${feastMarkers}
            ${fridaySunsetMarker}
            ${shabbatMarker}
            ${newMonthMarker}
          </div>

          <div class="calendar-day-note"></div>
        </div>
      `,
    });
  }

  // Render final
  grid.innerHTML = cells.map((c) => c.html).join("");

  // Click en día
  grid.querySelectorAll(".calendar-day[data-date]").forEach((el) => {
    el.addEventListener("click", () => {
      const dateKey = el.getAttribute("data-date");
      calendarState.selectedDate = dateKey;
      renderCalendar(data);
      renderDayDetail(dateKey, data);
    });
  });

  // Hover fiestas
  attachFeastHoverEvents();
}

function renderFeastsData(payload) {
  const data = payload.data;
  calendarState.feastData = data;

  const summary = document.getElementById("feastsSummary");

  if (summary) {
    summary.innerHTML = `
      <div class="feasts-summary">
        <p><strong>Fecha civil:</strong> ${data.civil_date}</p>
        <p><strong>Fecha bíblica:</strong> ${data.biblical_date}</p>
        <p><strong>Mes bíblico:</strong> ${data.biblical_month}</p>
        <p><strong>Día bíblico:</strong> ${data.biblical_day}</p>
        <p>${data.day_note || ""}</p>
      </div>
    `;
  }

  renderCalendar(data);
}

async function loadFeasts() {
  const summary = document.getElementById("feastsSummary");
  const detailPanel = document.getElementById("calendarDetailPanel");
  const grid = document.getElementById("calendarGrid");

  if (summary) {
    summary.innerHTML = `<div class="empty-state">Cargando información...</div>`;
  }

  if (detailPanel) {
    detailPanel.innerHTML = "Selecciona un día para ver detalles";
  }

  if (grid) {
    grid.innerHTML = "";
  }

  try {
    const [feastPayload] = await Promise.all([
      fetchJson("/api/feasts/jerusalem"),
      loadBiblicalMonthInfo(),
    ]);

    const civilDate = feastPayload?.data?.civil_date;
    calendarState.currentDate = civilDate ? normalizeDateInput(civilDate) : new Date();
    calendarState.selectedDate = civilDate || null;

    renderFeastsData(feastPayload);
  } catch (error) {
    console.error("Error cargando fiestas:", error);

    if (summary) {
      summary.innerHTML = `<div class="empty-state">${error.message}</div>`;
    }

    if (detailPanel) {
      detailPanel.innerHTML = "No se pudo cargar el detalle.";
    }

    if (grid) {
      grid.innerHTML = "";
    }
  }
}

function moveCalendarMonth(offset) {
  if (!calendarState.feastData) return;

  calendarState.currentDate = new Date(
    calendarState.currentDate.getFullYear(),
    calendarState.currentDate.getMonth() + offset,
    1
  );

  renderCalendar(calendarState.feastData);
}

function setupDetailToggle() {
  const panel = document.getElementById("calendarDetailPanel");
  const btn = document.getElementById("toggleDetailBtn");

  if (!panel || !btn) return;

  btn.addEventListener("click", () => {
    panel.classList.toggle("is-collapsed");

    if (panel.classList.contains("is-collapsed")) {
      btn.textContent = "+";
    } else {
      btn.textContent = "−";
    }
  });
}

document.addEventListener("DOMContentLoaded", () => {
  const btnRefreshTodayBiblical = document.getElementById("btnRefreshTodayBiblical");
  const btnLoadFeasts = document.getElementById("btnLoadFeasts");
  const btnPrevMonth = document.getElementById("btnPrevMonth");
  const btnNextMonth = document.getElementById("btnNextMonth");
  const btnTodayMonth = document.getElementById("btnTodayMonth");

  if (btnLoadFeasts) {
    btnLoadFeasts.addEventListener("click", loadFeasts);
  }

  if (btnRefreshTodayBiblical) {
    btnRefreshTodayBiblical.addEventListener("click", loadTodayBiblicalPanel);
  }

  if (btnPrevMonth) {
    btnPrevMonth.addEventListener("click", () => moveCalendarMonth(-1));
  }

  if (btnNextMonth) {
    btnNextMonth.addEventListener("click", () => moveCalendarMonth(1));
  }

  if (btnTodayMonth) {
    btnTodayMonth.addEventListener("click", () => {
      if (!calendarState.feastData) return;

      calendarState.currentDate = normalizeDateInput(calendarState.feastData.civil_date);
      calendarState.selectedDate = calendarState.feastData.civil_date;

      renderCalendar(calendarState.feastData);
      renderDayDetail(calendarState.selectedDate, groupFeastsByVisibleDay(calendarState.feastData));

      const todayVerseData = buildVerseContextForDate(calendarState.feastData.civil_date);
      renderBiblicalVersesSection(todayVerseData);
    });
  }
  loadTodayBiblicalPanel();
  setupDetailToggle();
  loadFeasts();
  window.addEventListener("resize", () => {
  if (calendarState.feastData) {
    renderCalendar(calendarState.feastData);
  }
});
});

