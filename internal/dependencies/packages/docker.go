package packages

import (
	"fmt"
	"os"

	"github.com/psprings/code-concierge/internal/utils"
)

// sudo apt-get install \
//     apt-transport-https \
//     ca-certificates \
//     curl \
//     gnupg-agent \
//     software-properties-common

func installDockerPackagePreReqs() ([]byte, error) {
	sudo := true
	preReqInstallCmd := []string{"apt-get", "install", "-y",
		"apt-transport-https",
		"ca-certificates",
		"curl",
		"gnupg-agent",
		"software-properties-common",
	}
	preReqInstallCmd = useSudo(preReqInstallCmd, sudo)
	binary, args := preReqInstallCmd[0], preReqInstallCmd[1:]
	stdout, err := utils.ShellCommandBufferedPrint(binary, args...)
	return stdout, err
}

func addDockerRepository() ([]byte, error) {
	sudo := true
	addDockerRepoCmd := []string{"add-apt-repository",
		`"deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"`,
	}

	addDockerRepoCmd = useSudo(addDockerRepoCmd, sudo)
	binary, args := addDockerRepoCmd[0], addDockerRepoCmd[1:]
	stdout, err := utils.ShellCommandBufferedPrint(binary, args...)
	return stdout, err
}

func addDockerGPGKey() ([]byte, error) {
	gpgCmd := []string{"curl", "-fsSL", "https://download.docker.com/linux/ubuntu/gpg", "|", "sudo", "apt-key", "add", "-"}
	binary, args := gpgCmd[0], gpgCmd[1:]
	stdout, err := utils.ShellCommandBufferedPrint(binary, args...)
	return stdout, err
}

func installDockerPreReqs() error {
	_, err := installDockerPackagePreReqs()
	if err != nil {
		return err
	}
	_, err = addDockerGPGKey()
	if err != nil {
		return err
	}
	_, err = addDockerRepository()
	if err != nil {
		return err
	}
	err = RunUpdate("")
	return err
}

// InstallDockerComposeCLI : install the docker-compose CLI
func InstallDockerComposeCLI(version string) error {
	if version == "" {
		version = "latest"
	}
	kernelName, err := utils.GetKernelName()
	if err != nil {
		return err
	}
	kernelVersion, err := utils.GetKernelVersion()
	if err != nil {
		return err
	}
	downloadFilename := "/usr/local/bin/docker-compose"
	downloadURL := fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/docker-compose-%s-%s", version, kernelName, kernelVersion)
	_, err = utils.DownloadFile(downloadFilename, downloadURL)
	if err != nil {
		return err
	}
	err = os.Chmod(downloadFilename, 0766)
	return err
}

// InstallDockerCLI : install the Docker CLI (not the daemon)
func InstallDockerCLI() error {
	err := installDockerPreReqs()
	if err != nil {
		return err
	}
	return Install("docker-ce-cli")
}

// InstallDocker : install Docker (with CLI and daemon)
func InstallDocker() error {
	err := installDockerPreReqs()
	if err != nil {
		return err
	}
	packagesToInstall := []string{"docker-ce", "docker-ce-cli", "containerd.io"}
	for _, packageName := range packagesToInstall {
		err = Install(packageName)
		if err != nil {
			return err
		}
	}
	return err
}
