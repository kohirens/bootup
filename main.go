package main

//go:generate git-tool-belt semver -save info.go -format go -packageName main -varName appConfig

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"log"
	"net/http"
	"os"
)

// TODO: Change name to tmplpress

const (
	AppName    = "tmpltoapp"
	gitConfDir = ".git"
)

var (
	// appConfig Runtime settings used throughout the application.
	appConfig = &cli.Config{
		UsrOpts: &cli.UserOptions{
			ExcludeFileExtensions: &[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"},
		},
	}
)

func init() {
	// Define all flags
	defineFlags(appConfig)
}

func main() {
	var mainErr error

	defer func() {
		if mainErr != nil {
			logf(cli.Errors.FatalHeader)
			log.Fatalln(mainErr.Error())
		}
		os.Exit(0)
	}()

	mainErr = parseFlags(appConfig)
	if mainErr != nil {
		return
	}

	mainErr = appConfig.Setup(AppName, cli.PS, cli.DirMode)
	if mainErr != nil {
		return
	}

	// Exit if we are just printing help usage
	if appConfig.Help {
		mainErr = Usage(appConfig)
		return
	}

	infof(cli.Messages.CurrentVersionInfo, appConfig.CurrentVersion, appConfig.CommitHash)

	// process sub-commands
	switch appConfig.SubCmd {
	case cli.CmdConfig:
		// store or get the key and return
		mainErr = cli.UpdateUserSettings(appConfig, cli.DirMode)
		return
	case cli.CmdManifest:
		// store or get the key and return
		fec, _ := stdlib.NewFileExtChecker(&[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"}, &[]string{})
		_, mainErr = cli.GenerateATemplateManifest(appConfig.SubCmdManifest.Path, fec, []string{})
		return
	}

	if appConfig.TmplType == "zip" {
		var zipFile string
		var iErr error
		zipFile = appConfig.TmplPath
		if appConfig.TmplLocation == "remote" {
			client := http.Client{}
			zipFile, iErr = cli.Download(appConfig.TmplPath, appConfig.UsrOpts.CacheDir, &client)
			if iErr != nil {
				mainErr = iErr
				return
			}
		}

		appConfig.Tmpl, iErr = cli.Extract(zipFile)
		if iErr != nil {
			mainErr = iErr
			return
		}
	}

	if appConfig.TmplType == "git" {
		var repo, commitHash string
		var err2 error

		if appConfig.Branch == "latest" {
			latestTag, e3 := getLatestTag(appConfig.TmplPath)
			// This error is informative, but not worth stopping the program.
			logf(e3.Error())
			if latestTag != "" {
				appConfig.Branch = latestTag
			}
		}

		// Determine the cache location
		repoDir := appConfig.UsrOpts.CacheDir + cli.PS + getRepoDir(appConfig.TmplPath, appConfig.Branch)
		infof(cli.Messages.OutRepoDir, repoDir)

		// Do a pull when the repo already exists. This will fail if it downloaded a zip.
		if stdlib.DirExist(repoDir + cli.PS + gitConfDir) {
			infof(cli.Messages.UsingCache, repoDir)
			repo, commitHash, err2 = gitCheckout(repoDir, appConfig.Branch)
		} else {
			infof(cli.Messages.CloningToCache, repoDir)
			repo, commitHash, err2 = gitClone(appConfig.TmplPath, repoDir, appConfig.Branch)
		}

		infof(cli.Messages.RepoInfo, repo, commitHash)
		if err2 != nil {
			mainErr = err2
			return
		}
		appConfig.Tmpl = repo
	}

	if !stdlib.DirExist(appConfig.Tmpl) {
		mainErr = fmt.Errorf(cli.Errors.InvalidTmplDir, appConfig.Tmpl)
		return
	}

	fec, err1 := stdlib.NewFileExtChecker(appConfig.UsrOpts.ExcludeFileExtensions, &[]string{})
	if err1 != nil {
		mainErr = fmt.Errorf(cli.Errors.CannotInitFileChecker, err1.Error())
	}

	// Require template directories to have a specific file in order to be processed to prevent processing directories unintentionally.
	tmplManifestFile := appConfig.Tmpl + cli.PS + cli.TmplManifest
	tmplManifest, errX := cli.ReadTemplateJson(tmplManifestFile)
	if errX != nil {
		mainErr = fmt.Errorf(cli.Errors.MissingTmplJson, cli.TmplManifest, tmplManifestFile, errX.Error())
		return
	}

	appConfig.TmplJson = tmplManifest
	appConfig.AnswersJson = cli.NewAnswerJson()

	if stdlib.PathExist(appConfig.AnswersPath) {
		appConfig.AnswersJson, mainErr = cli.LoadAnswers(appConfig.AnswersPath)
		if mainErr != nil {
			return
		}
	}

	// Checks for any missing placeholder values waits for their input from the CLI.
	if e := cli.GetPlaceholderInput(appConfig.TmplJson, &appConfig.AnswersJson.Placeholders, os.Stdin, appConfig.DefaultVal); e != nil {
		mainErr = fmt.Errorf(cli.Errors.GettingAnswers, e.Error())
	}

	cli.ShowAllPlaceholderValues(appConfig.TmplJson, &appConfig.AnswersJson.Placeholders)

	mainErr = cli.ParseDir(appConfig.Tmpl, appConfig.OutPath, appConfig.AnswersJson.Placeholders, fec, appConfig.TmplJson)
}
