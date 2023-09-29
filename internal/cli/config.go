package cli

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

// TODO: Change to RunTimeSettings to avoid confution with config sub command

// AppData runtime settings shared throughout the application.
type AppData struct {
	AnswersJson  *AnswersJson // data use for template processing
	DataDir      string       // Directory to store app data.
	Path         string       // Path to configuration file.
	SubCmd       string       // sub-command to execute
	Tmpl         string       // Path to template, this will be the cached path.
	TmplLocation string       // Indicates local or remote location to downloaded
	TmplJson     *TmplJson    // Data about the template such as placeholders, their descriptions, version, etc.
	UsrOpts      *UserOptions // options that can configured by the user.
}

// UserOptions Options the user can set
type UserOptions struct {
	ExcludeFileExtensions *[]string
	CacheDir              string
}

// Setup All application configuration.
func (cfg *AppData) Setup(appName, ps, tmplType, tmplPath string, dirMode os.FileMode) error {
	osDataDir, err1 := stdlib.AppDataDir() //os.UserHomeDir()
	log.Dbugf("app data dir = %q\n", osDataDir)
	if err1 != nil {
		return err1
	}

	// Make a hidden directory in userspace to store data.
	cfg.DataDir = osDataDir + ps + "." + appName
	if e := os.MkdirAll(cfg.DataDir, dirMode); e != nil {
		return e
	}

	cfg.UsrOpts.CacheDir = cfg.DataDir + ps + "cache"
	if e := os.MkdirAll(cfg.UsrOpts.CacheDir, dirMode); e != nil {
		return fmt.Errorf(msg.Stderr.CouldNotMakeCacheDir, e.Error())
	}

	cfg.Path = cfg.DataDir + ps + "config.json"
	// Make a configuration file when there is none.
	if e := cfg.initFile(); e != nil {
		return e
	}

	if e := cfg.LoadUserSettings(cfg.Path); e != nil {
		return e
	}

	// Determine if the template is on the local file system or a remote server.
	cfg.TmplLocation = getTmplLocation(tmplPath)

	if tmplType == "dir" { // TODO: Auto detect if the template is a git repo (look for .git), a zip (look for .zip), or dir (assume dir)
		cfg.Tmpl = filepath.Clean(tmplPath)
	}

	return nil
}

// Initialize a configuration file.
func (cfg *AppData) initFile() error {
	if stdlib.PathExist(cfg.Path) {
		log.Infof(msg.Stdout.ConfigFileExist, cfg.Path)
		return nil
	}

	f, err1 := os.Create(cfg.Path)
	if err1 != nil {
		return fmt.Errorf(msg.Stderr.CouldNotSaveConf, err1.Error())
	}

	data, err2 := json.Marshal(cfg.UsrOpts)
	if err2 != nil {
		return fmt.Errorf(msg.Stderr.CouldNotEncodeConfig, err2.Error())
	}

	b, err3 := f.Write(data)
	if err3 != nil {
		return fmt.Errorf(msg.Stderr.CouldNotWriteFile, cfg.Path, err3.Error())
	}

	if e := f.Close(); e != nil {
		return fmt.Errorf(msg.Stderr.CouldNotCloseFile, cfg.Path, e.Error())
	}

	log.Infof(msg.Stdout.MadeNewConfig, b, cfg.Path)

	return nil
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

// LoadUserSettings from a file, replacing the default built-in settings.
func (cfg *AppData) LoadUserSettings(filename string) error {
	log.Infof(msg.Stdout.ReadConfig, filename)
	content, er := os.ReadFile(filename)

	if os.IsNotExist(er) {
		return fmt.Errorf(msg.Stderr.CouldNot, er.Error())
	}

	if e := json.Unmarshal(content, &cfg.UsrOpts); e != nil {
		return fmt.Errorf(msg.Stderr.CouldNotDecode, filename, er.Error())
	}

	return nil
}

// LoadAnswers Load key/value pairs from a JSON file to fill in placeholders (provides that data for the Go templates).
func LoadAnswers(filename string) (*AnswersJson, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotReadAnswerFile, filename, err.Error())

	}

	var aj *AnswersJson
	if e := json.Unmarshal(content, &aj); e != nil {
		return nil, fmt.Errorf(msg.Stderr.CannotDecodeAnswerFile, filename, e.Error())
	}

	return aj, nil
}
