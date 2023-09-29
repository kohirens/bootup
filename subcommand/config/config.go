package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"os"
	"regexp"
	"strings"
)

type Arguments struct {
	Setting string // config setting
	Method  string // Method to call
	Value   string // value to update config setting
}

const (
	dirMode = 0774
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

	if len(ca) > 2 {
		args.Value = ca[2]
	}

	if args.Method == "set" && len(ca) != 3 {
		return fmt.Errorf(Stderr.ConfigValueNotSet)
	}

	log.Dbugf("args.Method = %v\n", args.Method)
	log.Dbugf("args.Setting = %v\n", args.Setting)
	log.Dbugf("args.Value = %v\n", args.Value)

	return nil
}

// Run Set or get a config setting.
func Run(ca []string, cfg *cli.AppData) error {
	if e := ParseInput(ca); e != nil {
		return e
	}

	if help {
		flags.Usage()
		return nil
	}

	log.Dbugf("args.Method = %v\n", args.Method)
	log.Dbugf("args.Key = %v\n", args.Setting)

	switch args.Method {
	case "set":
		log.Dbugf("args.Value = %v\n", args.Value)
		fmt.Printf("\n\n%v\n\n", args.Value)
		if e := set(args.Setting, args.Value, cfg); e != nil {
			return e
		}

		return save(cfg, dirMode)

	case "get":
		v, e := get(args.Setting, cfg)

		if e != nil {
			return e
		}

		log.Logf("%v", v)

	default:
		return fmt.Errorf(Stderr.InvalidConfigMethod, args.Method)
	}

	return nil
}

// get the value of a user setting.
func get(key string, cfg *cli.AppData) (interface{}, error) {
	var val interface{}

	switch key {
	case "CacheDir":
		val = cfg.UsrOpts.CacheDir
		break
	case "ExcludeFileExtensions":
		v2 := fmt.Sprintf("%v", val)
		ok, _ := regexp.Match("^[a-zA-Z0-9-.]+(?:,[a-zA-Z0-9-.]+)*", []byte(v2))
		if !ok {
			return nil, fmt.Errorf(Stderr.BadExcludeFileExt, val)
		}
		val = strings.Join(*cfg.UsrOpts.ExcludeFileExtensions, ",")
		break
	default:
		return "", fmt.Errorf("no setting %v found", key)
	}

	return val, nil
}

// save configuration file.
func save(cfg *cli.AppData, mode os.FileMode) error {
	data, err1 := json.Marshal(cfg.UsrOpts)

	log.Dbugf(Stdout.SaveData, data)

	if err1 != nil {
		return fmt.Errorf(Stderr.CouldNotEncodeConfig, err1.Error())
	}

	if e := os.WriteFile(cfg.Path, data, mode); e != nil {
		return e
	}

	return nil
}

// set the value of a user setting
func set(key, val string, cfg *cli.AppData) error {
	switch key {
	case "CacheDir":
		log.Dbugf("setting CacheDir = %q", val)
		cfg.UsrOpts.CacheDir = val
		break
	case "ExcludeFileExtensions":
		log.Dbugf("adding exclusions %q to config", val)
		tmp := strings.Split(val, ",")
		cfg.UsrOpts.ExcludeFileExtensions = &tmp
		break
	default:
		return fmt.Errorf("no %q setting found", key)
	}

	return nil
}
