package cli

import (
	"regexp"
)

// AppData runtime settings shared throughout the application.
type AppData struct {
	AnswersJson  *AnswersJson // data use for template processing
	DataDir      string       // Directory to store app data.
	Path         string       // Path to configuration file.
	SubCmd       string       // sub-command to execute
	Tmpl         string       // Path to template, this will be the cached path.
	TmplLocation string       // Indicates local or remote location to downloaded
	TmplJson     *TmplJson    // Data about the template such as placeholders, their descriptions, version, etc.
}

// getTmplLocation Determine if the template is on the local file system or a remote server.
func getTmplLocation(tmplPath string) string {
	regExpAbsolutePath := regexp.MustCompile(`^/([a-zA-Z._\-][a-zA-Z/._\-].*)?`)
	regExpRelativePath := regexp.MustCompile(`^(\.\.|\.|~)(/[a-zA-Z/._\-].*)?`)
	regExpWinDrive := regexp.MustCompile(`^[a-zA-Z]:\\[a-zA-Z/._\\-].*$`)
	pathType := "remote"

	if regExpAbsolutePath.MatchString(tmplPath) ||
		regExpRelativePath.MatchString(tmplPath) ||
		regExpWinDrive.MatchString(tmplPath) {
		pathType = "local"
	}

	return pathType
}
