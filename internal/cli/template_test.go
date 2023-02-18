package cli

import (
	"github.com/kohirens/tmpltoapp/internal/test"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kohirens/stdlib"
)

type HttpMock struct {
	Resp *http.Response
	Err  error
}

func (h HttpMock) Get(url string) (*http.Response, error) {
	return h.Resp, h.Err
}

func (h HttpMock) Head(url string) (*http.Response, error) {
	return h.Resp, h.Err
}

func TestDownload(runner *testing.T) {
	defer test.Silencer()()

	var err error
	fixtures := HttpMock{
		&http.Response{
			Body:       ioutil.NopCloser(strings.NewReader("200 OK")),
			StatusCode: 200,
		},
		err,
	}

	runner.Run("canDownload", func(t *testing.T) {
		got, err := Download("/fake_dl", TmpDir, &fixtures)
		if err != nil {
			t.Errorf("got %q, want nil", err.Error())
		}
		_, err = os.Stat(got)

		if os.IsNotExist(err) {
			t.Errorf("got %q, want nil", got)
		}
	})
}

func ExampleDownload() {
	client := http.Client{}
	_, err := Download(
		"https://github.com/kohirens/tmpltoapp-test-tpl/archive/main.zip",
		TmpDir,
		&client,
	)

	if err != nil {
		return
	}
}

func TestExtract(runner *testing.T) {
	runner.Run("canExtractDownload", func(t *testing.T) {
		wd, _ := os.Getwd()
		fixture := wd + PS + FixtureDir + PS + "001.zip"
		want := TmpDir + PS + "001"
		_, err := Extract(fixture)

		if err != nil {
			t.Errorf("could not Extract %s, error: %v", want, err.Error())
		}
	})
}

func ExampleExtract() {
	_, err := Extract(
		TmpDir + PS + "001.zip",
	)

	if err != nil {
		return
	}
}

func TestParseDir2(tester *testing.T) {
	defer test.Silencer()()

	fixtures := []struct {
		dstDir,
		name,
		srcDir string
		want func(error) bool
		vars tmplVars
	}{
		{
			TmpDir + PS + "template-04-out",
			"dir1IsEmpty",
			FixtureDir + "/template-04",
			func(e error) bool {
				return !stdlib.PathExist(TmpDir + PS + "template-04-out" + PS + "dir1" + PS + ".empty")
			},
			tmplVars{},
		},
	}

	fileChkr, _ := stdlib.NewFileExtChecker(&[]string{}, &[]string{"tpl"})
	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := ParseDir(fxtr.srcDir, fxtr.dstDir, fxtr.vars, fileChkr, &TmplJson{Excludes: []string{}})
			isAllGood := fxtr.want(err)

			if !isAllGood {
				test.Error("all is not good")
			}
		})
	}
}

func TestParse(tester *testing.T) {
	type wanted func(error) bool

	fixtures := []struct {
		name, tplFile, appDir string
		vars                  tmplVars
		gotWhatIWant          wanted
		failMsg               string
	}{
		{
			"emptyInput",
			"",
			"",
			tmplVars{"var1": "1234"},
			func(err error) bool { return err != nil },
			"failed with no input.",
		},
		{
			"validInput",
			FixtureDir + "/template-02/file-01.tpl",
			TmpDir + "/appDirParse-01",
			tmplVars{"var1": "1234"},
			func(err error) bool {
				f, _ := ioutil.ReadFile(TmpDir + "/appDirParse-01/file-01.tpl")
				s := string(f)
				return s == "testings 1234\n"
			},
			"failed with valid input",
		},
	}

	err := os.MkdirAll(TmpDir+"/appDirParse-01", os.FileMode(DirMode))
	if err != nil {
		tester.Errorf("Could not copy dir, err: %s", err.Error())
	}

	for _, fxtr := range fixtures {
		tester.Run(fxtr.name, func(test *testing.T) {
			err := parse(fxtr.tplFile, fxtr.appDir, fxtr.vars)

			if !fxtr.gotWhatIWant(err) {
				test.Error(fxtr.failMsg)
			}
		})
	}
}

