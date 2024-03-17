package manifest

var TmplJsonTmpl = `{
    "$schema": "https://github.com/kohirens/tmpltoapp/blob/main/template.schema.json",
    "version": "2.1.0",
	"placeholders": {{printf "%s" .Placeholders}},
	"excludes": {},
	"skip": {}
}
`
