package main

import (
	"bufio"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/google/go-github/v74/github"
)

// Project represents the data for a single project card.
type Project struct {
	Name      string
	Overview  template.HTML
	GitHubURL string
}

// PageData holds the data to be rendered in the HTML template.
type PageData struct {
	Projects []Project
}

func main() {
	// Create a new file server to serve the static files from the current directory.
	fs := http.FileServer(http.Dir("ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register a handler for the root URL that fetches and serves the project data.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		projects, err := getProjectsFromGitHub()
		if err != nil {
			http.Error(w, "Failed to get project data from GitHub.", http.StatusInternalServerError)
			log.Printf("Error getting projects: %v", err)
			return
		}

		// Prepare the data for the template.
		data := PageData{
			Projects: projects,
		}

		// Parse the index.html and project-template.html files.
		tmpl, err := template.ParseFiles("ui/html/pages/index.html", "ui/html/partials/project-template.html")
		if err != nil {
			http.Error(w, "Failed to parse template files.", http.StatusInternalServerError)
			log.Printf("Error parsing templates: %v", err)
			return
		}

		// Execute the template with the project data.
		err = tmpl.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			http.Error(w, "Failed to render template.", http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	})

	log.Println("Server listening on http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProjectsFromGitHub() ([]Project, error) {
	var projects []Project

	// Read the list of repositories from projects.txt.
	file, err := os.Open("projects.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new GitHub client without authentication.
	// This will work for public repos but is subject to rate limits.
	ctx := context.Background()
	client := github.NewClient(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		repoName := strings.TrimSpace(scanner.Text())
		if repoName == "" {
			continue
		}

		// Get the README content from the specified repository.
		readme, _, err := client.Repositories.GetReadme(ctx, "Matthew-Berthoud", repoName, nil)
		if err != nil {
			log.Printf("Error fetching README for %s: %v", repoName, err)
			continue
		}

		// Decode the base64 content.
		content, err := readme.GetContent()
		if err != nil {
			log.Printf("Error decoding content for %s: %v", repoName, err)
			continue
		}

		// Extract the "Overview" section.
		overviewMarkdown := extractOverview(content)

		// Convert Markdown to HTML.
		htmlBytes := markdown.ToHTML([]byte(overviewMarkdown), nil, nil)
		overviewHTML := template.HTML(htmlBytes)

		projects = append(projects, Project{
			Name:      repoName,
			Overview:  overviewHTML,
			GitHubURL: "https://github.com/Matthew-Berthoud/" + repoName,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func extractOverview(md string) string {
	// Simple string manipulation to find the content between "## Overview" and the next heading.
	overviewStart := "## Overview"
	startIdx := strings.Index(md, overviewStart)
	if startIdx == -1 {
		return "No overview section found."
	}
	startIdx += len(overviewStart)
	md = md[startIdx:]

	// Find the end of the overview section (next heading, or end of file).
	endIdx := strings.Index(md, "\n#")
	if endIdx == -1 {
		// If no next heading is found, take the rest of the content.
		return strings.TrimSpace(md)
	}

	return strings.TrimSpace(md[:endIdx])
}
