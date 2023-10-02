package cli

import (
	"testing"
)

func TestGetTmplLocation(runner *testing.T) {
	fixtures := []struct {
		name, want string
		tmplPath   string
	}{
		{"relative", "local", "./"},
		{"relative2", "local", "."},
		{"relativeUp", "local", ".."},
		{"absolute", "local", "/home/myuser"},
		{"windows", "local", "C:\\Temp"},
		{"http", "remote", "http://example.com/repo1"},
		{"https", "remote", "https://example.com/repo1"},
		{"git", "remote", "git://example.com/repo1"},
		{"file", "remote", "file://example.com/repo1"},
		{"hiddenRelative", "local", ".m/example.com/repo1"},
		{"tildeRelative", "local", "~/repo1.git"},
	}

	for _, tc := range fixtures {
		runner.Run(tc.name, func(t *testing.T) {
			got := getTmplLocation(tc.tmplPath)

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
