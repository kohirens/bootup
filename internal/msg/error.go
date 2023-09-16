package msg

var Errors = struct {
	AnswerFile404          string
	AppDataDir             string
	BadExcludeFileExt      string
	BadTmplType            string
	CannotDecodeAnswerFile string
	CannotInitFileChecker  string
	CannotReadAnswerFile   string
	Checkout               string
	Cloning                string
	CouldNot               string
	CouldNotCloseFile      string
	CouldNotDecode         string
	CouldNotEncodeConfig   string
	CouldNotMakeCacheDir   string
	CouldNotSaveConf       string
	CouldNotWriteFile      string
	CurrentBranch          string
	FatalHeader            string
	FlagOrderErr           string
	FileTooBig             string
	GettingAnswers         string
	GettingCommitHash      string
	GitCheckoutFailed      string
	GitFetchFailed         string
	GitExitErrCode         string
	GetLatestTag           string
	GetRemoteTags          string
	InvalidNoArgs          string
	InvalidNoSubCmdArgs    string
	InvalidTmplDir         string
	LocalOutPath           string
	MissingTmplJson        string
	NoGitTagFound          string
	NoInput                string
	OutPathCollision       string
	ParsingConfigArgs      string
	PathNotAllowed         string
	RunGitFailed           string
	TmplManifest404        string
	TmplOutput             string
	TmplPath               string
	UnhandledHttpErr       string
	ParsingFile            string
	PathNotExist           string
}{
	AnswerFile404:          "could not find the answer file %q, please specify a path to a valid answer file that exist: given %q",
	AppDataDir:             "the following error occurred trying to get the app data directory: %q",
	BadExcludeFileExt:      "invalid ExcludeFileExtensions, check format, for example: item1,item2,item3",
	BadTmplType:            "%q is an invalid value for flag tmplType, or it was not set, must be zip|git",
	CannotDecodeAnswerFile: "could not decode JSON in answer file %q, because of: %s",
	CannotInitFileChecker:  "cannot instantiate file extension checker: %v",
	CannotReadAnswerFile:   "there was an error reading the answer file %q: %s",
	Checkout:               "checkout failed for branch %q",
	Cloning:                "error cloning %v: %s",
	CouldNot:               "could not %s",
	CouldNotCloseFile:      "could not close file %v, %v",
	CouldNotDecode:         "could not decode %q, error: %s",
	CouldNotEncodeConfig:   "could not JSON encode user configuration settings, %v",
	CouldNotMakeCacheDir:   "could not make cache directory, error: %s",
	CouldNotSaveConf:       "could not save a config file, reason: %v",
	CouldNotWriteFile:      "could not write file %v, reason: %v",
	CurrentBranch:          "failed to get current for %s",
	FatalHeader:            "\nfatal error detected: ",
	FlagOrderErr:           "flag %v MUST come before any non-flag arguments, a fix would be to move this flag to the left of other input arguments",
	FileTooBig:             "template file too big to parse, must be less thatn %v bytes",
	GettingAnswers:         "problem getting answers; error %q",
	GettingCommitHash:      "error getting commit hash %v: %s",
	GitCheckoutFailed:      "git checkout failed: %s",
	GetLatestTag:           "failed to get latest tag from %v: %v",
	GetRemoteTags:          "could not get remote tags, please check for a typo, it exist, and is readable: %v",
	GitExitErrCode:         "git %v returned exit code %q",
	GitFetchFailed:         "fetch failed on %s and %s; %s",
	InvalidNoArgs:          "invalid number of arguments passed to the config command, please see config -help for usage",
	InvalidNoSubCmdArgs:    "subcommand %v takes at least %v arguments, run \"%[1]s -h\" for usage details",
	InvalidTmplDir:         "invalid template directory %q",
	LocalOutPath:           "enter a local path to output the app",
	MissingTmplJson:        "%s is a file that is required to be in the template, there was a problem reading %q; error %q",
	NoGitTagFound:          "no tag found in %v",
	NoInput:                "no input",
	OutPathCollision:       "-tmpl-path %q and -out-path %q cannot point to the same directory",
	ParsingConfigArgs:      "error parsing config command args: %v",
	PathNotAllowed:         "path/URL to template is not in the allow-list",
	RunGitFailed:           "error running git %v: %v\n%s",
	TmplManifest404:        "the required manifest %q file was not found",
	TmplOutput:             "template has NOT been cloned locally",
	TmplPath:               "please specify a path (or URL) to a template",
	UnhandledHttpErr:       "template Download aborted; I'm coded to NOT do anything when HTTP status is %q and status code is %d",
	ParsingFile:            "could not parse file %v, error: %v",
	PathNotExist:           "could not locate the path %q",
}
