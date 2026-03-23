package main

import (
	"log"
	"net/http"
	"os"

	"or-halevanah/internal/api"
)

func main() {
	// Puerto dinámico para Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Archivos estáticos
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Página principal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Solo responde index en la raíz exacta
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, "./web/templates/index.html")
	})

	// Registrar rutas API
	api.RegisterRoutes()

	log.Printf("Servidor corriendo en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
