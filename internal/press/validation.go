package press

import (
	"fmt"
	"github.com/kohirens/tmplpress/internal/msg"
	"regexp"
	"strconv"
)

type validator struct {
	expression string
	fields     []string
	rule       string
	message    string
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
