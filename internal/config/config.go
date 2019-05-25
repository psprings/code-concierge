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
	AdditionalPackages   []string
	InstallDockerCLI     bool
	SkipAutoInstalls     bool
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

func ensureGithubAPIURL(apiURL string, repoURL string) string {
	if apiURL != "" {
		return apiURL
	}
	inferredAPIURL := github.GetAPIURL(repoURL)
	return inferredAPIURL
}

func ensureAddExts(additionalExtensions string) string {
	if additionalExtensions != "" {
		return additionalExtensions
	}
	return os.Getenv("CODE_SERVER_INSTALL_EXTENSIONS")
}

func ensureAddPackage(additionalPackages string) string {
	if additionalPackages != "" {
		return additionalPackages
	}
	return os.Getenv("CODE_SERVER_INSTALL_PACKAGES")
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
	var addPackages string
	var installDockerCLI bool
	var skipAutoInstalls bool
	flag.StringVar(&apiURL, "api-url", "", "(optional) The base URL for the GitHub API")
	flag.StringVar(&apiToken, "api-token", "", "The token to use for authentication to GitHub")
	flag.StringVar(&repoURL, "repo-url", "", "The (https) URL of the GitHub repo to use")
	flag.StringVar(&addExts, "additional-extensions", "", "Comma separated list of extension IDs to install")
	flag.StringVar(&addPackages, "additional-packages", "", "Comma separated list of packages to install")
	flag.BoolVar(&installDockerCLI, "install-docker", false, "Whether to install the Docker CLI")
	flag.BoolVar(&skipAutoInstalls, "skip-auto-installs", false, "Whether to skip automatic install of extensions and packages")
	flag.Parse()
	additionalExtensions := ensureAddExts(addExts)
	additionalPackages := ensureAddExts(addPackages)
	eRepoURL := ensureRepoURL(repoURL)
	eAPIBaseURL := ensureGithubAPIURL(apiURL, eRepoURL)
	return &Config{
		APIBaseURL:           eAPIBaseURL,
		Token:                ensureAPIToken(apiToken),
		RepoURL:              eRepoURL,
		AdditionalExtensions: commaSplit(additionalExtensions),
		AdditionalPackages:   commaSplit(additionalPackages),
		InstallDockerCLI:     installDockerCLI,
		SkipAutoInstalls:     skipAutoInstalls,
	}
}
