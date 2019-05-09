package dependencies

import "testing"

func TestParseLanguagesDependencies(t *testing.T) {
	jsonString := `
	{
		"NodeJS": {
			"Extensions": [
				"eg2.vscode-npm-script",
				"esbenp.prettier-vscode"
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
			"Extensions": [
				"redhat.java"
			]
		}
	}
	`
	tt := []struct {
		name       string
		jsonString string
		parsed     map[string]Dependencies
	}{
		{"sample map", jsonString, map[string]Dependencies{
			"Go":   Dependencies{Extensions: []string{"ms-vscode.go"}},
			"Java": Dependencies{Extensions: []string{"redhat.java"}},
		},
		},
	}

	for _, tc := range tt {
		for language, deps := range tc.parsed {
			parsed, _ := ParseLanguagesDependencies(tc.jsonString)
			t.Run(tc.name, func(t *testing.T) {
				parsedDeps := parsed[language]
				if parsedDeps.Extensions[0] != deps.Extensions[0] {
					t.Errorf("%s: expected dependencies for language %s to be %s; got %s", tc.name, language, deps, parsedDeps)
				}
			})
		}
	}
}
