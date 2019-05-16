package utils

import (
	"archive/zip"
	"bufio"
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

// DownloadFile : download the contents from a URL to a specified filename
func DownloadFile(filename string, urlString string) (string, error) {
	return downloadFile(filename, urlString)
}

// DownloadTmpFile : download the contents from a given URL to a temporary file
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

// UnzipFile : given a .zip file name, unzip to a destination folder
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

// ShellCommand : execute shell commands
// Mostly just using for tests at the moment
// example usage:
// shellCommand("ls", "-la")
func ShellCommand(binary string, args ...string) ([]byte, error) {
	return shellCommand(false, binary, args...)
}

// ShellCommandBufferedPrint : print stdout as the command is running instead of waiting for it to run
func ShellCommandBufferedPrint(binary string, args ...string) ([]byte, error) {
	return shellCommand(true, binary, args...)
}

// func ShellCommand(binary string, args ...string) ([]byte, error) {
// 	cmd := exec.Command(binary, args...)
// 	cmd.Stderr
// 	stdout, err := cmd.Output()
// 	return stdout, err
// }

func shellCommand(bufferedLog bool, binary string, args ...string) ([]byte, error) {
	cmd := exec.Command(binary, args...)
	stdoutReader, err := cmd.StdoutPipe()
	cmd.Start()
	if bufferedLog {
		scanner := bufio.NewScanner(stdoutReader)
		// scanner.Split(bufio.ScanWords)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
		}
	}
	cmd.Wait()
	stdout, _ := ioutil.ReadAll(stdoutReader)
	return stdout, err
}

// GitConfigureAuth : inject API token into GitHub https calls
func GitConfigureAuth(repoURL string, apiToken string) {
	if apiToken == "" {
		return
	}
	parsedURL, _ := url.Parse(repoURL)
	parsedURL.Path = ""

	println(parsedURL)
	substituteValue := fmt.Sprintf("url.%s://\"${GITHUB_API_TOKEN}:x-oauth-basic@%s/\".insteadOf", parsedURL.Scheme, parsedURL.Host)
	forValue := fmt.Sprintf("%s://%s/", parsedURL.Scheme, parsedURL.Host)
	os.Setenv("GITHUB_API_TOKEN", apiToken)
	ShellCommand("git", "config", "--global", substituteValue, forValue)
}

func shellCommandPrint(binary string, args ...string) {
	stdout, err := ShellCommand(binary, args...)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(stdout))
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
	ShellCommandBufferedPrint("git", "init", ".")
	ShellCommandBufferedPrint("git", "remote", "add", "origin", repoURL)
	ShellCommandBufferedPrint("git", "pull", "origin", useBranch)
}

// UniqueStrings : deduplicate items in a string slice
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

// GetKernelName : returns the kernel name
func GetKernelName() (string, error) {
	stdout, error := ShellCommand("uname", "-s")
	return string(stdout), error
}

// GetKernelVersion : returns the kernel version
func GetKernelVersion() (string, error) {
	stdout, error := ShellCommand("uname", "-r")
	return string(stdout), error
}
