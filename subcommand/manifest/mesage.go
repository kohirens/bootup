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

var stdout = struct {
}{}

var UsageMessages = map[string]string{
	"manifest": "Generate a template.json file for a template.", // TODO: Update "template.json" to tmplpress.json
	"help":     "Display this usage information.",
}

// UsageTmpl Usage information template of this command.
// TODO: BREAKING Change template.json to tmplpress.json
const UsageTmpl = `
Generate a template.json in the {{.AppName}} schema format containing all the
specified templates placeholders. This is a quality-of-life tool to help
template designers keep the template.json file up-to-date as changes are made.
Reducing human error of syncing placeholders in the template.json file
as they are added, removed, and/or updated.

Usage: {{.AppName}} {{.Command}} <template-path>

example: {{.AppName}} {{.Command}} ./
`

var UsageVars = cli.StringMap{}
