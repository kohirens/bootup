package press

import (
	"github.com/kohirens/stdlib/path"
	"strings"
)

type replacements struct {
	Directory string   `json:"directory"`
	Files     []string `json:"files"`
}

// replaceWith Replace the current file with another.
func replaceWith(cf, ps, sp, tmplRoot string, replace *replacements) string {
	if replace == nil {
		return sp
	}

	for _, fileMap := range replace.Files {
		fileMap = path.Normalize(fileMap)
		fileAry := strings.Split(fileMap, ":")

		if strings.Contains(cf, fileAry[1]) {
			if cf == fileAry[1] {
				return tmplRoot + ps + replace.Directory + ps + fileAry[0]
			}
			// match by prefix

			return tmplRoot + ps + replace.Directory + ps + fileAry[0] + strings.Replace(cf, fileAry[0], "", 1)
		}
	}

	return sp
}
