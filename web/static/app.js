const btn = document.getElementById("btnLoad");
const output = document.getElementById("output");

btn.addEventListener("click", async () => {
  output.textContent = "Cargando...";

  try {
    const res = await fetch("/api/biblical/jerusalem/today");
    const data = await res.json();

    output.textContent = JSON.stringify(data, null, 2);
  } catch (err) {
    output.textContent = "Error al cargar datos";
    console.error(err);
  }
});