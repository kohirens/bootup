package main

var errors = struct {
	AnswerFile404    string
	BadTmplType      string
	LocalOutPath     string
	OutPathCollision string
	Path404          string
	TmplPath         string
}{
	AnswerFile404:    "could not find the answer file, please specify a path to a valid answer file that exist: given %q",
	BadTmplType:      "%q is an invalid value for flag tmplType, or it was not set, must be git",
	LocalOutPath:     "enter a local path to output the app",
	OutPathCollision: "invalid input; the template path and out path point to the same directory:\ntmpl path = %v\n out path = %v",
	Path404:          "problem with the path %v, please check the path exist and is readable: %v",
	TmplPath:         "please specify a path (or URL) to a template",
}

var stdout = struct {
	OutPathExist string
}{
	OutPathExist: "out-path already exits %q",
}

var um = map[string]string{
	"answer-path": "Path to a JSON file containing the values for placeholders (which are the keys) defined by a template.",
	"branch":      "Branch of the template to clone when tmplType=git.",
	"default-val": "Used for any unset placeholders and prevents the program waiting for input.",
	"help":        "Prints usage information and exit 0.",
	"out-path":    "Path to output the new project.",
	"tmpl-path":   "URL to a git repository or a local path to a directory.",
	"tmpl-type":   "Can be of git.",
	"verbosity":   "Set the level of information printed when running.",
	"version":     "Print build version information and exit 0.",
	"config":      "Set or get a configuration value.",
	"manifest":    "Generate a template.json file for a template.",
}
