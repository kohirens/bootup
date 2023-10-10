package config

import (
	"github.com/kohirens/stdlib/cli"
)

const UsageTmpl = `
Usage: config set|get <setting> [<value>]

Examples

	config set "CacheDir" "./path/to/a/directory"
	config get "CacheDir"

Settings

	CacheDir - Path to store template downloaded
	ExcludeFileExtensions - Files ending with these extensions will be excluded from parsing and copied as-is

Command

    set Update a configuration value.
    get Returns a configuration value.
`

var Stderr = struct {
	BadExcludeFileExt    string
	CouldNotEncodeConfig string
	ConfigValueNotSet    string
	InvalidConfigMethod  string
}{
	BadExcludeFileExt:    "invalid ExcludeFileExtensions, check format, for example: item1,item2,item3",
	CouldNotEncodeConfig: "could not JSON encode user configuration settings, %v",
	ConfigValueNotSet:    "no value passed in for the setting, try quotes to enter an empty string",
	InvalidConfigMethod:  "invalid config method %v",
}

var UsageMessages = cli.StringMap{
	"config":      "Set or get a configuration value.",
	"config_help": "Print config usage help.",
	"help":        "display this help message.",
}

var UsageVars = cli.StringMap{}
