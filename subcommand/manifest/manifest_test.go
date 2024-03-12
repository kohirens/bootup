package manifest

import (
	"encoding/json"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/press"
	"os"
	"reflect"
	"testing"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
)

func TestGenerateATemplateJson(runner *testing.T) {
	fec, _ := stdlib.NewFileExtChecker(&[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"}, &[]string{})

	testCases := []struct {
		name string
		repo string
		want map[string]string
	}{
		{"onlyDataEvaluations", "repo-06", map[string]string{
			"appTitle": "",
			"name":     "",
			"age":      "",
		}},
	}

	for _, tc := range testCases {
		runner.Run(tc.name, func(t *testing.T) {
			repoPath := git.CloneFromBundle(tc.repo, tmpDir, fixtureDir, ps)
			got, err := generateATemplateManifest(repoPath, fec, []string{})
			f := repoPath + ps + press.TmplManifestFile

			if err != nil {
				t.Errorf("want nil, got: %q", err.Error())
			}

			if !path.Exist(f) {
				t.Errorf("no template.json found in %v", repoPath)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		wantErr bool
		want    map[string]string
	}{
		{"case-1", "repo-07", false, map[string]string{"Placeholder1": ""}},
	}

	for _, tt := range tests {
		repoPath := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)
		Init()
		t.Run(tt.name, func(t *testing.T) {
			err := Run([]string{repoPath})
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := loadFile("tmp/repo-07/template.json")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

type mockManifest struct {
	Placeholders map[string]string `json:"placeholders"`
}

func loadFile(filename string) map[string]string {
	content, _ := os.ReadFile(filename)

	ph := &mockManifest{}

	_ = json.Unmarshal(content, ph)

	return ph.Placeholders
}
