package press

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmplpress/internal/msg"
	"os"
)

const (
	TmplManifestFile = "template.json" // TODO: BREAKING Rename to tmplpress.json
)

type AnswersJson struct {
	Placeholders map[string]string `json:"placeholders"`
}

type TmplManifest struct {
	// A list of files to exclude from processing through the template,
	// but still are output in the final output.
	CopyAsIs []string `json:"copyAsIs,omitempty"`

	// EmptyDirFile Name of a file that marks a directory as empty and has the
	// effect of "mkdir -p". This file allows you to add directories to Git but
	// have them made and empty when the template is pressed.
	EmptyDirFile string `json:"emptyDirFile"`

	// Values to supply to the template to fill in variables.
	Placeholders map[string]string `json:"placeholders,omitempty"`

	// Files that should not be processed through the template engine nor added
	// to the final output.
	Skip []string `json:"skip,omitempty"`

	// A path to a directory to overwrite other files/directories in the
	// template, before processing output.
	// Note that an empty directory can replace a directory with files.
	Substitute string `json:"substitute,omitempty"`

	// Optional validation to use when entering placeholder values from the CLI.
	Validation []validator `json:"validation,omitempty"`

	// The version of the schema to use, which serves only as an indicator of
	// the template engines features.
	Version string `json:"version"`
}

// LoadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func LoadAnswers(filename string) (*AnswersJson, error) {
	if !fsio.Exist(filename) {
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
	if !fsio.Exist(filePath) {
		return nil, fmt.Errorf(msg.Stderr.TmplManifest404, TmplManifestFile)
	}

	content, e1 := os.ReadFile(filePath)
	if e1 != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotReadFile, filePath, e1)
	}

	q, e2 := NewTmplManifest(content)
	if e2 != nil {
		return nil, e2
	}

	log.Dbugf(msg.Stdout.TemplateVersion, q.Version)
	if q.Version == "" {
		return nil, fmt.Errorf(msg.Stderr.MissingTmplJsonVersion)
	}

	// It is possible to have a template with no placeholders.
	log.Dbugf(msg.Stdout.TemplatePlaceholders, len(q.Placeholders))

	return q, nil
}

func NewTmplManifest(content []byte) (*TmplManifest, error) {
	tmf := &TmplManifest{}
	if e := json.Unmarshal(content, &tmf); e != nil {
		return nil, fmt.Errorf(msg.Stderr.NewManifest, e.Error())
	}

	return tmf, nil
}
