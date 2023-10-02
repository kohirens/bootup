package press

import (
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"os"
)

const (
	FilePerm       = 0774
	ConfigFileName = "config.json"
)

// AppData runtime settings shared throughout the application.
type AppData struct {
	AnswersJson           *AnswersJson // data use for template processing
	CacheDir              string       // Directory to store app data.
	DataDir               string       // Directory to store app data.
	ExcludeFileExtensions *[]string
	Path                  string // Path to configuration file.
	SubCmd                string // sub-command to execute
	Tmpl                  string // Path to template, this will be the cached path.
	TmplLocation          string // Indicates local or remote location to downloaded
	//	TmplJson              *TmplJson // Data about the template such as placeholders, their descriptions, version, etc.
}

type ConfigSaveData struct {
	CacheDir              string
	ExcludeFileExtensions *[]string
}

func NewAppData(sc *ConfigSaveData) (*AppData, error) {
	// Override defaults with user settings.
	ad := &AppData{
		CacheDir:              sc.CacheDir,
		ExcludeFileExtensions: sc.ExcludeFileExtensions,
	}

	return ad, nil
}

// InitConfig Save configuration file when it does not exist.
func InitConfig(filepath string) (*ConfigSaveData, error) {
	cd, err1 := buildCacheDirPath()
	if err1 != nil {
		return nil, fmt.Errorf(msg.Stderr.CouldNotMakeCacheDir, err1.Error())
	}

	sc := &ConfigSaveData{
		CacheDir:              cd,
		ExcludeFileExtensions: &[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"},
	}

	f, err2 := os.Create(filepath)
	if err2 != nil {
		return nil, fmt.Errorf(msg.Stderr.CouldNotSaveConf, err2.Error())
	}

	defer func() {
		if e := f.Close(); e != nil {
			panic(fmt.Errorf(msg.Stderr.CouldNotCloseFile, filepath, e.Error()))
		}
	}()

	data, err3 := json.MarshalIndent(sc, "", "\t")
	if err3 != nil {
		return nil, fmt.Errorf(msg.Stderr.CouldNotEncodeConfig, err3.Error())
	}

	b, err4 := f.Write(data)
	if err4 != nil {
		return nil, fmt.Errorf(msg.Stderr.CouldNotWriteFile, filepath, err4.Error())
	}

	log.Infof(msg.Stdout.MadeNewConfig, b, filepath)

	return sc, nil
}

// LoadConfig Load to the configuration file.
func LoadConfig(filepath string) (*ConfigSaveData, error) {
	log.Infof(msg.Stdout.ReadConfig, filepath)

	data, err1 := os.ReadFile(filepath)
	if os.IsNotExist(err1) {
		return nil, fmt.Errorf(msg.Stderr.CouldNot, err1.Error())
	}

	csd := &ConfigSaveData{}
	if e := json.Unmarshal(data, &csd); e != nil {
		return nil, fmt.Errorf(msg.Stderr.CouldNotDecode, filepath, e.Error())
	}

	return csd, nil
}

// SaveConfig Save any updates to the configuration file.
func SaveConfig(filepath string, cfg *ConfigSaveData) error {
	data, err1 := json.MarshalIndent(cfg, "", "\t")

	log.Dbugf(msg.Stdout.SaveData, data)

	if err1 != nil {
		return fmt.Errorf(msg.Stderr.CouldNotEncodeConfig, err1.Error())
	}

	if e := os.WriteFile(filepath, data, FilePerm); e != nil {
		return e
	}

	return nil
}

func BuildAppDataPath(appName string) (string, error) {
	osDataDir, err1 := os.UserConfigDir()
	if err1 != nil {
		return "", err1
	}

	// Make an app directory in user config space to store config data.
	appDataDir := osDataDir + PS + appName
	if e := os.MkdirAll(appDataDir, dirMode); e != nil {
		return "", e
	}

	log.Dbugf("app data dir = %q", appDataDir)

	return appDataDir, nil
}

func buildCacheDirPath() (string, error) {
	osCacheDir, err1 := os.UserCacheDir()
	if err1 != nil {
		return "", err1
	}

	// Make a directory in user cache space to download templates.
	appCacheDir := osCacheDir + PS + "cache"
	if e := os.MkdirAll(appCacheDir, dirMode); e != nil {
		return "", fmt.Errorf(msg.Stderr.CouldNotMakeCacheDir, e.Error())
	}

	log.Dbugf("app cache dir = %q\n", osCacheDir)

	return appCacheDir, nil
}
