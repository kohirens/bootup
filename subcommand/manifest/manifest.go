package manifest

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"github.com/kohirens/tmpltoapp/internal/press"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	txtParse "text/template/parse"
)

const (
	ps           = string(os.PathSeparator)
	GitConfigDir = ".git"
	Name         = "manifest"
	Summary      = "Generate a template.json file for a template."
)

type Arguments struct {
	Cmd  string // command to run.
	Path string // path to generate a manifest for.
}

var (
	input Arguments
	flags *flag.FlagSet
	help  bool
)

func Init() *flag.FlagSet {
	flags = flag.NewFlagSet(Name, flag.ExitOnError)

	flags.BoolVar(&help, "help", false, UsageMessages["help"])

	return flags
}

func parseInput(ca []string) error {
	if e := flags.Parse(ca); e != nil {
		return fmt.Errorf(msg.Stderr.ParsingConfigArgs, e.Error())
	}

	if help {
		flags.Usage()
		return nil
	}

	if len(ca) < 1 {
		flags.Usage()
		return fmt.Errorf(msg.Stderr.InvalidNoSubCmdArgs, Name, 1)
	}

	if len(ca) < 2 {
		// Fall back to the current working directory when no path is specified.
		cwd, e1 := os.Getwd()
		if e1 != nil {
			return fmt.Errorf(stderr.ListWorkingDirectory, e1.Error())
		}

		log.Dbugf(msg.Stdout.Cwd, cwd)

		ca = append(ca, cwd)
	}

	input.Cmd = ca[0]
	// clean up the path.
	p, e1 := filepath.Abs(ca[1])
	if e1 != nil {
		return fmt.Errorf(msg.Stderr.NoPath, e1.Error())
	}

	input.Path = p

	log.Dbugf(msg.Stdout.TemplatePath, input.Path)

	return nil
}

func Run(ca []string) error {
	if e := parseInput(ca); e != nil {
		return e
	}

	switch input.Cmd {
	default:
		flags.Usage()
		return fmt.Errorf(msg.Stderr.InvalidCmd, input.Cmd)
	case "generate":
		filename, e1 := generateATemplateManifest(input.Path)
		if e1 != nil {
			return e1
		}

		log.Logf(msg.Stdout.GeneratedManifest, filename)
	}

	return nil
}

// generateATemplateManifest Make a JSON file with your templates placeholders.
func generateATemplateManifest(tmplPath string) (string, error) {
	if !fsio.Exist(tmplPath) {
		return "", fmt.Errorf(msg.Stderr.PathNotExist, tmplPath)
	}

	filename := tmplPath + ps + press.TmplManifestFile
	// otherwise, start with the default template.json
	tm, e1 := press.NewTmplManifest([]byte(defaultJson))
	if e1 != nil {
		return "", e1
	}

	// check for existing template manifest and load it
	existing, e2 := press.ReadTemplateJson(filename)
	if e2 != nil {
		log.Infof(e2.Error())
	}

	if existing != nil { // merge old into the new updating it at the same time.
		tm.CopyAsIs = existing.CopyAsIs
		tm.EmptyDirFile = existing.EmptyDirFile
		tm.Placeholders = existing.Placeholders
		tm.Skip = existing.Skip
		tm.Substitute = existing.Substitute
		tm.Validation = existing.Validation
	}

	// Traverse the path recursively, filtering out files that should be excluded
	templates, e3 := parseDir(tmplPath, tm)
	if e3 != nil {
		return "", e3
	}

	actions := make(map[string]string)

	// Parse the file as a template and extract all actions from each file.
	for _, tmpl := range templates {
		fmt.Printf("checking %v\n", tmpl)

		t, e := template.ParseFiles(tmpl)
		if e != nil {
			return "", fmt.Errorf(msg.Stderr.ParsingFile, tmpl, e.Error())
		}

		listTemplateFields(t, actions)
	}

	tm.Placeholders = actions

	if e := saveFile(filename, tm); e != nil {
		return "", e
	}

	return filename, nil
}

// listTemplateFields list actions in Go templates. See SO answer: https://stackoverflow.com/a/40584967/419097
func listTemplateFields(t *template.Template, res map[string]string) {
	listNodeFields(t.Tree.Root, res)
}

// parseDir Recursively walk a directory parsing all files along the way as Go templates.
func parseDir(path string, tm *press.TmplManifest) ([]string, error) {
	// Normalize the path separator in these 2 variables before comparing them.
	nPath := fsio.Normalize(path)

	var files []string

	// Recursively walk the template directory.
	err := filepath.Walk(nPath, func(fPath string, info fs.FileInfo, err error) error {

		file, e1 := filterFile(fPath, nPath, info, err, tm)
		if err != nil {
			return e1
		}

		if file != "" {
			files = append(files, file)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// filterFile
func filterFile(sourcePath, nPath string, info os.FileInfo, wErr error, tm *press.TmplManifest) (string, error) {
	if wErr != nil {
		return "", wErr
	}

	if info.IsDir() {
		return "", nil
	}

	// skip certain .git files/directories
	if strings.Contains(sourcePath, ps+GitConfigDir+ps) {
		return "", nil
	}

	currFile := filepath.Base(sourcePath)

	// TODO: Add globbing is added. filepath.Glob(pattern)
	if currFile == tm.EmptyDirFile || currFile == press.TmplManifestFile { // Use an exclusion list, include every file by default.
		return "", nil
	}

	// Normalize the path separator in these 2 variables before comparing them.
	normSourcePath := strings.ReplaceAll(sourcePath, "/", ps)
	normSourcePath = strings.ReplaceAll(normSourcePath, "\\", ps)

	// Skip files that are listed in the excludes.
	if tm.CopyAsIs != nil {
		fileToCheck := strings.ReplaceAll(normSourcePath, nPath, "")
		fileToCheck = strings.ReplaceAll(fileToCheck, ps, "")

		for _, exclude := range tm.CopyAsIs {
			fileToCheckB := strings.ReplaceAll(exclude, "\\", "")
			fileToCheckB = strings.ReplaceAll(exclude, "/", "")

			if fileToCheckB == fileToCheck {
				return "", nil
			}
		}
	}

	return sourcePath, nil
}

// listTemplateFields list actions in Go templates. See SO answer: https://stackoverflow.com/a/40584967/419097
func listNodeFields(node txtParse.Node, res map[string]string) {
	if node.Type() == txtParse.NodeAction {
		res[strings.Trim(node.String(), "{}.")] = ""
	}

	if ln, ok := node.(*txtParse.ListNode); ok {
		for _, n := range ln.Nodes {
			listNodeFields(n, res)
		}
	}
}

type templateSchema struct {
	Placeholders []byte
}

// save configuration file.
func saveFile(jsonFile string, tm *press.TmplManifest) error {
	data, e1 := json.MarshalIndent(tm, "", "    ")

	if e1 != nil {
		return fmt.Errorf(stderr.EncodingJson, jsonFile, e1.Error())
	}

	// Write the template.json manifest to disk.
	if e := os.WriteFile(jsonFile, data, 0744); e != nil {
		return e
	}

	return nil
}
