package press

import (
	"regexp"
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
