package config

import (
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmplpress/internal/msg"
	"github.com/kohirens/tmplpress/internal/press"
)

type Arguments struct {
	Setting string // config setting
	Method  string // Method to call
	Value   string // value to update config setting
}

const (
	Summary = "Update or retrieve a configuration value"
	Name    = "config"
)

var (
	args  = &Arguments{}
	flags *flag.FlagSet
	help  bool
)

func Init() *flag.FlagSet {
	flags = flag.NewFlagSet(Name, flag.ExitOnError)

	flags.BoolVar(&help, "help", false, UsageMessages["help"])

	return flags
}

func ParseInput(ca []string) error {
	if e := flags.Parse(ca); e != nil {
		return fmt.Errorf(msg.Stderr.ParsingConfigArgs, e.Error())
	}

	if help {
		return nil
	}

	if len(ca) < 2 {
		return fmt.Errorf(msg.Stderr.InvalidNoArgs)
	}

	args.Method = ca[0]
	args.Setting = ca[1]

	log.Dbugf(msg.Stdout.ConfigMethodSetting, args.Method, args.Setting)

	if len(ca) > 2 {
		log.Dbugf(msg.Stdout.SetValue, args.Method, args.Value)
		args.Value = ca[2]
	}

	if args.Method == "set" && len(ca) != 3 {
		return fmt.Errorf(Stderr.ConfigValueNotSet)
	}

	return nil
}

// Run Set or get a config setting.
func Run(ca []string, appName string) error {
	if e := ParseInput(ca); e != nil {
		return e
	}

	if help {
		flags.Usage()
		return nil
	}

	appDataDir, err1 := press.BuildAppDataPath(appName)
	if err1 != nil {
		return err1
	}

	cp := appDataDir + fsio.PS + press.ConfigFileName

	switch args.Method {
	case "set":
		return set(args.Setting, args.Value, cp, appName)

	case "get":
		return get(args.Setting, cp)

	default:
		return fmt.Errorf(Stderr.InvalidConfigMethod, args.Method)
	}
}

// get the value of a user setting.
func get(key string, cp string) error {
	var val interface{}

	sc, err1 := press.LoadConfig(cp)
	if err1 != nil {
		return err1
	}

	switch key {
	case "CacheDir":
		val = sc.CacheDir
		log.Logf("%v", val)
		break

	default:
		return fmt.Errorf(msg.Stderr.NoSetting, key)
	}

	return nil
}

// set the value of a user setting
func set(key, val string, cp, appName string) error {
	sc := &press.ConfigSaveData{}

	switch key {
	case "CacheDir":
		sc.CacheDir = val
		break

	default:
		return fmt.Errorf(msg.Stderr.NoSetting, key)
	}

	if !fsio.Exist(cp) {
		_, e := press.InitConfig(cp, appName)
		return e
	}

	return press.SaveConfig(cp, sc)
}
