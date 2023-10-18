package press

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/path"
	test2 "github.com/kohirens/stdlib/test"
	"os"
	"path/filepath"
	"testing"
)

func xTestFindTemplates(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
		want    int
	}{
		{"1 file", test2.FixtureDir + PS + "find-me-01", false, 1},
		//{"2 files", test.FixturesDir + PS, false, 2},
		//{"5 files", test.FixturesDir + PS, false, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e1 := FindTemplates(tt.dir)

			if (e1 != nil) != tt.wantErr {
				t.Errorf("FindTemplates() error = %v, wantErr %v", e1.Error(), tt.wantErr)
			}

			if len(got) < tt.want {
				t.Errorf("got %v, want %v", len(got), tt.want)
			}
		})
	}
}

// TestParseDir2 Test the empty directory feature. The fixture directory
// "template-04" contains a directory "dir1" and has a file ".empty". After
// calling press.Print on the fixture directory, the output directory
// "template-04-out" should container "dir1" without any files or directories.
func TestEmptyDirectoryFeature(runner *testing.T) {
	// Abort in verbose mode.
	_, df := test2.TempFileSwap(&os.Stdout, runner.Name(), "out")
	defer df()

	fixtures := []struct {
		dstDir,
		name,
		srcDir string
		want func() bool
		vars cli.StringMap
	}{
		{
			test2.TmpDir + PS + "template-04-out",
			"dir1IsEmpty",
			test2.FixtureDir + "/template-04",
			func() bool {
				dir := test2.TmpDir + PS + "template-04-out" + PS + "dir1"
				fs, e1 := os.ReadDir(dir)
				if e1 != nil {
					return false
				}
				return len(fs) == 0
			},
			cli.StringMap{},
		},
	}

	fileChkr, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	for _, tc := range fixtures {
		runner.Run(tc.name, func(t *testing.T) {
			e1 := Print(tc.srcDir, tc.dstDir, tc.vars, fileChkr, &templateJson{Excludes: []string{}})

			if e1 != nil {
				t.Errorf("got error %v, want nil", e1.Error())
			}

			want := tc.want()
			fmt.Printf("want = %v\n", want)
			if !want {
				t.Error("all is not good")
			}
		})
	}
}

// TestPrinting Verify the print-head is printing.
//
//	Just making sure when calling press.Print that files make it through to
//	the template engine and	parsing happens as expected.
func TestPrinting(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(test2.FixtureDir + PS + "parse-dir-01")
	tmpDir, _ := filepath.Abs(test2.TmpDir)
	tests := []struct {
		name, tmplPath, outPath string
		tplVars                 cli.StringMap
		fileToCheck, want       string
	}{
		{
			"parse-dir-01", fixturePath1, tmpDir + PS + "parse-dir-01",
			cli.StringMap{"APP_NAME": "SolarPolar"},
			tmpDir + "/parse-dir-01/dir1/README.md", "SolarPolar\n",
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(test *testing.T) {
			fec, _ := stdlib.NewFileExtChecker(nil, &[]string{"md", "yml"})

			err := Print(tc.tmplPath, tc.outPath, tc.tplVars, fec, &templateJson{})

			if err != nil {
				test.Errorf("got an error %q", err.Error())
				return
			}

			got, err := os.ReadFile(tmpDir + "/parse-dir-01/dir1/README.md")

			if err != nil {
				test.Errorf("got an error %q", err.Error())
			}

			if string(got) != tc.want {
				test.Errorf("got %s, but want %v", got, tc.want)
			}
		})
	}
}

func TestSkipping(tester *testing.T) {
	repoFixture := "repo-10"
	outPath := test2.TmpDir + PS + "processed" + PS + repoFixture
	fecFixture, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	tc := struct {
		name    string
		absent  []string
		present []string
		answers cli.StringMap
		ph      *templateJson
	}{
		"pressTmplWithNoConfig",
		[]string{
			"dir-to-include/second-level/skip-me-as-well.md",
			"dir-to-skip",
			"skip-me-too.md",
			TmplManifestFile,
		},
		[]string{
			"dir-to-include/README.md",
			"dir-to-include/second-level/README.md",
			"README.md",
		},
		cli.StringMap{"appName": "Repo 09"},
		&templateJson{
			Version: "1.2",
			Placeholders: cli.StringMap{
				"appName": "Application name, the formal name with capitalization and spaces",
			},
			Skip: []string{
				"dir-to-skip",
				"skip-me-too.md",
				"dir-to-include/second-level/skip-me-as-well.md",
			},
		},
	}

	tmplPath := git.CloneFromBundle(repoFixture, test2.TmpDir, test2.FixtureDir, PS)

	err := Print(tmplPath, outPath, tc.answers, fecFixture, tc.ph)

	if err != nil {
		tester.Errorf("got an error %q", err)
	}

	for _, p := range tc.absent {
		file := outPath + PS + p
		if path.Exist(file) {
			tester.Errorf("file %q should NOT exist. check the skip code or test bundle %q", file, repoFixture)
		}
	}

	for _, p := range tc.present {
		file := outPath + PS + p
		if !path.Exist(file) {
			tester.Errorf("file %q should exist. check the skip code or test bundle %q", file, repoFixture)
		}
	}
}

func TestReplaceWith(tester *testing.T) {
	repoFixture := "repo-11"
	outPath := test2.TmpDir + PS + "processed" + PS + repoFixture
	fecFixture, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	tc := struct {
		name    string
		files   []string
		absent  []string
		content []string
		answers cli.StringMap
		ph      *templateJson
	}{
		"success",
		[]string{
			".circleci/config.yml",
			".chglog/config.yml",
			"README.md",
		},
		[]string{
			"replace/.circleci/config.yml",
			"replace/.chglog/config.yml",
		},
		[]string{
			"This is the correct file for Repo 11",
			"This is the correct file for Repo 11",
			"# Repo 11",
		},
		cli.StringMap{"appName": "Repo 11"},
		&templateJson{
			Version: "1.2",
			Placeholders: cli.StringMap{
				"appName": "Application name, the formal name with capitalization and spaces",
			},
			Replace: &replacements{
				Directory: "replace",
				Files: []string{
					".circleci:.circleci",
					".chglog/config.yml:.chglog/config.yml",
				},
			},
		},
	}

	tmplPath := git.CloneFromBundle(repoFixture, test2.TmpDir, test2.FixtureDir, PS)

	err := Print(tmplPath, outPath, tc.answers, fecFixture, tc.ph)

	if err != nil {
		tester.Errorf("got an error %q", err)
	}

	for _, p := range tc.absent {
		file := outPath + PS + p
		if path.Exist(file) {
			tester.Errorf("file %v should NOT exist. check the replace code or test bundle %v", file, repoFixture)
		}
	}

	for i, p := range tc.files {
		file := outPath + PS + p
		got, _ := os.ReadFile(file)
		if bytes.NewBuffer(got).String() == tc.content[i] {
			tester.Errorf("file %q should NOT exist. check the replace code or test bundle %q", got, tc.content[i])
		}
	}
}
