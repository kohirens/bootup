package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/path"
	stdt "github.com/kohirens/stdlib/test"
	"github.com/kohirens/tmpltoapp/internal/press"
	"github.com/kohirens/tmpltoapp/internal/test"
	"github.com/kohirens/tmpltoapp/subcommand/config"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	FixtureDir = "testdata"
	TmpDir     = "tmp"
)

func TestMain(m *testing.M) {
	stdt.RunMain(stdt.SubCmdFlags, main)

	stdt.ResetDir(TmpDir, 0774)
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}

func TestCallingMain(tester *testing.T) {
	dd := TmpDir + ps + tester.Name()
	_ = os.MkdirAll(dd, 0744)
	defer test.TmpSetParentDataDir(dd)()

	var tests = []struct {
		name     string
		wantCode int
		args     []string
	}{
		{"versionFlag", 0, []string{"-version"}},
		{"helpFlag", 0, []string{"-help"}},
		{
			"remoteGitTemplate",
			0,
			[]string{
				"-tmpl-path", "https://github.com/kohirens/tmpl-go-web.git",
				"-out-path", TmpDir + ps + "tmpl-go-web-03",
				"-answer-path", FixtureDir + ps + "answers-tmpl-go-web.json",
				"-tmpl-type", "git",
				"-branch", "refs/tags/0.3.0",
			},
		},
		{"manifest0", 0, []string{"manifest", "-h"}},
		{"manifest0", 1, []string{"manifest"}},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			_, _ = stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %v, want %v", got, tc.wantCode)
			}
		})
	}
}

// runMain execute main in a sub process
func runMain(testFunc string, args []string) *exec.Cmd {
	// Run the test binary and tell it to run just this test with environment set.
	cmd := exec.Command(os.Args[0], "-test.run", testFunc)

	subEnvVar := stdt.SubCmdFlags + "=" + strings.Join(args, " ")
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}

func TestSubCmdConfigExitCode(tester *testing.T) {
	var tests = []struct {
		name     string
		wantCode int
		args     []string
		expected string
	}{
		{"noArgs", 1, []string{config.Name}, "invalid number of arguments"},
		{"keyDoesNotExist", 1, []string{config.Name, "set", "key", "value"}, "no \"key\" setting found"},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			out, _ := stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			// using bytes.NewBuffer(out).String() allows the byte array to be
			// converted as-is and keeps chars such as newlines and quotes from
			// being converted to '\n' or '\"', etc.
			if !strings.Contains(bytes.NewBuffer(out).String(), tc.expected) {
				t.Errorf("std error from program did not contain %q", tc.expected)
			}
		})
	}
}

func TestSubCmdConfigHelp(tester *testing.T) {
	delayedFunc := test.TmpSetParentDataDir(TmpDir + "/TestSubCmdConfigSuccess")
	defer func() {
		delayedFunc()
	}()

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"help", 0, []string{config.Name, "-help"}, config.UsageTmpl},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			out, _ := stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			if !strings.Contains(bytes.NewBuffer(out).String(), tc.contains) {
				t.Errorf("did not contain %q", tc.contains)
			}
		})
	}
}

func TestSubCmdConfigSuccess(tester *testing.T) {
	delayedFunc := test.TmpSetParentDataDir(TmpDir + "/TestSubCmdConfigSuccess")
	defer func() {
		delayedFunc()
	}()

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		contains string
	}{
		{"getCache", 0, []string{config.Name, "get", "CacheDir"}, "tmp"},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			out, _ := stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			if !strings.Contains(bytes.NewBuffer(out).String(), tc.contains) {
				t.Errorf("did not contain %q", tc.contains)
			}
		})
	}
}

func TestSetUserOptions(tester *testing.T) {
	dd := TmpDir + "/TestSetUserOptions"
	_ = os.MkdirAll(dd, 0744)
	defer test.TmpSetParentDataDir(dd)()

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		want     string
	}{
		{"setCache", 0, []string{"-verbosity", "6", config.Name, "set", "CacheDir", "ABC123"}, "ABC123"},
		{
			"setExcludeFileExtensions",
			0,
			[]string{config.Name, "set", "ExcludeFileExtensions", "md,txt"},
			`"ExcludeFileExtensions":["md","txt"]`,
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {

			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			_, _ = stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			gotExit := cmd.ProcessState.ExitCode()

			if gotExit != tc.wantCode {
				t.Errorf("got %q, want %q", gotExit, tc.wantCode)
			}

			var confDir string

			switch runtime.GOOS {
			case "linux":
				confDir = dd + ps + ".config" + ps + AppName + ps + "config.json"
			default:
				confDir = dd + ps + AppName + ps + "config.json"
			}

			gotCfg, _ := press.LoadConfig(confDir)
			data, err1 := json.Marshal(gotCfg)

			if err1 != nil {
				t.Errorf("could not find config file %q", err1.Error())
			}
			if !strings.Contains(bytes.NewBuffer(data).String(), tc.want) {
				t.Errorf("the config %s did not contain %v set to %v", data, tc.args[2], tc.want)
			}
		})
	}
}

