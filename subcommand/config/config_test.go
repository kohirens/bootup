package config

import (
	"github.com/kohirens/tmpltoapp/internal/cli"
	"os"
	"testing"
)

const (
	tmpDir = "tmp"
)

func TestMain(m *testing.M) {
	// Set up a temporary dir for generate files
	_ = os.RemoveAll(tmpDir)
	_ = os.Mkdir(tmpDir, cli.DirMode) // set up a temporary dir for generate files

	// Run all tests
	exitCode := m.Run()

	// Clean up
	os.Exit(exitCode)
}
func TestFlagsAndArguments(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		ca       []string
		expected string
		wantErr  bool
	}{
		{"noArgs", 1, []string{}, "usage: config", true},
		{"keyDoesNotExist", 1, []string{"set", "key", "value"}, "no config setting \"keyDoesNotExist\" found", false},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			Init()
			e := ParseInput(tc.ca)
			if !tc.wantErr && e != nil {
				t.Error(e)
			}
		})
	}
}

func TestSubCmdConfigBadExit(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		ca       []string
		expected string
		cfg      *cli.AppData
	}{
		{
			"set-exclude-file-extensions",
			1,
			[]string{"set", "ExcludeFileExtensions", "jpg,gif"},
			"no config setting \"keyDoesNotExist\" found",
			&cli.AppData{
				Path:    tmpDir + "/config-set-exclude-file-extensions.json",
				UsrOpts: &cli.UserOptions{&[]string{}, ""},
			},
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			Init()

			e2 := Run(tc.ca, tc.cfg)
			if e2 != nil {
				t.Error(e2)
			}
		})
	}
}
