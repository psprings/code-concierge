package packages

import (
	"errors"
	"log"
	"strings"

	"github.com/psprings/code-concierge/internal/utils"
)

var EnableSudo = true

func isRisky(packageString string) bool {
	// If has space, there could be nefarious things
	return strings.ContainsAny(packageString, " ")
}

func installCommandWrapper(operatingSystem string, packageName string) []string {
	switch operatingSystem {
	case "ubuntu":
		return []string{"apt-get", "install", "-y", packageName}
	default:
		return []string{"apt-get", "install", "-y", packageName}
	}
}

func useSudo(command []string, sudo bool) []string {
	if !sudo {
		return command
	}
	command = append([]string{"sudo"}, command...)
	return command
}

// RunUpdate : package manager update
func RunUpdate(operatingSystem string) error {
	var command []string
	switch operatingSystem {
	case "ubuntu":
		command = []string{"apt-get", "update"}
	default:
		command = []string{"apt-get", "update"}
	}
	sudo := EnableSudo
	command = useSudo(command, sudo)
	binary, args := command[0], command[1:]
	// stdout, err := utils.ShellCommand(binary, args...)
	stdout, err := utils.ShellCommandBufferedPrint(binary, args...)
	log.Print(string(stdout))
	return err
}

// Install : install the given package on the host
func Install(packageString string) error {
	if isRisky(packageString) {
		return errors.New("contains patterns that could pose a security risk")
	}
	// Parameterize this later
	// currently needed to do installs via the `coder` user
	sudo := EnableSudo
	command := installCommandWrapper("", packageString)
	command = useSudo(command, sudo)
	binary, args := command[0], command[1:]
	stdout, err := utils.ShellCommand(binary, args...)
	log.Print(string(stdout))
	return err
}
