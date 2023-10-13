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
	e1 := filepath.Walk(normTplDir, func(sourcePath string, fi os.FileInfo, wErr error) (rErr error) {
		if wErr != nil {
			rErr = wErr
			return
		}

		log.Infof("adding %v", sourcePath)

		// Skip directories.
		if fi.IsDir() {
			return
		}

		return
	})

	if e1 != nil {
		return nil, e1
	}

	return files, nil
}

// Print templates to the output directory.
func Print(tplDir, outDir string, vars cli.StringMap, fec *stdlib.FileExtChecker, tmplJson *templateJson) (err error) {
	// Normalize the path separator in these 2 variables before comparing them.
	normTplDir := strings.ReplaceAll(tplDir, "/", PS)
	normTplDir = strings.ReplaceAll(normTplDir, "\\", PS)

	// Recursively walk the template directory.
	err = filepath.Walk(normTplDir, func(sourcePath string, fi os.FileInfo, wErr error) (rErr error) {
		if wErr != nil {
			rErr = wErr
			return
		}

		log.Infof("\nprocessing: %q", sourcePath)

		// Do not parse directories.
		if fi.IsDir() {
			return
		}

		// Stop processing files if a template file is too big.
		if fi.Size() > maxTmplSize {
			rErr = fmt.Errorf(msg.Stderr.FileTooBig, maxTmplSize)
			return
		}

		currFile := filepath.Base(sourcePath)
		if currFile != emptyFile && !fec.IsValid(sourcePath) { // Skip files by extension; Use an exclusion list, include every file by default.
			log.Infof(msg.Stdout.UnknownFileType, sourcePath)
			return
		}

		// Normalize the path separator in these 2 variables before comparing them.
		normSourcePath := path.Normalize(sourcePath)
		// Get the relative path of the file from root of the template and
		// append it to the output directory, so that files are placed in the
		// same subdirectories in the output directory.
		relativePath := strings.TrimLeft(strings.ReplaceAll(normSourcePath, normTplDir, ""), "\\/")
		saveDir := filepath.Clean(outDir + PS + filepath.Dir(relativePath))
		log.Infof("relativePath dir: %v", relativePath)
		log.Infof("save dir: %v", saveDir)

		// Skip template manifest file and the git config directory.
		if currFile == TmplManifestFile || strings.Contains(relativePath, gitConfigDir+PS) {
			log.Infof(msg.Stdout.SkipFile, relativePath)
			return
		}

		if inSkipArray(relativePath, tmplJson.Skip) { // Skip files in this list
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return
		}

		// skip the directory with replace files.
		if tmplJson.Replace != nil && tmplJson.Replace.Directory != "" && strings.Contains(relativePath, tmplJson.Replace.Directory) {
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return
		}

		// Make the subdirectories in the new savePath.
		err = os.MkdirAll(saveDir, dirMode)
		if err != nil || currFile == emptyFile {
			return
		}

		if tmplJson.Excludes != nil { // exclude from parsing, but copy as-is.
			// TODO: Replace with better method of comparing files.
			fileToCheck := strings.ReplaceAll(normSourcePath, normTplDir, "")
			log.Infof("fileToCheck: %q against excludes", fileToCheck)
			fileToCheck = strings.ReplaceAll(fileToCheck, PS, "")
			for _, exclude := range tmplJson.Excludes {
				fileToCheckB := strings.ReplaceAll(exclude, "\\", "")
				fileToCheckB = strings.ReplaceAll(exclude, "/", "")
				if fileToCheckB == fileToCheck {
					log.Infof("will copy as-is: %q", sourcePath)
					_, errC := copyToDir(sourcePath, saveDir, PS)
					return errC
				}
			}
		}

		sourcePath = replaceWith(relativePath, PS, sourcePath, normTplDir, tmplJson.Replace)

		rErr = parse(sourcePath, saveDir, vars)

		return
	})

	return
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
