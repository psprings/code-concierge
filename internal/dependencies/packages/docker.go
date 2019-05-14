package packages

import "github.com/psprings/code-concierge/internal/utils"

// sudo apt-get install \
//     apt-transport-https \
//     ca-certificates \
//     curl \
//     gnupg-agent \
//     software-properties-common

func installDockerPreReqs() ([]byte, error) {
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

// InstallDockerCLI :
func InstallDockerCLI() {
	installDockerPreReqs()
	addDockerGPGKey()
	addDockerRepository()
	RunUpdate("")
	Install("docker-ce-cli")
}
