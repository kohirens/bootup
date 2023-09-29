// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"github.com/kohirens/tmpltoapp/subcommand/config"
	"github.com/kohirens/tmpltoapp/subcommand/manifest"
)

type appFlags struct {
	AnswersPath    string // The path to a file containing values to variables to be parsed.
	Branch         string // The desired branch of the template to.
	CommitHash     string // Git commit hash of the current version.
	CurrentVersion string // Current semantic version of the application.
	DefaultVal     string // A default placeholder value when a placeholder is empty.
	Help           bool   // The usage for all flags.
	TmplPath       string // The URL or local template path to a template.
	TmplType       string // Indicate the type of package for a template, such as a zip to Extract or a repository to Download.
	OutPath        string // The location to save the processed template output.
	Version        bool   // The current version
	subcommands    map[string]*flag.FlagSet
}

// define All application flags.
func defineFlags(af *appFlags) {
	// Note: These are defined in alphabetical order.
	flag.StringVar(&af.AnswersPath, "answer-path", "", um["answer-path"]) // TODO: BREAKING Change to "answers"
	flag.StringVar(&af.Branch, "branch", "main", um["branch"])            // TODO: BREAKING Change git-ref, since refs alreay point to a complete SHA-1
	flag.StringVar(&af.DefaultVal, "default-val", " ", um["default-val"])
	flag.BoolVar(&af.Help, "help", false, um["help"])
	flag.BoolVar(&af.Help, "h", false, um["help"]+" (shorthand)")
	flag.StringVar(&af.OutPath, "out-path", "", um["out-path"])       // TODO: BREAKING remove this will be a required 2nd argument.
	flag.StringVar(&af.TmplPath, "tmpl-path", "", um["tmpl-path"])    // TODO: BREAKING remove this will be a required 1st argument.
	flag.StringVar(&af.TmplType, "tmpl-type", "git", um["tmpl-type"]) // TODO: BREAKING Remove, we only use git now.
	flag.IntVar(&verbosityLevel, "verbosity", log.VerboseLvlLog, um["verbosity"])
	flag.BoolVar(&af.Version, "version", false, um["version"])
}

// Parse and validate all global flags.
func parseCli(af *appFlags, cfg *cli.AppData) error {
	// Remember that flag parsing stops just before the first argument that does not have a "-" and is also NOT the
	// value of a flag or comes after the terminator "--".
	// It was planed to allow for flags/arguments in any order, but it may be less confusing to only support flag first
	// and then arguments; it may also require less code to debug and document for not very much gain.
	flag.Parse()

	log.Infof(msg.Stdout.VerboseLevelInfo, verbosityLevel)

	ca := flag.Args()

	log.Dbugf(msg.Stdout.NumNonFlagArgs, len(ca))
	log.Dbugf(msg.Stdout.ActualArgs, ca)
	log.Dbugf(msg.Stdout.NumParsedFlags, flag.NFlag())

	if verbosityLevel == verboseLvlDbug {
		fmt.Println(msg.Stdout.PrintAllFlags)
		flag.Visit(func(f *flag.Flag) {
			fmt.Printf(msg.Stdout.PrintFlag, f.Name, f.Value, f.DefValue)
		})
	}

	if af.Help {
		return nil
	}

	if len(ca) == 0 && flag.NFlag() > 0 {
		return nil
	}

	if len(ca) == 0 {
		return fmt.Errorf(msg.Stderr.NoInput)
	}

	switch flag.Arg(0) {
	case config.Name:
		cfg.SubCmd = config.Name
		if e := config.ParseFlags(ca[1:]); e != nil {
			return e
		}
	case manifest.Name:
		cfg.SubCmd = manifest.Name
		if e := manifest.ParseFlags(ca[1:]); e != nil {
			return e
		}
	}

	return nil
}
