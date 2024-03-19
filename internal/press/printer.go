// Package press serves as the print head to press the template into a working
// output of the intended purpose of the designer. So if the designer design a
// template for an application the output will produce just that.
package press

import (
	"bufio"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"github.com/ryanuber/go-glob"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	dirMode      = 0744
	gitConfigDir = ".git"
	maxTmplSize  = 1e+7
	PS           = string(os.PathSeparator)
)

// FindTemplates Recursively walk a directory looking for files along the way.
func FindTemplates(dir string) ([]string, error) {
	// Normalize the path separator in these 2 variables before comparing them.
	normTplDir := fsio.Normalize(dir)

	files := []string{}

	// Recursively walk the template directory.
	e1 := filepath.Walk(normTplDir, func(tmpl string, fi os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}

		// Skip directories.
		if fi.IsDir() {
			log.Dbugf(msg.Stdout.Skipping, tmpl)
			return nil
		}

		log.Infof(msg.Stdout.AddFile, tmpl)

		files = append(files, tmpl)

		return nil
	})

	if e1 != nil {
		return nil, e1
	}

	return files, nil
}

// Print templates to the output directory.
func Print(tplDir, outDir string, vars map[string]string, tmplJson *TmplManifest) error {
	if !fsio.Exist(tplDir) {
		return fmt.Errorf(msg.Stderr.PathNotExist, tplDir)
	}

	// Normalize the path separator in these 2 variables.
	normTplDir := fsio.Normalize(tplDir)
	normOutDir := fsio.Normalize(outDir)
	log.Infof("template: %v", normTplDir)
	log.Infof("output: %v", normOutDir)

	// Recursively walk the template directory.
	return filepath.Walk(normTplDir, func(sourcePath string, fi os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}

		// Do not parse directories.
		if fi.IsDir() {
			return nil
		}

		if strings.Contains(sourcePath, gitConfigDir+PS) {
			return nil
		}

		// Skip all files in the substitute directory.
		if hasParentDir(tmplJson.Substitute, sourcePath) {
			log.Infof("substitute skip %v", sourcePath)
			return nil
		}

		// Skip processing files if a template file is too big.
		if fi.Size() > maxTmplSize {
			return fmt.Errorf(msg.Stderr.FileTooBig, maxTmplSize)
		}

		log.Infof(msg.Stdout.Processing, sourcePath)

		currFile := filepath.Base(sourcePath)

		// Normalize the path separator before performing any operations.
		normSourcePath := fsio.Normalize(sourcePath)

		// Get the relative path of the file from root of the template and
		// append it to the output directory, so that files are placed in their
		// correct subdirectories in the output.
		relativePath := strings.TrimLeft(strings.ReplaceAll(normSourcePath, normTplDir, ""), "\\/")
		log.Infof(msg.Stdout.RelativeDir, relativePath)

		// Skip the template manifest file and the git config directory.
		if currFile == TmplManifestFile {
			log.Infof(msg.Stdout.Skipping, relativePath)
			return nil
		}

		// Don't do anything with the files in this list.
		if inSkipArray(relativePath, tmplJson.Skip) {
			log.Infof(msg.Stdout.Skipping, sourcePath)
			return nil
		}

		saveDir := filepath.Clean(normOutDir + PS + filepath.Dir(relativePath))
		log.Infof(msg.Stdout.SaveDir, saveDir)

		// Make all subdirectories in output path.
		if e := os.MkdirAll(saveDir, dirMode); e != nil {
			return e
		}

		// For empty directories, return here so that the directory is made and nothing else.
		if currFile == tmplJson.EmptyDirFile {
			return nil
		}

		// TODO: replace normSourcePath, normTplDir parameters to use relativePath.
		copied, e1 := copyAsIs(tmplJson.Excludes, relativePath, sourcePath, saveDir)
		if e1 != nil {
			return e1
		} else if copied {
			return nil
		}

		return parse(sourcePath, saveDir, vars)
	})
}

// copyAsIs Check a file matches a glob pattern, if so, then copy it to the
// output as-is (without template parsing).
func copyAsIs(files []string, relativePath, sourcePath, saveDir string) (bool, error) {
	if len(files) < 1 { // no-op
		return false, nil
	}

	for _, exclude := range files {
		// check if the file matches a pattern
		if glob.Glob(exclude, relativePath) {
			log.Infof(msg.Stdout.CopyAsIs, sourcePath)
			_, e := copyToDir(sourcePath, saveDir, PS)
			return true, e
		}
	}

	return false, nil
}

// GetPlaceholderInput Checks for any missing placeholder values waits for their input from the CLI.
func GetPlaceholderInput(placeholders *TmplManifest, tmplValues map[string]string, r *os.File, defaultVal string) error {
	tVals := tmplValues
	nPut := bufio.NewScanner(r)

	for placeholder, desc := range placeholders.Placeholders {
		a, answered := tVals[placeholder]
		// skip placeholder that have been supplied with an answer from an answer file.

		if answered {
			log.Infof(msg.Stdout.PlaceholderHasAnswer, desc, a)
			continue
		}

		// Just use the default value for all un-set placeholders.
		if defaultVal != " " {
			tVals[placeholder] = defaultVal
			log.Infof(msg.Stdout.VarDefaultValue, placeholder)
			continue
		}

		// Ask client for input.
		fmt.Printf("\n%v - %v: ", placeholder, desc)
		nPut.Scan()
		tVals[placeholder] = nPut.Text()
		log.Infof(msg.Stdout.Assignment, desc, tVals[placeholder])
		log.Infof(msg.Stdout.Assignment, placeholder, tVals[placeholder])
	}

	return nil
}

func ShowAllPlaceholderValues(tm *TmplManifest, tmplValues map[string]string) {
	if tm.Placeholders == nil {
		log.Logf(msg.Stdout.NoPlaceholders)
		return
	}

	log.Logf(msg.Stdout.ValuesProvided)
	for placeholder := range tm.Placeholders {
		log.Logf(msg.Stdout.Assignment, placeholder, tmplValues[placeholder])
	}
}

// copyToDir Copy a file to a directory.
func copyToDir(sourcePath, destDir, separator string) (int64, error) {
	//TODO: Move to stdlib.
	sFile, err1 := os.Open(sourcePath)
	if err1 != nil {
		return 0, err1
	}

	fileStats, err2 := os.Stat(sourcePath)
	if err2 != nil {
		return 0, err2
	}

	dstFile := destDir + separator + fileStats.Name()
	dFile, err3 := os.Create(dstFile)
	if err3 != nil {
		return 0, err3
	}

	return io.Copy(dFile, sFile)
}

// hasParentDir Detect if a strings is part of a path.
func hasParentDir(parent, dir string) bool {
	if parent != "" && strings.Contains(dir, parent) {
		return true
	}
	return false
}

// parse a file as a Go template.
func parse(tplFile, dstDir string, vars map[string]string) error {
	log.Infof(msg.Stdout.Parsing, tplFile)
	funcMap := template.FuncMap{
		"title":   strings.Title,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}

	tmplName := filepath.Base(tplFile)
	parser, err1 := template.New(tmplName).Funcs(funcMap).ParseFiles(tplFile)
	if err1 != nil {
		return err1
	}

	fileStats, err2 := os.Stat(tplFile)
	if err2 != nil {
		return err2
	}

	dstFile := dstDir + PS + fileStats.Name()
	file, err3 := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY, fileStats.Mode())
	if err3 != nil {
		return err3
	}

	if e := parser.Execute(file, vars); e != nil {
		return e
	}

	if e := file.Close(); e != nil {
		return e
	}

	return nil
}
