package codeserver

import (
	"os"
	"path/filepath"
)

// DefaultDirectory : the default directory where code-server assets such as extensions should be installed
var DefaultDirectory = defaultDirectory()

func defaultDirectory() string {
	homeDir, _ := os.UserHomeDir()
	codeserverDir := filepath.Join(homeDir, ".code-server")
	return codeserverDir
}
