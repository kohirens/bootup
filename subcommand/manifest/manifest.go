package manifest

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/stdlib/path"
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
	EmptyFile    = ".empty"
	GitConfigDir = ".git"
	Name         = "manifest"
	Summary      = "Generate a template.json file for a template."
)

type Arguments struct {
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
		// Fall back to the current working directory when no path is specified.
		cwd, e1 := os.Getwd()
		if e1 != nil {
			return fmt.Errorf(stderr.ListWorkingDirectory, e1.Error())
		}

		log.Dbugf("current working directory = %v", cwd)

		ca = append(ca, cwd)
	}

	// clean up the path.
	p, e1 := filepath.Abs(ca[0])
	if e1 != nil {
		return fmt.Errorf("invalid path, %v", e1.Error())
	}

	input.Path = p

	log.Dbugf("manifest.input.Path = %v", input.Path)

	return nil
}

func Run(ca []string) error {
	if e := parseInput(ca); e != nil {
		return e
	}

	// TODO: BREAKING Add this to the template.json, the template designer should be responsible for this; ".empty" should still be embedded in this app though.
	fec, _ := stdlib.NewFileExtChecker(&[]string{".empty", "exe", "gif", "jpg", "mp3", "pdf", "png", "tiff", "wmv"}, &[]string{})

	_, e1 := generateATemplateManifest(input.Path, fec, []string{})
	if e1 != nil {
		return e1
	}

	return nil
}

// generateATemplateManifest Make a JSON file with your templates placeholders.
func generateATemplateManifest(tmplPath string, fec *stdlib.FileExtChecker, excludes []string) (map[string]string, error) {
	if !path.Exist(tmplPath) {
		return nil, fmt.Errorf(msg.Stderr.PathNotExist, tmplPath)
	}

	// Traverse the path recursively, filtering out files that should be excluded
	templates, err := parseDir(tmplPath, fec, excludes)
	if err != nil {
		return nil, err
	}

	actions := make(map[string]string)

	// Parse the file as a template and extract all actions from each file.
	for _, tmpl := range templates {
		fmt.Printf("checking %v\n", tmpl)

		t, e := template.ParseFiles(tmpl)
		if e != nil {
			return nil, fmt.Errorf(msg.Stderr.ParsingFile, tmpl, e.Error())
		}

		listTemplateFields(t, actions)
	}

	if e := saveFile(tmplPath+ps+press.TmplManifestFile, actions); e != nil {
		return nil, e
	}

	return actions, nil
}

// listTemplateFields list actions in Go templates. See SO answer: https://stackoverflow.com/a/40584967/419097
func listTemplateFields(t *template.Template, res map[string]string) {
	listNodeFields(t.Tree.Root, res)
}

// parseDir Recursively walk a directory parsing all files along the way as Go templates.
func parseDir(path string, fec *stdlib.FileExtChecker, excludes []string) ([]string, error) {
	// Normalize the path separator in these 2 variables before comparing them.
	nPath := strings.ReplaceAll(path, "/", ps)
	nPath = strings.ReplaceAll(nPath, "\\", ps)

	var files []string
	i := 0
	// Recursively walk the template directory.
	err := filepath.Walk(nPath, func(fPath string, info fs.FileInfo, err error) error {
		i++
		//fmt.Printf("%-2d %v\n", i, fPath)

		file, e1 := filterFile(fPath, nPath, info, err, excludes)
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
func filterFile(sourcePath, nPath string, info os.FileInfo, wErr error, excludes []string) (string, error) {
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

	// Skip files by extension.
	// TODO: Add globbing is added. filepath.Glob(pattern)
	if currFile == EmptyFile || currFile == press.TmplManifestFile { // Use an exclusion list, include every file by default.
		return "", nil
	}

	// Normalize the path separator in these 2 variables before comparing them.
	normSourcePath := strings.ReplaceAll(sourcePath, "/", ps)
	normSourcePath = strings.ReplaceAll(normSourcePath, "\\", ps)

	// Skip files that are listed in the excludes.
	if excludes != nil {
		fileToCheck := strings.ReplaceAll(normSourcePath, nPath, "")
		fileToCheck = strings.ReplaceAll(fileToCheck, ps, "")

		for _, exclude := range excludes {
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
func saveFile(jsonFile string, actions map[string]string) error {
	data, e1 := json.Marshal(actions)

	if e1 != nil {
		return fmt.Errorf(stderr.EncodingJson, jsonFile, e1.Error())
	}

	tmpl := template.Must(template.New(press.TmplManifestFile).Parse(TmplJsonTmpl))

	f, e2 := os.Create(jsonFile)
	if e2 != nil {
		return e2
	}

	// Write the template.json manifest to disk.
	if e := tmpl.Execute(f, templateSchema{Placeholders: data}); e != nil {
		return fmt.Errorf(stderr.SavingManifest, jsonFile, e.Error())
	}

	return nil
}
