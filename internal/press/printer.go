// Package press serves as the print head to press the template into a working
// output of the intended purpose of the designer. So if the designer design a
// template for an application the output will produce just that.
package press

import (
	"bufio"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const (
	dirMode      = 0744
	emptyFile    = ".empty"
	gitConfigDir = ".git"
	maxTmplSize  = 1e+7
	PS           = string(os.PathSeparator)
)

// FindTemplates Recursively walk a directory looking for files along the way.
func FindTemplates(dir string) ([]string, error) {
	// Normalize the path separator in these 2 variables before comparing them.
	normTplDir := path.Normalize(dir)

	files := []string{}

	// Recursively walk the template directory.
	e1 := filepath.Walk(normTplDir, func(tmpl string, fi os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}

		// Skip directories.
		if fi.IsDir() {
			log.Dbugf("skipping directory %v", tmpl)
			return nil
		}

		log.Infof("adding file %v", tmpl)

		files = append(files, tmpl)

		return nil
	})

	if e1 != nil {
		return nil, e1
	}

	return files, nil
}

// Print templates to the output directory.
func Print(tplDir, outDir string, vars cli.StringMap, fec *stdlib.FileExtChecker, tmplJson *templateJson) error {
	if !path.Exist(tplDir) {
		return fmt.Errorf(msg.Stderr.PathNotExist, tplDir)
	}

	// Normalize the path separator in these 2 variables.
	normTplDir := path.Normalize(tplDir)
	normOutDir := path.Normalize(outDir)
	log.Infof("%v -> %v", tplDir, normTplDir)
	log.Infof("%v -> %v", outDir, normOutDir)

	// Recursively walk the template directory.
	return filepath.Walk(normTplDir, func(sourcePath string, fi os.FileInfo, wErr error) error {
		if wErr != nil {
			return wErr
		}

		log.Infof("processing: %v", sourcePath)

		// Do not parse directories.
		if fi.IsDir() {
			return nil
		}

		// Skip processing files if a template file is too big.
		if fi.Size() > maxTmplSize {
			return fmt.Errorf(msg.Stderr.FileTooBig, maxTmplSize)

		}

		currFile := filepath.Base(sourcePath)

		// Skip when the extension matches a file extension ignore list.
		// these files are included in the output but skipped by the template
		// processor.
		if !fec.IsValid(sourcePath) { // Skip files by extension; Use an exclusion list, include every file by default.
			log.Infof(msg.Stdout.UnknownFileType, sourcePath)
			return nil
		}

		// Normalize the path separator in these 2 variables before comparing them.
		normSourcePath := path.Normalize(sourcePath)

		// Get the relative path of the file from root of the template and
		// append it to the output directory, so that files are placed in the
		// same subdirectories in the output directory.
		relativePath := strings.TrimLeft(strings.ReplaceAll(normSourcePath, normTplDir, ""), "\\/")
		log.Infof("relativePath dir: %v", relativePath)

		saveDir := filepath.Clean(normOutDir + PS + filepath.Dir(relativePath))
		log.Infof("save dir: %v", saveDir)

		// Skip template manifest file and the git config directory.
		if currFile == TmplManifestFile || strings.Contains(relativePath, gitConfigDir+PS) {
			log.Infof(msg.Stdout.SkipFile, relativePath)
			return nil
		}

		if inSkipArray(relativePath, tmplJson.Skip) { // Skip files in this list
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return nil
		}

		// skip the directory with replace files.
		if tmplJson.Substitute != "" {
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return nil
		}

		// Make the subdirectories in the new savePath.
		if e := os.MkdirAll(saveDir, dirMode); e != nil {
			return e
		}

		// we skip the designated empty file here so that the directory is made.
		if currFile == emptyFile {
			return nil
		}

		// TODO: Replace with better method of comparing files.
		if tmplJson.Excludes != nil { // exclude from parsing, but copy as-is.
			fileToCheck := strings.ReplaceAll(normSourcePath, normTplDir, "")
			fileToCheck = strings.ReplaceAll(fileToCheck, PS, "")

			for _, exclude := range tmplJson.Excludes {
				// we want to see if the paths are the same, so we remove
				fileToCheckB := strings.ReplaceAll(exclude, "\\", "")
				fileToCheckB = strings.ReplaceAll(exclude, "/", "")
				if fileToCheckB == fileToCheck {
					log.Infof("file %v will be copied as-is", sourcePath)
					_, errC := copyToDir(sourcePath, saveDir, PS)
					return errC
				}
			}
		}

		return parse(sourcePath, saveDir, vars)
	})
}

// GetPlaceholderInput Checks for any missing placeholder values waits for their input from the CLI.
func GetPlaceholderInput(placeholders *templateJson, tmplValues cli.StringMap, r *os.File, defaultVal string) error {
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
			log.Infof("using default value for placeholder %v", placeholder)
			continue
		}

		// Ask client for input.
		fmt.Printf("\n%v - %v: ", placeholder, desc)
		nPut.Scan()
		tVals[placeholder] = nPut.Text()
		log.Infof(msg.Stdout.PlaceholderAnswer, desc, tVals[placeholder])
		log.Infof("%v = %q\n", placeholder, tVals[placeholder])
	}

	return nil
}

func ShowAllPlaceholderValues(placeholders *templateJson, tmplValues *cli.StringMap) {
	tVals := *tmplValues
	log.Logf("the following values have been provided\n")
	for placeholder := range placeholders.Placeholders {
		log.Logf(msg.Stdout.PlaceholderAnswer, placeholder, tVals[placeholder])
	}
}

func copyAsIs(tmplRoot, file, out string, excludedFiles []string) (rVal bool) {
	rVal = false

	if excludedFiles == nil {
		return
	}

	// Using the relative path should always remove the drive letter on windows
	// and make it uniform across all OS.
	rp := strings.Replace(file, tmplRoot, "", 1)

	if runtime.GOOS == "windows" { // I dislike the fact that Windows is case-insensitive.
		rp = strings.ToLower(rp)
	}

	log.Logf("fileToCheck: %q against excludes", rp)

	for _, ef := range excludedFiles {
		ef = path.Normalize(ef)
		if runtime.GOOS == "windows" { // I dislike the fact that Windows is case-insensitive.
			ef = strings.ToLower(ef)
		}

		if ef == rp {
			rVal = true
			break
		}
	}

	return
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

// parse a file as a Go template.
func parse(tplFile, dstDir string, vars cli.StringMap) error {
	log.Infof("parsing %v", tplFile)
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
