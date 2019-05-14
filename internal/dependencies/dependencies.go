package dependencies

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/psprings/code-concierge/internal/codeserver/extensions"
	"github.com/psprings/code-concierge/internal/config"
	"github.com/psprings/code-concierge/internal/dependencies/packages"
	"github.com/psprings/code-concierge/internal/github"
	"github.com/psprings/code-concierge/internal/utils"
)

// DefaultDependenciesString : JSON representation of the default dependencies map
var DefaultDependenciesString = `
{
    "NodeJS": {
        "Packages": [
            "nodejs",
            "npm"
        ],
        "Extensions": [
            "eg2.vscode-npm-script",
            "esbenp.prettier-vscode"
        ]
    },
    "Python": {
        "Extensions": [
            "ms-python.python"
        ]
    },
    "Typescript": {
        "Inherit": "NodeJS",
        "Extensions": [
            "ms-vscode.typescript-javascript-grammar"
        ]
    },
    "CSS": {
        "Extensions": [
            "pranaygp.vscode-css-peek"
        ]
    },
    "Dockerfile": {
		"Packages": [
            "docker-ce-cli"
        ],
        "Extensions": [
            "peterjausovec.vscode-docker"
        ]
    },
    "Vue": {
        "Extensions": [
            "octref.vetur"
        ]
    },
    "HTML": {
        "Extensions": [
            "ritwickdey.liveserver"
        ]
    },
    "Go": {
        "Packages": [
            "golang-1.12-go"
        ],
        "Extensions": [
            "ms-vscode.go"
        ]
    },
    "Markdown": {
        "Extensions": [
            "davidanson.vscode-markdownlint",
            "yzhang.markdown-all-in-one"
        ]
    },
    "Shell": {},
    "C++": {
        "Extensions": [
            "ms-vscode.cpptools"
        ]
    },
    "Java": {
        "Packages": [
            "default-jdk"
        ],
        "Extensions": [
            "redhat.java"
        ]
    }
}
`

// Dependencies :
type Dependencies struct {
	Extensions []string
	Packages   []string
}

// ParseLanguagesDependencies : given a JSON string, return a map of languages
// and extensions
func ParseLanguagesDependencies(jsonString string) (map[string]Dependencies, error) {
	var langDepMap map[string]Dependencies
	err := json.Unmarshal([]byte(jsonString), &langDepMap)
	return langDepMap, err
}

func getDependenciesFromURL(externalURL string) (string, error) {
	resp, err := http.Get(externalURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	return bodyString, err
}

// GetLanguagesDependencies :
func GetLanguagesDependencies(externalURL string) (map[string]Dependencies, error) {
	dependenciesMapString := DefaultDependenciesString
	if externalURL != "" {
		dms, err := getDependenciesFromURL(externalURL)
		if err != nil {
			log.Fatal(err)
		}
		dependenciesMapString = dms
	}
	return ParseLanguagesDependencies(dependenciesMapString)
}

func installExtension(extensionID string) {
	ec := extensions.GetConfig(extensionID)
	err := ec.Install()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Installed extension: %s", extensionID)
	}
}

func installExtensions(extensionIDs []string) {
	extensionIDs = utils.UniqueStrings(extensionIDs)
	for _, extensionID := range extensionIDs {
		installExtension(extensionID)
	}
}

func installPackages(packageList []string) {
	if len(packageList) < 1 {
		return
	}
	packageList = utils.UniqueStrings(packageList)
	packages.RunUpdate("")
	for _, packageName := range packageList {
		err := packages.Install(packageName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Install :
func Install() {
	c := config.Retrieve()
	gc := github.Config{
		APIBaseURL: c.APIBaseURL,
		Token:      c.Token,
	}
	languages, err := gc.GetLanguagesFromURL(c.RepoURL)
	if err != nil {
		log.Fatal(err)
	}

	langDepMap, err := GetLanguagesDependencies(c.DependenciesURL)

	var allExtensions []string
	var allPackages []string

	// Iterate through all languages discovered in GitHub
	for language := range languages {
		log.Printf("Language: %s", language)
		currentDeps := langDepMap[language]
		if _, ok := langDepMap[language]; !ok {
			continue
		}
		// Populate extension ID list
		for _, extensionID := range currentDeps.Extensions {
			allExtensions = append(allExtensions, extensionID)
		}
		// Populate package install list
		for _, packageName := range currentDeps.Packages {
			allPackages = append(allPackages, packageName)
		}
	}
	// Merge automatic and user provided extension lists
	allExtensions = append(allExtensions, c.AdditionalExtensions...)
	// De-duplicate and install extensions
	installExtensions(allExtensions)
	// Merge automatic and user provided packages
	allPackages = append(allPackages, c.AdditionalPackages...)
	// De-duplicate and install packages
	installPackages(allPackages)
	// Install Docker CLI
	if c.InstallDockerCLI {
		packages.InstallDockerCLI()
	}

	// Clone repo from GitHub
	utils.GitClone(c.RepoURL, c.Token)
}
