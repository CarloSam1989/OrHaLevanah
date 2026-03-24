package main

import (
	"log"
	"net/http"
	"os"

	"or-halevanah/internal/api"
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

	// Página principal
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Solo responde index en la raíz exacta
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./web/templates/index.html")
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