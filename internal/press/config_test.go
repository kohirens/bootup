package press

import (
	"bytes"
	"encoding/json"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitConfig(tr *testing.T) {
	testCases := []struct {
		name     string
		filepath string
		cfg      *ConfigSaveData
		wantErr  bool
	}{
		{"file_exist", test.TmpDir + "/TestInitConfig/config-01.json", &ConfigSaveData{}, false},
	}

	for _, tc := range testCases {
		tr.Run(tc.name, func(t *testing.T) {
			_ = os.MkdirAll(filepath.Dir(tc.filepath), dirMode)
			_, err := InitConfig(tc.filepath, "test1")
			if (err != nil) != tc.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !fsio.Exist(tc.filepath) {
				t.Errorf("InitConfig did not save config file %v", tc.filepath)
			}
		})
	}
}

func TestSaveConfig(tr *testing.T) {
	testCases := []struct {
		name     string
		filepath string
		cfg      *ConfigSaveData
		wantErr  bool
	}{
		{
			"file_exist",
			test.TmpDir + "/TestSaveConfig/config-01.json",
			&ConfigSaveData{},
			false,
		},
	}

	for _, tc := range testCases {
		tr.Run(tc.name, func(t *testing.T) {
			_ = os.MkdirAll(filepath.Dir(tc.filepath), dirMode)
			_ = os.WriteFile(tc.filepath, []byte{}, dirMode)

			if err := SaveConfig(tc.filepath, tc.cfg); (err != nil) != tc.wantErr {
				t.Errorf("SaveConfig() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !fsio.Exist(tc.filepath) {
				t.Errorf("SaveConfig did not save config file %v", tc.filepath)
			}
		})
	}
}

func TestLoadConfig(tr *testing.T) {
	testCases := []struct {
		name         string
		filepath     string
		cfg          *ConfigSaveData
		wantCacheDir string
		wantErr      bool
	}{
		{
			"file_exist",
			test.FixturesDir + "/load-config-test-01.json",
			&ConfigSaveData{},
			"/tmp/abc123",
			false,
		},
	}

	for _, tc := range testCases {
		tr.Run(tc.name, func(t *testing.T) {
			got, gotErr := LoadConfig(tc.filepath)

			if (gotErr != nil) != tc.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", gotErr, tc.wantErr)
			}

			if got.CacheDir == "" {
				t.Errorf("got %v, want %v", got.CacheDir, tc.wantCacheDir)
			}
		})
	}
}

func TestLoadUserSettings(tester *testing.T) {
	delayedFunc := test.TmpSetParentDataDir(test.TmpDir + "/TestLoadUserSettings")
	defer delayedFunc()

	var tests = []struct {
		name     string
		filename string
		want     string
		wantErr  bool
	}{
		{
			"good_file",
			test.FixturesDir + PS + "good-config-01.json",
			`"CacheDir":"/tmp/test"`,
			false,
		},
		{
			"bad_file",
			test.FixturesDir + PS + "bad-config-01.json",
			"",
			true,
		},
	}

	for _, tc := range tests {
		tester.Run(tc.name, func(t *testing.T) {
			gotCfg, err := LoadConfig(tc.filename)

			if (err != nil) != tc.wantErr { // test bad values
				t.Errorf("LoadConfig() error %v, want %v", err.Error(), tc.wantErr)
			}

			if !tc.wantErr {
				data, _ := json.Marshal(gotCfg)

				if !strings.Contains(bytes.NewBuffer(data).String(), tc.want) {
					t.Errorf("the config %s did not contain %v", data, tc.want)
				}
			}
		})
	}
}
