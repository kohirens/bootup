package command

import (
	"github.com/kohirens/stdlib"
	"github.com/kohirens/tmpltoapp/internal"
	"reflect"
	"testing"
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
			repoPath := internal.SetupARepository(tc.repo)
			got, err := GenerateATemplateManifest(repoPath, fec, []string{})
			f := repoPath + internal.PS + "template.json"

			if err != nil {
				t.Errorf("want nil, got: %q", err.Error())
			}

			if !stdlib.PathExist(f) {
				t.Errorf("no template.json found in %v", repoPath)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
