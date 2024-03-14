package msg

var Stderr = struct {
	AnswerFile404          string
	AppDataDir             string
	BadExcludeFileExt      string
	CannotCopyDirToDir     string
	CannotDecodeAnswerFile string
	CannotInitFileChecker  string
	CannotReadFile         string
	CannotReadAnswerFile   string
	CannotRemoveDir        string
	CouldNot               string
	CouldNotCloseFile      string
	CouldNotDecode         string
	CouldNotEncodeConfig   string
	CouldNotMakeCacheDir   string
	CouldNotSaveConf       string
	CouldNotWriteFile      string
	FatalHeader            string
	FlagOrderErr           string
	FileTooBig             string
	GettingAnswers         string
	GitFetchFailed         string
	GetLatestTag           string
	GetRemoteTags          string
	InvalidNoArgs          string
	InvalidNoSubCmdArgs    string
	InvalidRegExp          string
	InvalidTmplDir         string
	MissingTmplJson        string
	MissingTmplJsonVersion string
	NoGitTagFound          string
	NoInput                string
	NoPath                 string
	NoSetting              string
	ParseBool              string
	ParseInt               string
	ParseUInt              string
	ParsingConfigArgs      string
	PathNotAllowed         string
	PlaceholdersProperty   string
	RunGitFailed           string
	TmplManifest404        string
	TmplOutput             string
	UnhandledHttpErr       string
	ParsingFile            string
	PathNotExist           string
}{
	AnswerFile404:          "could not find the answer file, please specify a path to a valid answer file that exist: given %q",
	AppDataDir:             "the following error occurred trying to get the app data directory: %q",
	BadExcludeFileExt:      "invalid ExcludeFileExtensions, check format, for example: item1,item2,item3",
	CannotCopyDirToDir:     "could not copy %v to %v: %v",
	CannotDecodeAnswerFile: "could not decode JSON in answer file %q, because of: %s",
	CannotInitFileChecker:  "cannot instantiate file extension checker: %v",
	CannotReadAnswerFile:   "there was an error reading the answer file %q: %s",
	CannotReadFile:         "could not read file %v: %v",
	CannotRemoveDir:        "could not remove dir %v: %v",
	CouldNot:               "could not %s",
	CouldNotCloseFile:      "could not close file %v, %v",
	CouldNotDecode:         "could not decode %q, error: %s",
	CouldNotEncodeConfig:   "could not JSON encode user configuration settings, %v",
	CouldNotMakeCacheDir:   "could not make cache directory, error: %s",
	CouldNotSaveConf:       "could not save a config file, reason: %v",
	CouldNotWriteFile:      "could not write file %v, reason: %v",
	FatalHeader:            "\nfatal error detected: ",
	FlagOrderErr:           "flag %v MUST come before any non-flag arguments, a fix would be to move this flag to the left of other input arguments",
	FileTooBig:             "template file too big to parse, must be less thatn %v bytes",
	GettingAnswers:         "problem getting answers; error %q",
	GetLatestTag:           "failed to get latest tag from %v: %v",
	InvalidNoArgs:          "invalid number of arguments passed to the config command, please see config -help for usage",
	InvalidNoSubCmdArgs:    "subcommand %v takes at least %v arguments, run \"%[1]s -h\" for usage details",
	InvalidRegExp:          "invalid regular expression \"%v\", %v",
	InvalidTmplDir:         "invalid template directory %q",
	MissingTmplJson:        "%s is a file that is required to be in the template, there was a problem reading %q; error %q",
	MissingTmplJsonVersion: "missing the Version property in template.json",
	NoGitTagFound:          "no tag found in %v",
	NoInput:                "no input",
	NoPath:                 "unable to determine absolute path for %v, because %v",
	NoSetting:              "no setting %v found",
	ParseBool:              "%v is not a valid boolean value",
	ParseInt:               "could not parse %v as a integer, %v",
	ParseUInt:              "could not parse %v as a natural number, %v",
	ParsingConfigArgs:      "error parsing config command args: %v",
	PathNotAllowed:         "path/URL to template is not in the allow-list",
	PlaceholdersProperty:   "missing the placeholders property in template.json",
	TmplManifest404:        "the required manifest %q file was not found",
	TmplOutput:             "template has NOT been cloned locally",
	UnhandledHttpErr:       "template Download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
	ParsingFile:            "could not parse file %v, error: %v",
	PathNotExist:           "could not locate the path %v",
}
