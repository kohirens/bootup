package press

import (
	"github.com/kohirens/stdlib/path"
	"strings"
)

func inSkipArray(p string, skips []string) bool {
	found := false

	for _, skip := range skips {
		skip = path.Normalize(skip)

		if strings.Contains(p, skip) {
			found = true
			break
		}
	}

	return found
}
