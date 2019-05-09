package github

import (
	"testing"
)

func TestGetAPIURLFromRepo(t *testing.T) {
	tt := []struct {
		name    string
		repoURL string
		value   string
	}{
		{"publicgithub", "https://github.com/psprings/code-concierge", "https://api.github.com"},
		{"enterprisegithub", "https://github.pretendenterprise.com/psprings/code-concierge", "https://github.pretendenterprise.com/api/v3"},
	}

	for _, tc := range tt {
		testAPIURL := getAPIURLFromRepo(tc.repoURL)
		t.Run(tc.name, func(t *testing.T) {
			if testAPIURL != tc.value {
				t.Errorf("%s: expected value to be %s; got %s", tc.name, tc.value, testAPIURL)
			}
		})
	}
}
func TestOrgRepoFromURL(t *testing.T) {
	tt := []struct {
		name    string
		repoURL string
		org     string
		repo    string
	}{
		{"publicgithub", "https://github.com/psprings/code-concierge", "psprings", "code-concierge"},
		{"enterprisegithub", "https://github.pretendenterprise.com/foo/bar", "foo", "bar"},
	}

	for _, tc := range tt {
		org, repo := orgRepoFromURL(tc.repoURL)
		t.Run(tc.name, func(t *testing.T) {
			if org != tc.org && repo != tc.repo {
				t.Errorf("%s: expected org to be %s and repo to be %s; got org %s and repo %s", tc.name, tc.org, tc.repo, org, repo)
			}
		})
	}
}
