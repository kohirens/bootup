package press

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"os"
)

const (
	TmplManifestFile = "template.json" // TODO: BREAKING Rename to tmplpress.json
)

// LoadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func LoadAnswers(filename string) (*cli.AnswersJson, error) {
	if !path.Exist(filename) {
		return nil, fmt.Errorf(msg.Stderr.AnswerFile404, filename)
	}

	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotReadAnswerFile, filename, err.Error())
	}

	var aj *cli.AnswersJson
	if e := json.Unmarshal(content, &aj); e != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotDecodeAnswerFile, filename, e.Error())
	}

	return aj, nil
}
