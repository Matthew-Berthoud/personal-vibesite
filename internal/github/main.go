package github

import (
	"context"
	"html/template"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/gomarkdown/markdown"
	"github.com/google/go-github/v74/github"
)

type Project struct {
	Name      string
	Overview  template.HTML
	GitHubURL string
}

type GithubConnection struct {
	Context context.Context
	Client  *github.Client
	User    string
}

func NewGithubConnection(user string) *GithubConnection {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Println("GITHUB_TOKEN environment variable not set. Using unauthenticated client.")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GithubConnection{
		Context: ctx,
		Client:  client,
		User:    user,
	}
}

func (g *GithubConnection) GetProjects(repos []string) ([]Project, error) {
	var projects []Project
	for _, repo := range repos {
		readme, err := g.GetReadMe(repo)
		if err != nil {
			log.Printf("Failed to get ReadMe for %s", repo)
			continue
		}

		overviewMarkdown := extractOverview(readme)
		htmlBytes := markdown.ToHTML([]byte(overviewMarkdown), nil, nil)
		overviewHTML := template.HTML(htmlBytes)

		projects = append(projects, Project{
			Name:      repo,
			Overview:  overviewHTML,
			GitHubURL: "https://github.com/" + g.User + "/" + repo,
		})
	}
	return projects, nil
}

func (g *GithubConnection) GetAboutMe() (template.HTML, error) {
	aboutMeMD, err := g.GetReadMe(g.User)
	if err != nil {
		return "", err
	}

	bodyMarkdown := extractSection(aboutMeMD, "Hi there 👋")
	htmlBytes := markdown.ToHTML([]byte(bodyMarkdown), nil, nil)

	return template.HTML(htmlBytes), nil
}

func (g *GithubConnection) GetReadMe(repo string) (string, error) {
	readmeEncoded, _, err := g.Client.Repositories.GetReadme(g.Context, g.User, repo, nil)
	if err != nil {
		return "", err
	}

	readme, err := readmeEncoded.GetContent()
	if err != nil {
		return "", err
	}

	return readme, err
}

func extractSection(md string, sectionTitle string) string {
	titleStart := "# " + sectionTitle
	startIdx := strings.Index(md, titleStart)
	if startIdx == -1 {
		return "No " + sectionTitle + " section found."
	}
	startIdx += len(titleStart)
	md = md[startIdx:]

	// Find the end of the section (next heading, or end of file).
	endIdx := strings.Index(md, "\n#")
	if endIdx == -1 {
		// If no next heading is found, take the rest of the content.
		return strings.TrimSpace(md)
	}

	return strings.TrimSpace(md[:endIdx])
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
