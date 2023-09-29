package main

var usageTmpl2 = `
{{- define "optionHeader"}}
Options:
{{end -}}

{{- define "option"}}
{{printf "  -%-11s %v" .OptionName .OptionInfo}}{{with .DefaultValue }} (default = {{.}}){{end}}
{{end -}}

{{- define "commandHeader"}}
Commands
{{end -}}

{{- define "subcommand"}}
  {{.Command}} - {{.Summary}}

    usage: {{.AppName}} [global options] {{.Command}} [options] <args>

    See {{.AppName}} {{.Command}} -help
{{end}}
Usage: {{.AppName}} -[options] <args>

example: {{.AppName}} [options] <args>
`
