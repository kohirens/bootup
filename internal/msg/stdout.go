package msg

// Stdout helpful info to std out
var Stdout = struct {
	ActualArgs            string
	AddFile               string
	AppCacheDir           string
	AppDataDir            string
	Assignment            string
	CloningToCache        string
	ConfigMethodSetting   string
	CopyAsIs              string
	CurrentVersion        string
	CurrentVersionInfo    string
	Cwd                   string
	GeneratedManifest     string
	MadeNewConfig         string
	NoPlaceholders        string
	NumNonFlagArgs        string
	NumParsedFlags        string
	Parsing               string
	PlaceholderAnswerStat string
	PlaceholderHasAnswer  string
	PrintAllFlags         string
	PrintFlag             string
	Processing            string
	ProvideValues         string
	ReadConfig            string
	RelativeDir           string
	RepoDir               string
	RepoInfo              string
	SaveData              string
	SaveDir               string
	SetValue              string
	Skipping              string
	TemplatePath          string
	TemplatePlaceholders  string
	TemplateVersion       string
	UnknownFileType       string
	UsageHeader           string
	UsingCache            string
	ValuesProvided        string
	VarDefaultValue       string
	VerboseLevelInfo      string
}{
	ActualArgs:            "actual arguments passed in: %v",
	AddFile:               "adding file %v",
	AppCacheDir:           "app cache dir = %v",
	AppDataDir:            "app data dir is %v",
	Assignment:            "%v = %q",
	CloningToCache:        "no cache; cloning %v to %v",
	CopyAsIs:              "file %v will be copied as-is",
	ConfigMethodSetting:   "config.%v(%v)",
	CurrentVersion:        "%v, %v",
	CurrentVersionInfo:    "version: %v, %v",
	Cwd:                   "current working directory is %v",
	GeneratedManifest:     "manifest generated %v",
	MadeNewConfig:         "saved %d bytes to a new config %v",
	NoPlaceholders:        "this template contains no placeholders/actions, which is ok",
	NumNonFlagArgs:        "number of non-flag arguments passed in: %d",
	NumParsedFlags:        "number of parsed flags = %v",
	Parsing:               "parsing %v",
	PlaceholderAnswerStat: "please provide values for %v placeholders",
	PlaceholderHasAnswer:  "placeholder %v has a value of %q, so skipping",
	PrintAllFlags:         "printing all flags set:",
	PrintFlag:             "\t%v = %v (default= %v)",
	Processing:            "processing %v",
	ProvideValues:         "note that entering no value will render the placeholder with an empty string",
	ReadConfig:            "reading config file %v",
	RelativeDir:           "relativePath dir: %v",
	RepoDir:               "repoDir = %q",
	RepoInfo:              "repo = %q; %q",
	SaveData:              "save data: %s",
	SaveDir:               "save dir: %v",
	SetValue:              "%v value = %v",
	Skipping:              "skipping: %v",
	TemplatePath:          "template manifest path: %v",
	TemplatePlaceholders:  "TmplJson.Placeholders = %v",
	TemplateVersion:       "TmplJson.Version = %v",
	UsageHeader:           "Usage: %v -[options] [args]",
	UsingCache:            "using cache located at %v",
	UnknownFileType:       "will skip and not process through template engine; could not detect file type for %v",
	ValuesProvided:        "the following values have been provided",
	VarDefaultValue:       "using default value for placeholder %v",
	VerboseLevelInfo:      "verbose level: %v",
}