func TestParseDir(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(FixtureDir + PS + "parse-dir-01")
	tmpDir, _ := filepath.Abs(TmpDir)

	fixtures := []struct {
		name, tmplPath, outPath string
		tplVars                 tmplVars
		fileToCheck, want       string
	}{
		{
			"parse-dir-01", fixturePath1, tmpDir + PS + "parse-dir-01",
			tmplVars{"APP_NAME": "SolarPolar"},
			tmpDir + "/parse-dir-01/dir1/README.md", "SolarPolar\n",
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		fec, _ := stdlib.NewFileExtChecker(nil, &[]string{"md", "yml"})

		err := ParseDir(fxtr.tmplPath, fxtr.outPath, fxtr.tplVars, fec, &TmplJson{})

		if err != nil {
			test.Errorf("got an error %q", err.Error())
			return
		}

		got, err := ioutil.ReadFile(tmpDir + "/parse-dir-01/dir1/README.md")

		if err != nil {
			test.Errorf("got an error %q", err.Error())
		}

		//TODO: Verify the files in the parsed dir was processed
		if string(got) != fxtr.want {
			test.Errorf("got %q, but want %q", string(got), fxtr.want)
		}
	})
}

func TestReadTemplateJson(tester *testing.T) {
	fixturePath1, _ := filepath.Abs(FixtureDir + "/template-03")

	fixtures := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			"canBeFound",
			&Config{
				TmplPath: fixturePath1,
			},
			false,
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		got, err := ReadTemplateJson(fxtr.config.TmplPath + PS + TmplManifest)

		if fxtr.shouldErr && err == nil {
			test.Errorf("expected an error, but got nil")
			return
		}

		if !fxtr.shouldErr && err != nil {
			test.Errorf("did not expect an error, but got %s", err.Error())
			return
		}

		if got.Version != "1.0.0" {
			test.Error("could not get version from template.json")
			return
		}
	})
}

func TestPlaceholderInput(tester *testing.T) {
	defer test.Silencer()()
	// Use a temp file to simulate input on the command line.
	tmpFile, err := ioutil.TempFile(TmpDir, "qi-01")
	if err != nil {
		tester.Errorf("failed to make temp file %v", err.Error())
	}

	defer func() {
		if e := tmpFile.Close(); e != nil {
			tester.Errorf("failed to close tmp file: %v", e.Error())
		}
		if e := os.Remove(tmpFile.Name()); e != nil {
			tester.Errorf("failed to remove tmp file: %v", e.Error())
		}
	}()

	if _, e := tmpFile.Write([]byte("1\n")); e != nil {
		tester.Errorf("failed to write content to temp file %v, error: %v", tmpFile.Name(), e.Error())
	}

	resetTmpFile := func() {
		if _, e := tmpFile.Seek(0, 0); e != nil {
			tester.Errorf("failed to reset to beginning of temp file %v, error: %v", tmpFile.Name(), e.Error())
		}
	}

	fixtures := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			"missingAnAnswer",
			&Config{
				AnswersJson: &AnswersJson{
					Placeholders: tmplVars{"var1": "", "var2": ""},
				},
				TmplJson: &TmplJson{
					Version:      "0.1.0",
					Placeholders: tmplVars{"var1": "var1", "var2": "var2", "var3": "var3"},
					Excludes:     nil,
				},
			},
			"var3",
		},
		{
			"noMissingAnswers",
			&Config{
				AnswersJson: &AnswersJson{
					Placeholders: tmplVars{"var1": "1", "var2": "2", "var3": "3"},
				},
				TmplJson: &TmplJson{
					Version:      "0.1.0",
					Placeholders: tmplVars{"var1": "var1", "var2": "var2", "var3": "var3"},
					Excludes:     nil,
				},
			},
			"var3",
		},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(test *testing.T) {
		resetTmpFile()
		err := GetPlaceholderInput(fxtr.config.TmplJson, &fxtr.config.AnswersJson.Placeholders, tmpFile, " ")

		if err != nil {
			test.Errorf("got an error %q", err.Error())
		}

		if fxtr.config.AnswersJson.Placeholders[fxtr.want] != "1" {
			test.Errorf("failed to get a value for placeholder %v using file as input", fxtr.want)
		}
	})
}
