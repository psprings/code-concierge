package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DownloadFile :
func DownloadFile(filename string, urlString string) (string, error) {
	return downloadFile(filename, urlString)
}

// DownloadTmpFile :
func DownloadTmpFile(urlString string) (string, error) {
	return downloadFile("", urlString)
}

func tmpFileWrapper(filename string) (*os.File, error) {
	if filename != "" {
		file, err := os.Create(filename)
		return file, err
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), "code-concierge-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}
	// defer os.Remove(tmpFile.Name())
	return tmpFile, err
}

// downloadFile :
func downloadFile(filename string, urlString string) (string, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := tmpFileWrapper(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return file.Name(), err
}

// UnzipFile :
func UnzipFile(filename string, destination string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(filename)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, file := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(destination, file.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		// Create directory
		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Create file
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		createFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(createFile, rc)

		// Close the file immediately because of for loop
		createFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

// Mostly just using for tests at the moment
// example usage:
// shellCommand("ls", "-la")
// ShellCommand : execute shell commands
func ShellCommand(binary string, args ...string) {
	cmd := exec.Command(binary, args...)
	stdout, err := cmd.Output()
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Println(string(stdout))
}

// GitConfigureAuth : inject API token into GitHub https calls
func GitConfigureAuth(repoURL string, apiToken string) {
	if apiToken == "" {
		return
	}
	parsedURL, _ := url.Parse(repoURL)
	parsedURL.Path = ""

	println(parsedURL)
	substituteValue := fmt.Sprintf("url.%s://\"${GITHUB_TOKEN}:x-oauth-basic@%s/\".insteadOf", parsedURL.Scheme, parsedURL.Host)
	forValue := fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host)
	os.Setenv("GITHUB_API_TOKEN", apiToken)
	ShellCommand("git", "config", "--global", substituteValue, forValue)
}

// GitClone : execute `git clone`
func GitClone(repoURL string, apiToken string) {
	if apiToken != "" {
		GitConfigureAuth(repoURL, apiToken)
	}
	useBranch := os.Getenv("GITHUB_REPO_BRANCH")
	if useBranch == "" {
		useBranch = "master"
	}
	ShellCommand("git", "init", ".")
	ShellCommand("git", "remote", "add", "origin", repoURL)
	ShellCommand("git", "pull", "origin", useBranch)
}

// UniqueStrings :
func UniqueStrings(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