// Verify that same path error is thrown when the template path and the output
// are the same.
func TestTmplAndOutPathMatch(tester *testing.T) {
	dd := TmpDir + ps + tester.Name()
	_ = os.MkdirAll(dd, 0744)
	defer test.TmpSetParentDataDir(dd)()

	fixture := "repo-08"
	fixtureDir := dd + ps + "remotes" + ps + fixture

	var testCases = []struct {
		name     string
		wantCode int
		args     []string
		want     string
	}{
		{
			"inputOutputCollision",
			1,
			[]string{
				"-tmpl-path", fixtureDir,
				"-out-path", fixtureDir,
			},
			"template path and output path point to the same directory",
		},
	}

	for _, tc := range testCases {
		tester.Run(tc.name, func(t *testing.T) {
			tc.args[1] = git.CloneFromBundle(fixture, dd+ps+"remotes", FixtureDir, ps)

			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			sout, _ := stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			gotCode := cmd.ProcessState.ExitCode()

			if gotCode != tc.wantCode {
				t.Errorf("got %q, want %q", gotCode, tc.wantCode)
			}

			fmt.Println(string(sout))

			if strings.Contains(string(sout), tc.want) {
				t.Errorf("got %q, want %q", gotCode, tc.wantCode)
			}
		})
	}
}

// Check for the bug that occurs when there is no config file, and you
// want to process a template. This occurs most the first time you use the CLI.
func TestFirstTimeRun(tester *testing.T) {
	dd := TmpDir + ps + tester.Name()
	_ = os.MkdirAll(dd, 0744)
	defer test.TmpSetParentDataDir(dd)()

	fixture := "repo-07"

	var tests = []struct {
		name     string
		wantCode int
		args     []string
	}{
		{
			"pressTmplWithNoConfig",
			0,
			[]string{
				"-answer-path", FixtureDir + ps + fixture + "-answers.json",
				"-tmpl-path", "will be replace below",
				"-out-path", dd + ps + "processed" + ps + fixture,
			},
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			tc.args[3] = git.CloneFromBundle(fixture, dd+ps+"remotes", FixtureDir, ps)

			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			_, _ = stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %q, want %q", got, tc.wantCode)
			}

			gotPressedTmpl := tc.args[5]
			if !path.DirExist(gotPressedTmpl) {
				t.Errorf("output directory %q was not found", gotPressedTmpl)
			}
		})
	}
}

func TestSkipFeature(tester *testing.T) {
	dd := TmpDir + ps + tester.Name()
	_ = os.MkdirAll(dd, 0744)
	defer test.TmpSetParentDataDir(dd)()

	fixture := "repo-09"
	outPath := TmpDir + ps + "processed" + ps + fixture
	tc := struct {
		name     string
		wantCode int
		args     []string
		absent   []string
		present  []string
	}{
		"pressTmplWithNoConfig",
		0,
		[]string{
			"-answer-path", FixtureDir + ps + fixture + "-answers.json",
			"-out-path", outPath,
			"-tmpl-path", "will be replace below",
		},
		[]string{
			"dir-to-include/second-level/skip-me-as-well.md",
			"dir-to-skip",
			"skip-me-too.md",
			press.TmplManifestFile,
		},
		[]string{
			"dir-to-include/README.md",
			"dir-to-include/second-level/README.md",
			"README.md",
		},
	}

	tc.args[5] = git.CloneFromBundle(fixture, dd+ps+"remotes", FixtureDir, ps)

	cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

	_, _ = stdt.VerboseSubCmdOut(cmd.CombinedOutput())

	got := cmd.ProcessState.ExitCode()

	if got != tc.wantCode {
		tester.Errorf("got %q, but want %q", got, tc.wantCode)
	}

	for _, p := range tc.absent {
		file := outPath + ps + p
		if path.Exist(file) {
			tester.Errorf("file %q should NOT exist. check the skip code or test bundle %q", file, fixture)
		}
	}

	for _, p := range tc.present {
		file := outPath + ps + p
		if !path.Exist(file) {
			tester.Errorf("file %q should exist. check the skip code or test bundle %q", file, fixture)
		}
	}
}

func TestTmplPress(tester *testing.T) {
	op := TmpDir + ps + "out" + ps + "app-parse-dir-02"
	cd := TmpDir + ps + tester.Name()

	var tests = []struct {
		name     string
		wantCode int
		args     []string
		want     bool
	}{
		{
			"local",
			0,
			[]string{
				"-answer-path", FixtureDir + ps + "answers-parse-dir-02.json",
				"-out-path", op,
				"-tmpl-type", "git",
			},
			true,
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			rd := git.CloneFromBundle("parse-dir-02", cd, FixtureDir, ps)

			defer test.TmpSetParentDataDir(filepath.Dir(rd))()
			tc.args = append(tc.args, "-tmpl-path", rd)
			cmd := stdt.GetTestBinCmd(stdt.SubCmdFlags, tc.args)

			_, _ = stdt.VerboseSubCmdOut(cmd.CombinedOutput())

			// get exit code.
			got := cmd.ProcessState.ExitCode()

			if got != tc.wantCode {
				t.Errorf("got %v, want %v", got, tc.wantCode)
			}

			if path.Exist(op) != tc.want {
				t.Errorf("got %v, want %v", path.Exist(op), tc.want)
			}
		})
	}
}
