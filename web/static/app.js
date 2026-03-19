const btnLoadToday = document.getElementById("btnLoadToday");
const btnLoadMonth = document.getElementById("btnLoadMonth");
const outputToday = document.getElementById("outputToday");
const outputMonth = document.getElementById("outputMonth");

async function loadJson(url, outputElement) {
  outputElement.textContent = "Cargando...";

  try {
    const res = await fetch(url);
    const data = await res.json();
    outputElement.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    outputElement.textContent = "Error al cargar datos";
    console.error(err);
  }
}

btnLoadToday?.addEventListener("click", () => {
  loadJson("/api/biblical/jerusalem/today", outputToday);
});

btnLoadMonth?.addEventListener("click", () => {
  loadJson("/api/biblical/jerusalem/month", outputMonth);
});