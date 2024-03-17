package manifest

var defaultJson = `{
    "$schema": "https://github.com/kohirens/tmpltoapp/blob/main/template.schema.json",
    "version": "2.2.0",
	"excludes": [
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
	"skip": []
}
`
