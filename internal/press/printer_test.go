package press

import (
	"bytes"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"path/filepath"
	"testing"
)

func TestFindTemplates(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
		want    int
	}{
		{"1 file", test.FixturesDir + PS, false, 1},
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

func TestParseDir2(tester *testing.T) {
	defer test.Silencer()()

	fixtures := []struct {
		dstDir,
		name,
		srcDir string
		want func(error) bool
		vars cli.StringMap
	}{
		{
			test.TmpDir + PS + "template-04-out",
			"dir1IsEmpty",
			test.FixturesDir + "/template-04",
			func(e error) bool {
				return !path.Exist(test.TmpDir + PS + "template-04-out" + PS + "dir1" + PS + ".empty")
			},
			cli.StringMap{},
		},
	}

	fileChkr, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := Print(fxtr.srcDir, fxtr.dstDir, fxtr.vars, fileChkr, &templateJson{Excludes: []string{}})
			isAllGood := fxtr.want(err)

			if !isAllGood {
				test.Error("all is not good")
			}
		})
	}
}

func TestParseDir(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(test.FixturesDir + PS + "parse-dir-01")
	tmpDir, _ := filepath.Abs(test.TmpDir)
	fixtures := []struct {
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

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		fec, _ := stdlib.NewFileExtChecker(nil, &[]string{"md", "yml"})

		err := Print(fxtr.tmplPath, fxtr.outPath, fxtr.tplVars, fec, &templateJson{})

		if err != nil {
			test.Errorf("got an error %q", err.Error())
			return
		}

		got, err := os.ReadFile(tmpDir + "/parse-dir-01/dir1/README.md")

		if err != nil {
			test.Errorf("got an error %q", err.Error())
		}

		if string(got) != fxtr.want {
			test.Errorf("got %q, but want %q", string(got), fxtr.want)
		}
	})
}

func TestSkipping(tester *testing.T) {
	repoFixture := "repo-10"
	outPath := test.TmpDir + PS + "processed" + PS + repoFixture
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

	tmplPath := test.SetupARepository(repoFixture, test.TmpDir, test.FixturesDir, PS)

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
	outPath := test.TmpDir + PS + "processed" + PS + repoFixture
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

	tmplPath := test.SetupARepository(repoFixture, test.TmpDir, test.FixturesDir, test.PS)

	err := Print(tmplPath, outPath, tc.answers, fecFixture, tc.ph)

	if err != nil {
		tester.Errorf("got an error %q", err)
	}

	for _, p := range tc.absent {
		file := outPath + test.PS + p
		if path.Exist(file) {
			tester.Errorf("file %v should NOT exist. check the replace code or test bundle %v", file, repoFixture)
		}
	}

	for i, p := range tc.files {
		file := outPath + test.PS + p
		got, _ := os.ReadFile(file)
		if bytes.NewBuffer(got).String() == tc.content[i] {
			tester.Errorf("file %q should NOT exist. check the replace code or test bundle %q", got, tc.content[i])
		}
	}
}
