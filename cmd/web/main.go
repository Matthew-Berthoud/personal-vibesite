package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"

	"personal-vibesite/internal/github"
)

const GITHUB_USERNAME = "Matthew-Berthoud"
const PROJECT_NAMES = "projects.txt"

// PageData holds the data to be rendered in the HTML template.
type PageData struct {
	Projects []github.Project
	AboutMe  template.HTML
}

func ReadLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func main() {
	// Create a new file server to serve the static files from the current directory.
	gh := github.NewGithubConnection(GITHUB_USERNAME)
	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register a handler for the root URL that fetches and serves the project data.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		repos, err := ReadLines(PROJECT_NAMES)
		if err != nil {
			log.Fatalf("failed to read repos: %v", err)
		}

		projects, err := gh.GetProjects(repos)
		if err != nil {
			http.Error(w, "Failed to get project data from GitHub.", http.StatusInternalServerError)
			log.Printf("Error getting projects: %v", err)
			return
		}

		aboutMeMD, err := gh.GetReadMe(GITHUB_USERNAME)
		if err != nil {
			http.Error(w, "Failed to get About Me data from GitHub.", http.StatusInternalServerError)
			log.Printf("Error getting About Me: %v", err)
			return
		}
		htmlBytes := markdown.ToHTML([]byte(aboutMeMD), nil, nil)
		aboutMe := template.HTML(htmlBytes)

		data := PageData{
			Projects: projects,
			AboutMe:  aboutMe,
		}

		tmpl, err := template.ParseFiles("ui/html/pages/index.html", "ui/html/partials/project-template.html", "ui/html/partials/about-me.html")
		if err != nil {
			http.Error(w, "Failed to parse template files.", http.StatusInternalServerError)
			log.Printf("Error parsing templates: %v", err)
			return
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
