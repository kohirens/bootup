package manifest

import "github.com/kohirens/stdlib/cli"

var stderr = struct {
	EncodingJson         string
	ListWorkingDirectory string
	SavingManifest       string
}{
	EncodingJson:         "could not marshall actions in file %v, error: %v",
	ListWorkingDirectory: "could not get current working directory, %v",
	SavingManifest:       "could not save file %v, error: %v",
}

var UsageMessages = map[string]string{
	"manifest": "Perform operations on the template manifest file.",
	"help":     "Display this usage information.",
}

// UsageTmpl Usage information template of this command.
const UsageTmpl = `
Generate a template manifest in the {{.AppName}} schema format containing any
placeholders found in the directory. This is a quality-of-life tool to help
build new or update an existing template manifest file as changes to the
template are made. Reducing human error of syncing placeholders as they are
added, removed, or updated.

Usage: {{.AppName}} {{.Command}} generate <template-path>

example: {{.AppName}} {{.Command}} generate ./
`

var UsageVars = cli.StringMap{}
