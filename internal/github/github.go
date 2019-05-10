package github

import (
	"context"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// RepoURL : the GitHub repository URL taken from the environment variable `GITHUB_REPO_URL`
var RepoURL = os.Getenv("GITHUB_REPO_URL")

// APIToken : the GitHub personal access / API token taken from the environment variable `GITHUB_API_TOKEN`
var APIToken = os.Getenv("GITHUB_API_TOKEN")

// DefaultAPIURL : the public GitHub API URL to use
var DefaultAPIURL = "https://api.github/com"

func getAPIURLFromRepo(repoURL string) string {
	parsedRepoURL, err := url.Parse(repoURL)
	if err != nil {
		log.Fatal("error: ", err)
	}
	if parsedRepoURL.Host == "github.com" {
		parsedRepoURL.Host = "api.github.com"
		parsedRepoURL.Path = ""
		return parsedRepoURL.String()
	}
	parsedRepoURL.Path = "/api/v3"
	return parsedRepoURL.String()
}

// GetAPIURL :
func GetAPIURL(repoURL string) string {
	githubEnvAPIURL := os.Getenv("GITHUB_API_URL")
	if repoURL == "" && githubEnvAPIURL == "" {
		return DefaultAPIURL
	}
	if len(githubEnvAPIURL) > 0 {
		return githubEnvAPIURL
	}
	return getAPIURLFromRepo(repoURL)
}

func ownerRepoFromURL(repoURL string) (string, string) {
	parsedRepoURL, err := url.Parse(repoURL)
	if err != nil {
		log.Fatal("error: ", err)
	}
	repoPath := parsedRepoURL.Path
	pathSplit := strings.Split(repoPath, "/")
	return pathSplit[1], pathSplit[2]
}

func formatAPIBaseURL(apiURL string) string {
	// The GitHub library being used requires the base url to end in a "/"
	if strings.HasSuffix(apiURL, "/") {
		return apiURL
	}
	return apiURL + "/"
}

// Config :
type Config struct {
	APIBaseURL string
	Token      string
}

// GetLanguagesFromURL : given a (HTTPS) GitHub repo URL, return the languages
func (c *Config) GetLanguagesFromURL(repoURL string) (map[string]int, error) {
	owner, repo := ownerRepoFromURL(repoURL)
	return c.ListLanguages(owner, repo)
}

// ListLanguages : retrieve the code languages associated with a given GitHub repository
func (c *Config) ListLanguages(owner string, repo string) (map[string]int, error) {
	// tp := github.BasicAuthTransport{
	// 	Username: strings.TrimSpace(username),
	// 	Password: strings.TrimSpace(password),
	// }
	// client := github.NewClient(tp.Client())

	var client *github.Client
	ctx := context.Background()
	if c.Token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: c.Token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}
	if c.APIBaseURL != DefaultAPIURL {
		formattedBaseURL := formatAPIBaseURL(c.APIBaseURL)
		u, err := url.Parse(formattedBaseURL)
		if err != nil {
			log.Fatal(err)
		}
		client.BaseURL = u
	}

	languages, _, err := client.Repositories.ListLanguages(ctx, owner, repo)
	return languages, err
}
