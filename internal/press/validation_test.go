package press

import (
	"testing"
)

func TestValidateAlphaNumeric(t *testing.T) {
	rule := "alphaNumeric"
	testCases := []struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		{
			"letters",
			"abc",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"hyphen",
			"a-bc",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"underscore",
			"a_bc",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"numbers",
			"123",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"lettersAndNumbers",
			"acb123",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"specialChars",
			"*.(#",
			"var1",
			[]validator{{
				expression: "",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := validate(tc.ui, tc.ph, tc.v)
			if got != tc.want {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestValidateRegExp(t *testing.T) {
	rule := "regExp"
	testCases := []struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		{
			"compilesAndValidInput",
			"abc",
			"var1",
			[]validator{{
				expression: "[a-z]",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"compilesAndInvalidInput",
			"ABC",
			"var1",
			[]validator{{
				expression: "[a-z]",
				fields:     []string{"var1"},
				rule:       rule,
				message:    "var1 failed to validate",
			}},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := validate(tc.ui, tc.ph, tc.v)
			if got != tc.want {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestValidateRegExpCompileError(t *testing.T) {
	rule := "regExp"

	tc := struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		"doesNotCompiles",
		"ABC",
		"var1",
		[]validator{{
			expression: "[a-z",
			fields:     []string{"var1"},
			rule:       rule,
			message:    "var1 failed to validate",
		}},
		false,
	}

	got, e := validate(tc.ui, tc.ph, tc.v)
	if got != tc.want {
		t.Errorf("got %v want %v", got, tc.want)
	}
	if e == nil {
		t.Errorf("got %v want an error", e)
	}
}
