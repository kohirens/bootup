package press

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"os"
)

const (
	TmplManifestFile = "template.json" // TODO: BREAKING Rename to tmplpress.json
)

type AnswersJson struct {
	Placeholders cli.StringMap `json:"placeholders"`
}

type TmplManifest struct {
	// A list of files to exclude from processing through the template,
	// but still are output in the final output.
	Excludes []string `json:"excludes"`

	IgnoreExtensions *[]string `json:"ignoreExtensions"`

	// Values to supply to the template to fill in variables.
	Placeholders cli.StringMap `json:"placeholders"`

	// Files that should not be processed through the template engine nor added
	// to the final output.
	Skip []string `json:"skip"`

	// A path to a directory to overwrite other files/directories in the
	// template, before processing output.
	// Note that an empty directory can replace a directory with files.
	Substitute string `json:"substitute"`

	// Optional validation to use when entering placeholder values from the CLI.
	Validation []validator `json:"validation"`

	// The version of the schema to use, which serves only as an indicator of
	// the template engines features.
	Version string `json:"version"`
}

// LoadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func LoadAnswers(filename string) (*AnswersJson, error) {
	if !path.Exist(filename) {
		return nil, fmt.Errorf(msg.Stderr.AnswerFile404, filename)
	}

	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotReadAnswerFile, filename, err.Error())
	}

	var aj *AnswersJson
	if e := json.Unmarshal(content, &aj); e != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotDecodeAnswerFile, filename, e.Error())
	}

	return aj, nil
}

// ReadTemplateJson read variables needed from the template.json file.
func ReadTemplateJson(filePath string) (*TmplManifest, error) {
	log.Dbugf(msg.Stdout.TemplatePath, filePath)

	// Verify the TMPL_MANIFEST file is present.
	if !path.Exist(filePath) {
		return nil, fmt.Errorf(msg.Stderr.TmplManifest404, TmplManifestFile)
	}

	content, e1 := os.ReadFile(filePath)
	if e1 != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotReadFile, filePath, e1)
	}

	q := TmplManifest{}
	if err2 := json.Unmarshal(content, &q); err2 != nil {
		return nil, err2
	}

	log.Dbugf(msg.Stdout.TemplateVersion, q.Version)
	if q.Version == "" {
		return nil, fmt.Errorf(msg.Stderr.MissingTmplJsonVersion)
	}

	log.Dbugf(msg.Stdout.TemplatePlaceholders, len(q.Placeholders))
	if q.Placeholders == nil {
		return nil, fmt.Errorf(msg.Stderr.PlaceholdersProperty)
	}

	return &q, nil
}
