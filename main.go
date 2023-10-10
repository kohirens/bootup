package main

//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName flags

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	stdc "github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"github.com/kohirens/tmpltoapp/internal/press"
	"github.com/kohirens/tmpltoapp/subcommand/config"
	"github.com/kohirens/tmpltoapp/subcommand/manifest"
	"os"
	"regexp"
)

// TODO: Change name to tmplpress

const (
	AppName    = "tmpltoapp"
	gitConfDir = ".git"
	Summary    = "Generate an application from a template."
)

var (
	// appConfig Runtime settings used throughout the application.
	appConfig = &cli.AppData{}

	appData *press.AppData

	flags = &appFlags{
		subcommands: map[string]*flag.FlagSet{},
	}
)

func init() {
	// Define all flags
	defineFlags(flags)

	usg := stdc.NewUsage(AppName, um, nil, Summary, usageTmpl2)
	usg.Command.AddCommand(
		config.Init(),
		config.Name,
		config.UsageMessages,
		config.UsageVars,
		config.Summary,
		config.UsageTmpl,
	)
	usg.Command.AddCommand(
		manifest.Init(),
		manifest.Name,
		manifest.UsageMessages,
		manifest.UsageVars,
		manifest.Summary,
		manifest.UsageTmpl,
	)
}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			logf(msg.Stderr.FatalHeader)
			log.Fatf(mainErr.Error())
		}
		os.Exit(0)
	}()

	mainErr = parseCli(flags)
	if mainErr != nil {
		return
	}

	// Exit if we are just printing help usage
	if flags.Help {
		flag.Usage()
		fmt.Println()
		return
	}

	if flags.Version {
		log.Logf(msg.Stdout.CurrentVersion, flags.CurrentVersion, flags.CommitHash)
		os.Exit(0)
	}

	ap, err1 := press.BuildAppDataPath(AppName)
	if err1 != nil {
		mainErr = err1
		return
	}

	cf := ap + cli.PS + press.ConfigFileName
	var sc *press.ConfigSaveData

	if path.Exist(cf) {
		sc, mainErr = press.LoadConfig(cf)
	} else {
		// Make a configuration file when there is none.
		sc, mainErr = press.InitConfig(cf, AppName)
	}
	if mainErr != nil {
		return
	}

	appData, mainErr = press.NewAppData(sc)
	if mainErr != nil {
		return
	}

	ca := flag.Args()

	if len(ca) > 0 {
		switch ca[0] {
		case config.Name:
			// store or get the key and return
			mainErr = config.Run(ca[1:], AppName)
			return
		case manifest.Name:
			if e := manifest.ParseFlags(ca[1:]); e != nil {
				mainErr = e
				return
			}
			ma := &manifest.Arguments{}

			// TODO: BREAKING Add this to the template.json, the template designer should be responsible for this; ".empty" should still be embedded in this app though.
			fec, _ := stdlib.NewFileExtChecker(&[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"}, &[]string{})
			_, mainErr = manifest.GenerateATemplateManifest(ma.Path, fec, []string{})
			return
		}
	}

	if e := parseMainArgs(flags, ca); e != nil {
		mainErr = e
		return
	}

	var tmplToPress string

	if flags.TmplType == "git" {
		var repo, commitHash string
		var err2 error

		if flags.Branch == "latest" {
			latestTag, e3 := git.LatestTag(flags.TmplPath)
			log.Infof(e3.Error())
			if latestTag != "" {
				flags.Branch = latestTag
			}
		}

		// Determine the cache location
		repoDir := appData.CacheDir + cli.PS + getRepoDir(flags.TmplPath, flags.Branch)
		log.Infof(msg.Stdout.RepoDir, repoDir)

		// Do a pull when the repo already exists.
		if path.DirExist(repoDir + cli.PS + gitConfDir) {
			infof(msg.Stdout.UsingCache, repoDir)
			repo, commitHash, err2 = git.Checkout(repoDir, flags.Branch)
		} else {
			infof(msg.Stdout.CloningToCache, repoDir, flags.TmplPath)
			repo, commitHash, err2 = git.Clone(flags.TmplPath, repoDir, flags.Branch)
		}

		infof(msg.Stdout.RepoInfo, repo, commitHash)
		if err2 != nil {
			mainErr = err2
			return
		}
		tmplToPress = repo
	}

	if !path.DirExist(tmplToPress) {
		mainErr = fmt.Errorf(msg.Stderr.InvalidTmplDir, tmplToPress)
		return
	}

	fec, err1 := stdlib.NewFileExtChecker(appData.ExcludeFileExtensions, &[]string{})
	if err1 != nil {
		mainErr = fmt.Errorf(msg.Stderr.CannotInitFileChecker, err1.Error())
	}

	// Require template directories to have a specific file in order to be processed to prevent processing directories unintentionally.
	tmplManifestFile := tmplToPress + cli.PS + press.TmplManifestFile
	tmplManifest, errX := cli.ReadTemplateJson(tmplManifestFile)
	if errX != nil {
		mainErr = fmt.Errorf(msg.Stderr.MissingTmplJson, press.TmplManifestFile, tmplManifestFile, errX.Error())
		return
	}

	appConfig.TmplJson = tmplManifest
	appConfig.AnswersJson = cli.NewAnswerJson()

	if path.Exist(flags.AnswersPath) {
		appConfig.AnswersJson, mainErr = press.LoadAnswers(flags.AnswersPath)
		if mainErr != nil {
			return
		}
	}

	// Checks for any missing placeholder values waits for their input from the CLI.
	if e := cli.GetPlaceholderInput(appConfig.TmplJson, &appConfig.AnswersJson.Placeholders, os.Stdin, flags.DefaultVal); e != nil {
		mainErr = fmt.Errorf(msg.Stderr.GettingAnswers, e.Error())
	}

	cli.ShowAllPlaceholderValues(appConfig.TmplJson, &appConfig.AnswersJson.Placeholders)

	mainErr = cli.Press(tmplToPress, flags.OutPath, appConfig.AnswersJson.Placeholders, fec, appConfig.TmplJson)
}

// / TODO: Move this to parseCli
func parseMainArgs(af *appFlags, pArgs []string) error {
	// throw an error when a flag comes after any arguments.
	for i := 0; i < len(pArgs); i++ {
		v := pArgs[i]
		if v[0] == '-' {
			return fmt.Errorf(msg.Stderr.FlagOrderErr, v)
		}
	}

	numArgs := len(pArgs)
	if numArgs >= 1 {
		af.TmplPath = pArgs[0]
	}
	if numArgs >= 2 {
		af.OutPath = pArgs[1]
	}
	if numArgs >= 3 {
		af.AnswersPath = pArgs[3]
	}

	if e := validateMainArgs(af); e != nil {
		return e
	}
	return nil
}

// validateMainArgs parses command line flags into program options.
func validateMainArgs(af *appFlags) error {
	if af.TmplPath == "" {
		return fmt.Errorf(errors.TmplPath)
	}

	if af.OutPath == "" {
		return fmt.Errorf(errors.LocalOutPath)
	}

	if af.TmplPath == af.OutPath {
		return fmt.Errorf(errors.OutPathCollision, af.TmplPath, af.OutPath)
	}

	if path.DirExist(af.OutPath) {
		return fmt.Errorf(stdout.OutPathExist, af.OutPath)
	}

	if af.AnswersPath != "" && !path.Exist(af.AnswersPath) {
		return fmt.Errorf(errors.AnswerFile404, af.AnswersPath)
	}

	regExpTmplType := regexp.MustCompile("^(zip|git|dir)$")

	if !regExpTmplType.MatchString(af.TmplType) {
		return fmt.Errorf(errors.BadTmplType, af.TmplType)
	}

	return nil
}
