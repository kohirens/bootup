package press

import (
	"github.com/kohirens/stdlib/fsio"
	"github.com/ryanuber/go-glob"
)

// Get to know glob before you question the correctness of my program:
// https://man7.org/linux/man-pages/man7/glob.7.html
func inSkipArray(pathToFile string, skips []string) bool {
	skip := false

	for _, pattern := range skips {
		pattern = fsio.Normalize(pattern)

		if glob.Glob(pattern, pathToFile) {
			skip = true
			break
		}
	}

	return skip
}
