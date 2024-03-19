package manifest

var defaultJson = `{
    "$schema": "https://github.com/kohirens/tmplpress/blob/main/template.schema.json",
    "version": "2.2.0",
	"copyAsIs": [
		"*.exe",
		"*.gif",
		"*.jpg",
		"*.mp3",
		"*.pdf",
		"*.png",
		"*.tiff",
		"*.wmv"
	],
	"emptyDirFile": ".empty",
	"skip": [],
	"substitute": "replace"
}
`
