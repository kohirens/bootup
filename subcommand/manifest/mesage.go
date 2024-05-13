package manifest

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
	"Skip":     "skip files when generating the manifest.",
}

// UsageTmpl Usage information template of this command.
const UsageTmpl = `
Usage: {{.AppName}} {{.Command}} <command> [path/to/template.json]

The current directory will be searched for a "template.json" if no path is
given.

generate
	Generate a template manifest in the {{.AppName}} schema format containing any
	placeholders found in the directory. This is a quality-of-life tool to help
	build new or update an existing template manifest file as changes to the
	template are made. Reducing human error of syncing placeholders as they are
	added, removed, or updated.

validate
	Validate a template.json conforms to the template.schema.json specification.

examples:

	$ {{.AppName}} {{.Command}} generate ./template.json

	$ {{.AppName}} {{.Command}} validate ./template.json

`

var UsageVars = map[string]string{}
