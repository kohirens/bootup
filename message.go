package main

var errors = struct {
	AnswerFile404    string
	BadTmplType      string
	LocalOutPath     string
	OutPathCollision string
	TmplPath         string
}{}

var stdout = struct {
	OutPathExist string
}{}

var um = map[string]string{
	"answer-path": "Path to a JSON file containing the values for placeholders (which are the keys) defined by a template.",
	"branch":      "Branch of the template to clone when tmplType=git.",
	"default-val": "Used for any unset placeholders and prevents the program waiting for input.",
	"help":        "Prints usage information and exit 0.",
	"out-path":    "Path to output the new project.",
	"tmpl-path":   "URL to a zip or a local path to a directory.",
	"tmpl-type":   "Can be of git|zip.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "Print build version information and exit 0.",
	"config":      "Set or get a configuration value.",
	"manifest":    "Generate a template.json file for a template.",
}
