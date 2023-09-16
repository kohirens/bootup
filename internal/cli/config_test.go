package cli

import (
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"strings"
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

func TestGetSettings(t *testing.T) {

	t.Run("configNotFound", func(t *testing.T) {
		cfgFixture := &AppData{}
		gotErr := cfgFixture.LoadUserSettings("does-not-exist")

		if !strings.Contains(gotErr.Error(), "could not open") {
			t.Errorf("got %q, want %q", gotErr, "could not open")
		}
	})

	t.Run("canReadConfig", func(t *testing.T) {
		cfgFixture := &AppData{}
		err := cfgFixture.LoadUserSettings(FixtureDir + PS + "config-01.json")
		if err != nil {
			t.Errorf("got an unexpected error %v", err.Error())
		}
	})
}

func TestInitConfigFile(t *testing.T) {
	var testCases = []struct {
		name string
		want error
		cfg  *AppData
	}{
		{"NotExist", nil, &AppData{Path: TmpDir + PS + "config-fix-01.json"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := tc.cfg.initFile()

			if tc.want != gotErr {
				t.Errorf("got %t; want %t", gotErr, tc.want)
			}

			_, err := os.Stat(tc.cfg.Path)

			if err != nil {
				t.Errorf("got %v; want %v", gotErr, tc.want)
			}
		})
	}
}

func TestLoadAnswers(tester *testing.T) {
	var fixtures = []struct {
		name, file, want string
	}{
		{"goodJson", FixtureDir + PS + "answers-01.json", "value1"},
		{"badJson", FixtureDir + PS + "answers-02.json", ""},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(t *testing.T) {
		got, err := LoadAnswers(fxtr.file)

		if err == nil && got.Placeholders["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got.Placeholders["var1"], fxtr.want)
		}
	})

	fxtr = fixtures[1]
	tester.Run(fxtr.name, func(t *testing.T) {
		_, err := LoadAnswers(fxtr.file)

		if err == nil {
			t.Error("did not get an error")
		}
	})
}

func xTestLoadUserSettings(tester *testing.T) {
	delayedFunc := test.TmpSetParentDataDir(TmpDir)
	defer delayedFunc()

	var tests = []struct {
		name     string
		filename string
		want     *AppData
	}{
		{
			"goodFile",
			FixtureDir + PS + "good-config-01.json",
			&AppData{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
		{
			"badFile",
			FixtureDir + PS + "bad-config-01.json",
			&AppData{
				UsrOpts: &UserOptions{
					ExcludeFileExtensions: &[]string{""},
					CacheDir:              "",
				},
			},
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			gotCfg := &AppData{}
			err := gotCfg.LoadUserSettings(tc.filename)

			if err != nil { // test bad values

			} else { // test good values

			}
		})
	}
}
