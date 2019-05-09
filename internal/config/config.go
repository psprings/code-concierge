package config

import (
	"flag"
	"os"
	"strings"

	"github.com/psprings/code-concierge/internal/github"
)

// Config :
type Config struct {
	APIBaseURL string
	Token      string
	// Username   string
	// Password   string
	RepoURL              string
	DependenciesURL      string
	AdditionalExtensions []string
}

func ensureAPIToken(token string) string {
	if token != "" {
		return token
	}
	return github.APIToken
}

func ensureRepoURL(repoURL string) string {
	if repoURL != "" {
		return repoURL
	}
	return github.RepoURL
}

func ensureAddExts(additionalExtensions string) string {
	if additionalExtensions != "" {
		return additionalExtensions
	}
	return os.Getenv("CODE_SERVER_INSTALL_EXTENSIONS")
}

func commaSplit(commaSeparated string) []string {
	var separated []string
	if commaSeparated == "" {
		return separated
	}
	return strings.Split(commaSeparated, ",")
}

// Retrieve :
func Retrieve() *Config {
	var apiURL string
	var apiToken string
	var repoURL string
	var addExts string
	flag.StringVar(&apiURL, "api-url", "https://api.github.com", "(optional) The base URL for the GitHub API")
	flag.StringVar(&apiToken, "api-token", "", "The token to use for authentication to GitHub")
	flag.StringVar(&repoURL, "repo-url", "", "The (https) URL of the GitHub repo to use")
	flag.StringVar(&addExts, "additional-extensions", "", "Comma separated list of extension IDs to install")
	flag.Parse()
	additionalExtensions := ensureAddExts(addExts)
	return &Config{
		APIBaseURL:           apiURL,
		Token:                ensureAPIToken(apiToken),
		RepoURL:              ensureRepoURL(repoURL),
		AdditionalExtensions: commaSplit(additionalExtensions),
	}
}
