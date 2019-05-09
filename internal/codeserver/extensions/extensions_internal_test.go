package extensions

import (
	"testing"
)

func TestGetConfig(t *testing.T) {
	tt := []struct {
		name        string
		extensionID string
		config      *Config
	}{
		{"golang extension", "ms-vscode.go", &Config{Publisher: "ms-vscode", Extension: "go", Version: "latest"}},
		{"javascript extension - version", "ms-vscode.typescript-javascript-grammar@1.0.0", &Config{Publisher: "ms-vscode", Extension: "typescript-javascript-grammar", Version: "1.0.0"}},
	}

	for _, tc := range tt {
		config := GetConfig(tc.extensionID)
		t.Run(tc.name, func(t *testing.T) {
			if *config != *tc.config {
				t.Errorf("%s: expected %#v; got %#v", tc.name, tc.config, config)
			}
		})
	}
}

func TestExtensionDirectoryName(t *testing.T) {
	tt := []struct {
		name          string
		config        *Config
		directoryName string
	}{
		{"[extension directory] golang extension", &Config{Publisher: "ms-vscode", Extension: "go", Version: "latest"}, "ms-vscode.go-latest"},
		{"[extension directory] javascript extension (version)", &Config{Publisher: "ms-vscode", Extension: "typescript-javascript-grammar", Version: "1.0.0"}, "ms-vscode.typescript-javascript-grammar-1.0.0"},
	}

	for _, tc := range tt {
		directoryName := tc.config.extensionDirectoryName()
		t.Run(tc.name, func(t *testing.T) {
			if directoryName != tc.directoryName {
				t.Errorf("%s: expected %s; got %s", tc.name, tc.directoryName, directoryName)
			}
		})
	}
}
