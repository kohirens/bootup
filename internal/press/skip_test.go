package press

import (
	"path/filepath"
	"testing"
)

func Test_inSkipArray(t *testing.T) {
	tests := []struct {
		name     string
		f        string
		patterns []string
		want     bool
	}{
		{
			"lvl-1-wildcard",
			"README.md",
			[]string{"*.md"},
			true,
		},
		{
			"exact-match",
			"README.md",
			[]string{"README.md"},
			true,
		},
		{
			"nested",
			filepath.Join("dir1", "dir2", "dir3", "fake.txt"),
			[]string{filepath.Join("dir1", "**", "dir3", "*.txt")},
			true,
		},
		{
			"nested",
			filepath.Join("dir1", "dir2", "dir3", "skip-me-too.md"),
			[]string{"*skip-me-too.md"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inSkipArray(tt.f, tt.patterns); got != tt.want {
				t.Errorf("inSkipArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
