package extensions

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/psprings/code-concierge/internal/codeserver"
	"github.com/psprings/code-concierge/internal/utils"
)

// Config : basic properties of an extension
type Config struct {
	Publisher string
	Extension string
	Version   string
}

// GetConfig : return extension configuration
func GetConfig(extensionID string) *Config {
	extensionIDBase := extensionID
	version := "latest"
	splitVersion := strings.Split(extensionIDBase, "@")
	if len(splitVersion) > 1 {
		extensionIDBase = splitVersion[0]
		version = splitVersion[1]
	}
	splitExtension := strings.Split(extensionIDBase, ".")
	publisher := splitExtension[0]
	name := splitExtension[1]

	return &Config{
		Publisher: publisher,
		Extension: name,
		Version:   version,
	}
}

func getDownloadURL(publisher string, extension string, version string) string {
	if version == "" {
		version = "latest"
	}
	return fmt.Sprintf("https://%s.gallery.vsassets.io/_apis/public/gallery/publisher/%s/extension/%s/%s/assetbyname/Microsoft.VisualStudio.Services.VSIXPackage", publisher, publisher, extension, version)
}

func getExtensionsDirectory() string {
	// for testing
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	codeserverDir := codeserver.DefaultDirectory
	return filepath.Join(codeserverDir, "extensions")
}

func (c *Config) extensionDirectoryName() string {
	return fmt.Sprintf("%s.%s-%s", c.Publisher, c.Extension, c.Version)
}

func (c *Config) destinationDirectory() string {
	extensionsDirectory := getExtensionsDirectory()
	extensionFolderDir := c.extensionDirectoryName()
	makeDirectory := filepath.Join(extensionsDirectory, extensionFolderDir)
	os.MkdirAll(makeDirectory, os.ModePerm)
	return makeDirectory
}

// Install :
func (c *Config) Install() error {
	extensionURL := getDownloadURL(c.Publisher, c.Extension, c.Version)
	// Download the file
	filename, err := utils.DownloadTmpFile(extensionURL)
	defer os.Remove(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// For temporary testing
	destination := c.destinationDirectory()

	// Unzip to extensions directory
	_, err = utils.UnzipFile(filename, destination)

	return err
}
