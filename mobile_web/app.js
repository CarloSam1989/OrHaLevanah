const API_BASE = "https://orhalevanah-production.up.railway.app/";

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

document.addEventListener("DOMContentLoaded", async () => {
  const fecha = document.getElementById("fecha");
  const mensaje = document.getElementById("mensaje");

  try {
    mensaje.innerText = "Probando conexión...";

    const payload = await fetchJson(`${API_BASE}/api/biblical/jerusalem/today`);
    const data = payload?.data || {};

    fecha.innerText = data.civil_date || "Sin fecha";
    mensaje.innerText = data.day_note || "Conexión correcta";
  } catch (error) {
    console.error("Error fetch:", error);
    fecha.innerText = "Error de conexión";
    mensaje.innerText = error.message || "Failed to fetch";
  }
});

function setView(view) {
  alert("Vista: " + view);
}