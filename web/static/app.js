const btnLoadToday = document.getElementById("btnLoadToday");
const btnLoadMonth = document.getElementById("btnLoadMonth");
const btnLoadFeasts = document.getElementById("btnLoadFeasts");
const btnLoadAvivStatus = document.getElementById("btnLoadAvivStatus");
const btnSendAvivReport = document.getElementById("btnSendAvivReport");
const btnConfirmAviv = document.getElementById("btnConfirmAviv");
const btnRejectAviv = document.getElementById("btnRejectAviv");
const btnResetAviv = document.getElementById("btnResetAviv");

const outputToday = document.getElementById("outputToday");
const outputMonth = document.getElementById("outputMonth");
const outputFeasts = document.getElementById("outputFeasts");
const outputAvivStatus = document.getElementById("outputAvivStatus");
const outputAvivReport = document.getElementById("outputAvivReport");
const outputAvivAdmin = document.getElementById("outputAvivAdmin");

async function loadJson(url, outputElement) {
  outputElement.textContent = "Cargando...";

  try {
    const res = await fetch(url);
    const text = await res.text();

    try {
      const data = JSON.parse(text);
      outputElement.textContent = JSON.stringify(data, null, 2);
    } catch {
      outputElement.textContent = text || "Respuesta no JSON";
    }
  } catch (err) {
    outputElement.textContent = "Error al cargar datos";
    console.error(err);
  }
}

async function postJson(url, body, outputElement) {
  outputElement.textContent = "Procesando...";

  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(body)
    });

    const text = await res.text();

    try {
      const data = JSON.parse(text);
      outputElement.textContent = JSON.stringify(data, null, 2);
    } catch {
      outputElement.textContent = text || "Respuesta no JSON";
    }
  } catch (err) {
    outputElement.textContent = "Error al procesar solicitud";
    console.error(err);
  }
}

btnLoadToday?.addEventListener("click", () => {
  loadJson("/api/biblical/jerusalem/today", outputToday);
});

btnLoadMonth?.addEventListener("click", () => {
  loadJson("/api/biblical/jerusalem/month", outputMonth);
});

btnLoadFeasts?.addEventListener("click", () => {
  loadJson("/api/feasts/jerusalem", outputFeasts);
});

btnLoadAvivStatus?.addEventListener("click", () => {
  loadJson("/api/aviv/status", outputAvivStatus);
});

btnSendAvivReport?.addEventListener("click", () => {
  const userName = document.getElementById("avivUserName").value;
  const confirmed = document.getElementById("avivConfirmed").value === "true";
  const notes = document.getElementById("avivNotes").value;
  const evidenceURL = document.getElementById("avivEvidenceURL").value;

  postJson("/api/aviv/report", {
    user_name: userName,
    confirmed: confirmed,
    notes: notes,
    evidence_url: evidenceURL
  }, outputAvivReport);
});

btnConfirmAviv?.addEventListener("click", () => {
  const adminName = document.getElementById("adminName").value;

  postJson("/api/admin/aviv/confirm", {
    admin_name: adminName
  }, outputAvivAdmin);
});

btnRejectAviv?.addEventListener("click", () => {
  const adminName = document.getElementById("adminName").value;

  postJson("/api/admin/aviv/reject", {
    admin_name: adminName
  }, outputAvivAdmin);
});

btnResetAviv?.addEventListener("click", () => {
  postJson("/api/admin/aviv/reset", {}, outputAvivAdmin);
});