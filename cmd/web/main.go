package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	data := GatherData()

	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles("ui/html/pages/index.html", "ui/html/partials/project-template.html", "ui/html/partials/about-me.html")
		if err != nil {
			log.Fatalf("Error parsing templates: %v", err)
		}

		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, "Failed to render template.", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	})

	log.Println("Server listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
