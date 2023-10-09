package main

import (
	"testing"
)

func TestGetRepoDir(tester *testing.T) {
	var testCases = []struct {
		name    string
		repo    string
		refName string
		want    string
	}{
		{"local", "/repo-01", "", "repo-01"},
		{"url", "https://example.com/repo-01.git", "", "repo-01"},
		{"bareRepo", "/repo-01.git", "", "repo-01"},
		{"bareRepo", "/repo-02.git", "refs/heads/main", "repo-02-refs-heads-main"},
		{"bareRepo", "/repo-02.git", "refs/tags/0.1.0", "repo-02-refs-tags-0.1.0"},
		{"bareRepo", "/repo-02.git", "refs/remotes/origin/HEAD", "repo-02-refs-remotes-origin-HEAD"},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {

			got := getRepoDir(tc.repo, tc.refName)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
