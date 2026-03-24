const API_BASE = "https://orhalevanah-production.up.railway.app";

const MONTHS_ES = [
  "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
  "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"
];

const state = {
  currentView: "inicio",
  currentMonthDate: new Date(),
  selectedDate: null,
  todayPayload: null,
  todayFeastsPayload: null
};

async function fetchJson(url) {
  const response = await fetch(url);
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

function formatBiblicalDate(data) {
  if (!data) return "-";

  if (data.biblical_day && data.month_name) {
    return `${data.biblical_day} de ${data.month_name}`;
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
      <div class="list-card-title">${safeText(item.name, "Fiesta")}</div>
      <div class="list-card-meta">
        <span>Fecha: ${formatFeastDate(item)}</span>
      </div>
      ${
        item.description
          ? `<div class="list-card-text">${item.description}</div>`
          : ""
      }
    </div>
  `).join("");
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

async function getTodayBiblical() {
  return fetchJson(`${API_BASE}/api/biblical/jerusalem/today`);
}

async function getBiblicalByDate(date) {
  const candidates = [
    `${API_BASE}/api/biblical/jerusalem/date?date=${encodeURIComponent(date)}`,
    `${API_BASE}/api/biblical/date?date=${encodeURIComponent(date)}`,
    `${API_BASE}/api/biblical?date=${encodeURIComponent(date)}`
  ];

  let lastError = null;

  for (const url of candidates) {
    try {
      return await fetchJson(url);
    } catch (error) {
      lastError = error;
    }
  }

  throw lastError || new Error("No se pudo consultar la fecha bíblica");
}

async function getFeastsByDate(date) {
  const candidates = [
    `${API_BASE}/api/feasts?date=${encodeURIComponent(date)}`,
    `${API_BASE}/api/feasts/current?date=${encodeURIComponent(date)}`,
    `${API_BASE}/api/biblical/feasts?date=${encodeURIComponent(date)}`
  ];

  let lastError = null;

  for (const url of candidates) {
    try {
      return await fetchJson(url);
    } catch (error) {
      lastError = error;
    }
  }

  return {
    success: false,
    data: {
      current_feasts: [],
      upcoming_feasts: []
    },
    error: lastError
  };
}

function fillTodayUI(todayData) {
  document.getElementById("inicioCivilDate").textContent = safeText(todayData.civil_date, "Sin fecha");
  document.getElementById("inicioBiblicalDate").textContent = formatBiblicalDate(todayData);
  document.getElementById("inicioJerusalemTime").textContent = safeText(todayData.jerusalem_time, "No disponible");
  document.getElementById("inicioSunset").textContent = safeText(todayData.sunset_time, "No disponible");
  document.getElementById("inicioDayNote").textContent = safeText(todayData.day_note, "Información cargada correctamente.");

  document.getElementById("hoyCivilDate").textContent = safeText(todayData.civil_date, "Sin fecha");
  document.getElementById("hoyBiblicalDate").textContent = formatBiblicalDate(todayData);
  document.getElementById("hoyJerusalemTime").textContent = safeText(todayData.jerusalem_time, "No disponible");
  document.getElementById("hoySunsetTime").textContent = safeText(todayData.sunset_time, "No disponible");
  document.getElementById("hoyDayNote").textContent = safeText(todayData.day_note, "Sin nota disponible.");
}

async function loadTodayData() {
  try {
    const payload = await getTodayBiblical();
    const todayData = payload?.data || {};
    state.todayPayload = todayData;

    fillTodayUI(todayData);

    const todayDate = todayData.civil_date || formatDateToYMD(new Date());
    state.selectedDate = todayDate;

    const feastPayload = await getFeastsByDate(todayDate);
    state.todayFeastsPayload = feastPayload?.data || {
      current_feasts: [],
      upcoming_feasts: []
    };

    renderFeastList(
      "hoyCurrentFeasts",
      state.todayFeastsPayload.current_feasts || [],
      "No hay fiestas activas para este día."
    );

    renderFeastList(
      "hoyUpcomingFeasts",
      state.todayFeastsPayload.upcoming_feasts || [],
      "No se encontraron próximas fiestas."
    );

    renderFeastList(
      "fiestasCurrentList",
      state.todayFeastsPayload.current_feasts || [],
      "No hay fiestas activas para este día."
    );

    renderFeastList(
      "fiestasUpcomingList",
      state.todayFeastsPayload.upcoming_feasts || [],
      "No se encontraron próximas fiestas."
    );

    await loadSelectedDate(todayDate);
  } catch (error) {
    console.error("Error cargando hoy:", error);

    document.getElementById("inicioCivilDate").textContent = "Error";
    document.getElementById("inicioBiblicalDate").textContent = "Error";
    document.getElementById("inicioJerusalemTime").textContent = "Error";
    document.getElementById("inicioSunset").textContent = "Error";
    document.getElementById("inicioDayNote").textContent = error.message || "Failed to fetch";

    document.getElementById("hoyCivilDate").textContent = "Error";
    document.getElementById("hoyBiblicalDate").textContent = "Error";
    document.getElementById("hoyJerusalemTime").textContent = "Error";
    document.getElementById("hoySunsetTime").textContent = "Error";
    document.getElementById("hoyDayNote").textContent = error.message || "Failed to fetch";

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

  selectedCivilDate.textContent = dateString;
  selectedBiblicalDate.textContent = "Cargando...";
  selectedJerusalemTime.textContent = "Cargando...";
  selectedSunsetTime.textContent = "Cargando...";
  selectedDayNote.textContent = "Consultando detalle del día...";

  try {
    const payload = await getBiblicalByDate(dateString);
    const data = payload?.data || {};

    selectedCivilDate.textContent = safeText(data.civil_date, dateString);
    selectedBiblicalDate.textContent = formatBiblicalDate(data);
    selectedJerusalemTime.textContent = safeText(data.jerusalem_time, "No disponible");
    selectedSunsetTime.textContent = safeText(data.sunset_time, "No disponible");
    selectedDayNote.textContent = safeText(data.day_note, "Sin nota disponible.");
  } catch (error) {
    console.error("Error cargando fecha seleccionada:", error);
    selectedBiblicalDate.textContent = "No disponible";
    selectedJerusalemTime.textContent = "No disponible";
    selectedSunsetTime.textContent = "No disponible";
    selectedDayNote.textContent = error.message || "No se pudo cargar el detalle.";
  }
}

function renderCalendar() {
  const grid = document.getElementById("calendarGrid");
  const monthLabel = document.getElementById("calendarMonthLabel");
  const monthSubLabel = document.getElementById("calendarMonthSubLabel");

  const current = state.currentMonthDate;
  const year = current.getFullYear();
  const month = current.getMonth();

  monthLabel.textContent = `${MONTHS_ES[month]} ${year}`;
  monthSubLabel.textContent = "Calendario civil con referencia bíblica";

  grid.innerHTML = "";

  const firstDay = new Date(year, month, 1).getDay();
  const lastDate = new Date(year, month + 1, 0).getDate();
  const todayString = state.todayPayload?.civil_date || formatDateToYMD(new Date());

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

    cell.innerHTML = `
      <span class="calendar-day-number">${day}</span>
      <span class="calendar-day-label">${ymd}</span>
    `;

    cell.addEventListener("click", async () => {
      state.selectedDate = ymd;
      renderCalendar();
      await loadSelectedDate(ymd);
    });

    grid.appendChild(cell);
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
    state.selectedDate = state.todayPayload?.civil_date || formatDateToYMD(now);
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
  renderCalendar();

  if ("serviceWorker" in navigator) {
    navigator.serviceWorker.register("./sw.js").catch((error) => {
      console.warn("No se pudo registrar el Service Worker:", error);
    });
  }
});

// Mantener disponible globalmente por si algo del HTML lo usa
window.setView = setView;