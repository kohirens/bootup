package press

import (
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/tmplpress/internal/msg"
	"golang.org/x/mod/semver"
	"path/filepath"
	"regexp"
	"strconv"
)

type validator struct {
	expression string
	fields     []string
	rule       string
	message    string
}

// ValidateManifest Read a template manifest and report any errors. This is a
// quality-of-life tool for template designers.
func ValidateManifest(aFile string) error {
	// check for existing template manifest and load it
	tm, e1 := ReadTemplateJson(aFile)
	if e1 != nil {
		return fmt.Errorf(msg.Stderr.CannotReadFile, aFile, e1.Error())
	}

	if e := checkVersion(tm.Version); e != nil {
		return e
	}

	if e := checkCopyAsIs(tm.CopyAsIs); e != nil {
		return e
	}

	if len(tm.EmptyDirFile) > 1 {
		if e := checkFilename(tm.EmptyDirFile); e != nil {
			return fmt.Errorf(msg.Stderr.EmptyDirFilename, aFile, e.Error())
		}
	}

	if e := checkVarName(tm.Placeholders); e != nil {
		return fmt.Errorf(msg.Stderr.PlaceholdersProperty, aFile, e.Error())
	}

	if e := checkFilePatterns(tm.Skip); e != nil {
		return fmt.Errorf(msg.Stderr.CannotReadFile, aFile, e.Error())
	}

	if e := checkSubstitute(aFile, tm.Substitute); e != nil {
		return fmt.Errorf(msg.Stderr.CannotReadFile, aFile, e.Error())
	}

	if e := checkValidationRules(tm.Placeholders, tm.Validation); e != nil {
		return fmt.Errorf(msg.Stderr.CannotReadFile, aFile, e.Error())
	}

	return nil
}

func checkCopyAsIs(filePatterns []string) error {
	if e := checkFilePatterns(filePatterns); e != nil {
		return e
	}

	return nil
}

// checkFilename Verify a filename is valid. Make use of Unicode for better
// language compatability. See https://www.regular-expressions.info/unicode.html
func checkFilename(filename string) error {
	re := regexp.MustCompile(`^([\p{L}\p{N}\-\\_./]|\P{M}\p{M}*)+$`)

	if !re.MatchString(filename) {
		return fmt.Errorf(msg.Stderr.Filename, filename)
	}

	return nil
}

func checkFilePatterns(files []string) error {
	for _, aFile := range files {
		if e := checkFilename(aFile); e != nil {
			return e
		}
	}

	return nil
}

func checkSubstitute(filename, dir string) error {
	substituteDir := filepath.Dir(filename) + PS + dir

	log.Infof("checking for substiture dir %v", substituteDir)

	if dir != "" && !fsio.Exist(substituteDir) {
		return fmt.Errorf(msg.Stderr.NoDir, substituteDir)
	}

	return nil
}

// checkValidationRules Verify rules apply and are of some correctness.
// 1. Each rule maps to existing placeholders.
// 2. Each regex rule will compile.
func checkValidationRules(placeholders map[string]string, validators []*validator) error {
	for _, vldtr := range validators {
		// verify each field is a placeholder.
		for _, name := range vldtr.fields {
			_, ok := placeholders[name]
			if !ok {
				return fmt.Errorf(msg.Stderr.NoPlaceholder, name)
			}
		}

		if vldtr.rule == "regExp" {
			if vldtr.expression == "" {
				return fmt.Errorf(msg.Stderr.EmptyRegExp, vldtr.expression)
			}

			_, e := regexp.Compile(vldtr.expression)
			if e != nil {
				return fmt.Errorf(msg.Stderr.InvalidRegExp, vldtr.expression, e.Error())
			}
		}
	}

	return nil
}

func checkVarName(vars map[string]string) error {
	return nil
}

func checkVersion(semantic string) error {
	ver := "v" + semantic
	if !semver.IsValid(ver) {
		return fmt.Errorf(msg.Stderr.MissingTmplJsonVersion)
	}

	c := semver.Compare(ver, "v"+schemaVersion)
	if c == 1 {
		return fmt.Errorf(msg.Stderr.InvalidManifestVersion, semantic)
	}

	return nil
}

// findValidator locate the validator for a placeholder
func findValidator(placeholder string, validators []validator) (validator, bool) {
	var val validator
	found := false

	// locate the validator
	for _, x := range validators {
		for _, y := range x.fields {
			if y == placeholder {
				val = x
				found = true
			}
		}
	}

	return val, found
}

func isBoolean(userInput string) (bool, error) {
	if userInput != "true" && userInput != "false" {
		//return false, fmt.Errorf(msg.Stderr.ParseBool, userInput)
		return false, nil
	}

	return true, nil
}

func isInt(userInput string) (bool, error) {
	_, e := strconv.ParseInt(userInput, 10, 64)
	if e != nil {
		return false, fmt.Errorf(msg.Stderr.ParseInt, userInput, e.Error())
	}

	return true, nil
}

func isUInt(userInput string) (bool, error) {
	_, e := strconv.ParseUint(userInput, 10, 64)
	if e != nil {
		//return false, fmt.Errorf(msg.Stderr.ParseInt, userInput, e.Error())
		return false, nil
	}

	return true, nil
}

func runRegex(expression, userInput string) (bool, error) {
	re, e := regexp.Compile(expression)
	if e != nil {
		return false, fmt.Errorf(msg.Stderr.InvalidRegExp, expression, e.Error())
	}

	return re.MatchString(userInput), nil
}

// validate user input for placeholders
func validate(userInput, placeholder string, validators []validator) (bool, error) {
	val, found := findValidator(placeholder, validators)

	if found { // perform validation
		switch val.rule {
		case "alphaNumeric":
			re := regexp.MustCompile("^[a-zA-Z0-9]+$")
			return re.MatchString(userInput), nil
		case "bool":
			return isBoolean(userInput)
		case "int":
			return isInt(userInput)
		case "unsigned":
			return isUInt(userInput)
		case "regExp":
			re, e := regexp.Compile(val.expression)
			if e == nil {
				return re.MatchString(userInput), nil
			}
			return false, e
		}
	}

	return false, nil
}
