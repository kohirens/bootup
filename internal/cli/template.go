package cli

// TODO: move to pkg/press
import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/stdlib/path"
	"github.com/kohirens/tmpltoapp/internal/msg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	MaxTplSize       = 1e+7
	EmptyFile        = ".empty"
	GitConfigDir     = ".git"
	TmplManifestFile = "template.json" // TODO: BREAKING Rename to tmplpress.json
)

type AnswersJson struct {
	Placeholders tmplVars `json:"placeholders"`
}

// Client specify the methods required by an HTTP client
type Client interface {
	Get(url string) (*http.Response, error)
	Head(url string) (*http.Response, error)
}

type TmplJson struct {
	Excludes     []string    `json:"excludes"`
	Placeholders tmplVars    `json:"placeholders"`
	Skip         []string    `json:"skip"`
	Replace      ReplaceWith `json:"replace"`
	Validation   []validator `json:"validation"`
	Version      string      `json:"version"`
}

type ReplaceWith struct {
	Directory string   `json:"directory"`
	Files     []string `json:"files"`
}

type tmplVars map[string]string

// GetPlaceholderInput Checks for any missing placeholder values waits for their input from the CLI.
func GetPlaceholderInput(placeholders *TmplJson, tmplValues *tmplVars, r *os.File, defaultVal string) error {
	numPlaceholder := len(placeholders.Placeholders)
	numValues := len(*tmplValues)

	//log.Logf(msg.Stdout.PlaceholderAnswerStat, numPlaceholder)

	if numPlaceholder == numValues {
		return nil
	}

	//log.Logf(msg.Stdout.ProvideValues)

	tVals := *tmplValues
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

// copyToDir copy a file to a directory
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

func NewAnswerJson() *AnswersJson {
	return &AnswersJson{
		Placeholders: make(tmplVars),
	}
}

// Deprecated: Do NOT USE, will be removed.
func Extract(archivePath string) (string, error) {
	tmplDir := ""
	zipParentDir := ""
	dest := strings.ReplaceAll(archivePath, ".zip", "")
	// Get resource to zip archive.
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return tmplDir, fmt.Errorf("could not open archive %q, error: %v", archivePath, err.Error())
	}

	err = os.MkdirAll(dest, DirMode)
	if err != nil {
		return tmplDir, fmt.Errorf("could not write dest %q, error: %v", dest, err.Error())
	}

	log.Infof("extracting %v to %v\n", archivePath, dest)
	for _, file := range archive.File {
		sourceFile, fErr := file.Open()
		if fErr != nil {
			return tmplDir, fmt.Errorf("failed to Extract archive %q to dest %q, error: %v", archivePath, dest, file.Name)
		}

		extractionDir := filepath.Join(dest, file.Name)
		// trying to figure out the
		if zipParentDir == "" {
			// TODO: Document the fact that template archives MUST be zip format and contain all template files in a single directory at the root of the zip.
			zipParentDir = extractionDir
		}

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(extractionDir, filepath.Clean(dest)+PS) {
			return tmplDir, fmt.Errorf("illegal file path: %s", extractionDir)
		}

		if file.FileInfo().IsDir() {
			ferr := os.MkdirAll(extractionDir, file.Mode())
			if ferr != nil {
				return tmplDir, ferr
			}
		} else {
			dh, ferr := os.OpenFile(extractionDir, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())

			if ferr != nil {
				return tmplDir, ferr
			}

			_, ferr = io.Copy(dh, sourceFile)
			if ferr != nil {
				return tmplDir, ferr
			}

			ferr = dh.Close()
			if ferr != nil {
				panic(ferr)
			}
		}

		fErr = sourceFile.Close()
		if fErr != nil {
			return tmplDir, fmt.Errorf("unsuccessful extracting archive %q, error: %v", archivePath, fErr.Error())
		}
	}

	err = archive.Close()
	tmplDir = zipParentDir
	log.Dbugf("zipParentDir = %v", zipParentDir)

	return tmplDir, nil
}

// parse a file as a Go template.
func parse(tplFile, dstDir string, vars tmplVars) error {
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

// Press Recursively walk a directory parsing all files along the way as Go templates.
// Deprecated: See press.FindTemplates
func Press(tplDir, outDir string, vars tmplVars, fec *stdlib.FileExtChecker, tmplJson *TmplJson) (err error) {
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
		if fi.Size() > MaxTplSize {
			rErr = fmt.Errorf(msg.Stderr.FileTooBig, MaxTplSize)
			return
		}

		currFile := filepath.Base(sourcePath)
		if currFile != EmptyFile && !fec.IsValid(sourcePath) { // Skip files by extension; Use an exclusion list, include every file by default.
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
		if currFile == TmplManifestFile || strings.Contains(relativePath, GitConfigDir+PS) {
			log.Infof(msg.Stdout.SkipFile, relativePath)
			return
		}

		if inSkipArray(relativePath, tmplJson.Skip) { // Skip files in this list
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return
		}

		// skip the directory with replace files.
		if tmplJson.Replace.Directory != "" && strings.Contains(relativePath, tmplJson.Replace.Directory) {
			log.Infof(msg.Stdout.SkipFile, sourcePath)
			return
		}

		// Make the subdirectories in the new savePath.
		err = os.MkdirAll(saveDir, DirMode)
		if err != nil || currFile == EmptyFile {
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

// ReadTemplateJson read variables needed from the template.json file.
func ReadTemplateJson(filePath string) (*TmplJson, error) {
	log.Dbugf("\ntemplate manifest path: %q\n", filePath)

	// Verify the TMPL_MANIFEST file is present.
	if !path.Exist(filePath) {
		return nil, fmt.Errorf(msg.Stderr.TmplManifest404, TmplManifestFile)
	}

	content, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		return nil, err1
	}

	log.Infof("content = %s \n", content)

	q := TmplJson{}
	if err2 := json.Unmarshal(content, &q); err2 != nil {
		return nil, err2
	}

	log.Dbugf("TmplJson.Version = %v", q.Version)
	if q.Version == "" {
		return nil, fmt.Errorf("missing the Version propery in template.json")
	}

	log.Dbugf("TmplJson.Placeholders = %v", len(q.Placeholders))
	if q.Placeholders == nil {
		return nil, fmt.Errorf("missing the placeholders propery in template.json")
	}

	return &q, nil
}

func ShowAllPlaceholderValues(placeholders *TmplJson, tmplValues *tmplVars) {
	tVals := *tmplValues
	log.Logf("the following values have been provided\n")
	for placeholder := range placeholders.Placeholders {
		log.Logf(msg.Stdout.PlaceholderAnswer, placeholder, tVals[placeholder])
	}
}

func inSkipArray(p string, skips []string) bool {
	for _, skip := range skips {
		skip = path.Normalize(skip)
		if strings.Contains(p, skip) {
			return true
		}
	}
	return false
}

// replaceWith Replace the current file with another.
func replaceWith(cf, ps, sp, tmplRoot string, replace ReplaceWith) string {
	for _, fileMap := range replace.Files {
		fileMap = path.Normalize(fileMap)
		fileAry := strings.Split(fileMap, ":")

		if strings.Contains(cf, fileAry[1]) {
			if cf == fileAry[1] {
				return tmplRoot + ps + replace.Directory + ps + fileAry[0]
			}
			// match by prefix

			return tmplRoot + ps + replace.Directory + ps + fileAry[0] + strings.Replace(cf, fileAry[0], "", 1)
		}
	}

	return sp
}
