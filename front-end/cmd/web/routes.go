package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func routes(brokerURL string) *chi.Mux {
	router := chi.NewRouter()

	// Set up middleware
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
	}))

	// Set up routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Define a closure function that captures the brokerURL variable
		renderWithBrokerURL := func(w http.ResponseWriter, t string) {
			render(w, t, map[string]interface{}{
				"brokerURL": brokerURL,
			})
		}

		// Call the closure function instead of the render function
		renderWithBrokerURL(w, "test.page.gohtml")
	})

	return router
}

func render(w http.ResponseWriter, t string, data map[string]interface{}) {
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	templateSlice := append([]string{fmt.Sprintf("./cmd/web/templates/%s", t)}, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}