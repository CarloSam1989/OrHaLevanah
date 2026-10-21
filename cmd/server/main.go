package main

import (
	"log"
	"net/http"
	"os"

	"or-halevanah/internal/api"
	"or-halevanah/internal/sighting"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	stopSightings := make(chan struct{})
	sighting.StartAutoRefresh(stopSightings)
	defer close(stopSightings)

	// Puerto dinámico para Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Crear mux principal
	mux := http.NewServeMux()

	// Archivos estáticos
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Página pública de la comunidad y aplicación existente.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			http.ServeFile(w, r, "./web/templates/community.html")
		case "/app":
			http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
		case "/app/":
			http.ServeFile(w, r, "./web/templates/index.html")
		default:
			http.NotFound(w, r)
		}
	})

	// Registrar rutas API sobre el mux global temporalmente
	// y luego redirigirlas al mux principal
	api.RegisterRoutes()

	// Copiar las rutas registradas en DefaultServeMux al mux principal
	mux.Handle("/api/", http.DefaultServeMux)

	// Aplicar CORS
	handler := corsMiddleware(mux)

	log.Printf("Servidor corriendo en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
