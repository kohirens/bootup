package msg

// Messages helpful info to std out
var Messages = struct {
	ActualArgs            string
	CloningToCache        string
	ConfigFileExist       string
	CurrentVersion        string
	CurrentVersionInfo    string
	GitCheckout           string
	MadeNewConfig         string
	NumNonFlagArgs        string
	NumParsedFlags        string
	OutPathExist          string
	OutRepoDir            string
	ProvideValues         string
	PrintAllFlags         string
	PrintFlag             string
	PlaceholderAnswer     string
	PlaceholderAnswerStat string
	PlaceholderHasAnswer  string
	RefInfo               string
	RemoteTagDbug1        string
	ReadConfig            string
	RepoDir               string
	RepoInfo              string
	RunningCommand        string
	SaveData              string
	SkipFile              string
	SubCommands           string
	UnknownFileType       string
	UsageHeader           string
	UsingCache            string
	VerboseLevelInfo      string
}{
	ActualArgs:            "actual arguments passed in: %v",
	CloningToCache:        "no cache; cloning %v to %v",
	ConfigFileExist:       "config file %q exist",
	CurrentVersion:        "%v, %v",
	CurrentVersionInfo:    "version: %v, %v",
	GitCheckout:           "git checkout %s",
	MadeNewConfig:         "saved %d bytes to a new config %q",
	NumNonFlagArgs:        "number of non-flag arguments passed in: %d",
	NumParsedFlags:        "number of parsed flags = %v",
	OutRepoDir:            "repoDir = %v",
	OutPathExist:          "out-path already exits %q",
	ProvideValues:         "note: entering no value will render the placeholder with an empty string",
	PrintAllFlags:         "printing all flags set:",
	PrintFlag:             "\t%s = %v (default= %v)",
	PlaceholderAnswer:     "%v = %q",
	PlaceholderAnswerStat: "please provide values for %v placeholders",
	PlaceholderHasAnswer:  "placeholder %v has a value of %q, so skipping",
	RefInfo:               "ref = %v ",
	RemoteTagDbug1:        "remote tag: %v",
	ReadConfig:            "reading config file %v",
	RepoDir:               "repoDir = %q",
	RepoInfo:              "repo = %q; %q",
	RunningCommand:        "running command %s",
	SaveData:              "%save data: s",
	SkipFile:              "skipping: %v",
	SubCommands:           "sub-commands:\n",
	UsageHeader:           "Usage: %v -[options] [args]\n",
	UsingCache:            "using cache %v",
	UnknownFileType:       "will skip and not process through template engine; could not detect file type for %v",
	VerboseLevelInfo:      "verbose level: %v",
}
