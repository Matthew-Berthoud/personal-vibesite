package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"personal-vibesite/internal/github"
)

const GITHUB_USERNAME = "Matthew-Berthoud"
const PROJECT_NAMES = "projects.txt"

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

func GatherData() *PageData {
	repos, err := ReadLines(PROJECT_NAMES)
	if err != nil {
		log.Fatalf("failed to read repos: %v", err)
	}

	gh := github.NewGithubConnection(GITHUB_USERNAME)

	projects, err := gh.GetProjects(repos)
	if err != nil {
		log.Fatalf("Error getting projects: %v", err)
	}

	aboutMe, err := gh.GetAboutMe()
	if err != nil {
		log.Fatalf("Error getting About Me: %v", err)
	}

	return &PageData{
		Projects: projects,
		AboutMe:  aboutMe,
	}
}
