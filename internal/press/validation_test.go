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
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"hyphen",
			"a-bc",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"underscore",
			"a_bc",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"numbers",
			"123",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"lettersAndNumbers",
			"acb123",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"specialChars",
			"*.(#",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
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

func TestValidateBoolean(t *testing.T) {
	rule := "bool"
	testCases := []struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		{
			"true",
			"true",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"false",
			"false",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"false",
			"false",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"badInput",
			"1",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
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

func TestValidateInt(t *testing.T) {
	rule := "int"
	testCases := []struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		{
			"number",
			"21",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"negative",
			"-32",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"NaN",
			"NaN",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"decimal",
			"1.00",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
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

func TestValidateUInt(t *testing.T) {
	rule := "unsigned"
	testCases := []struct {
		name string
		ui   string
		ph   string
		v    []validator
		want bool
	}{
		{
			"number",
			"21",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"negative",
			"-32",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"NaN",
			"NaN",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			false,
		},
		{
			"decimal",
			"1.00",
			"var1",
			[]validator{{
				Expression: "",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
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
				Expression: "[a-z]",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
			}},
			true,
		},
		{
			"compilesAndInvalidInput",
			"ABC",
			"var1",
			[]validator{{
				Expression: "[a-z]",
				Fields:     []string{"var1"},
				Rule:       rule,
				Message:    "var1 failed to validate",
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
			Expression: "[a-z",
			Fields:     []string{"var1"},
			Rule:       rule,
			Message:    "var1 failed to validate",
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

func Test_checkFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{"case-01", "file-01", false},
		{"case-01", "/file-01", false},
		{"case-01", "/\u0061\u0300-file-02", false},
		{"case-01", "\u0300-file-02", true},
		{"case-01", "/\u00E0-file-03", false},
		{"case-01", "[!a-z]-0*.html?", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkFilename(tt.filename); (err != nil) != tt.wantErr {
				t.Errorf("checkFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunValidate(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		cmd      string
		wantErr  bool
	}{
		{"case-1", fixtureDir + PS + "template-01.json", "validate", true},
		{"case-2", fixtureDir + PS + "template-02.json", "validate", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateManifest(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkSubstitute(t *testing.T) {

	tests := []struct {
		name     string
		filename string
		dir      string
		wantErr  bool
	}{
		{"case-1", fixtureDir + PS + t.Name(), "replace", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkSubstitute(tt.filename, tt.dir); (err != nil) != tt.wantErr {
				t.Errorf("checkSubstitute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkValidationRules(t *testing.T) {
	tests := []struct {
		name         string
		placeholders map[string]string
		rules        []*validator
		wantErr      bool
	}{
		{
			"non-existing-placeholder",
			map[string]string{"var1": ""},
			[]*validator{
				{
					Fields:  []string{"var1", "var2"},
					Rule:    "alphaNumeric",
					Message: "var1 failed to validate",
				},
			},
			true,
		},
		{
			"empty-regexp",
			map[string]string{"var1": ""},
			[]*validator{
				{
					Expression: "",
					Fields:     []string{"var1"},
					Rule:       "regExp",
					Message:    "var1 failed to validate",
				},
			},
			true,
		},
		{
			"invalid-regexp",
			map[string]string{"var1": ""},
			[]*validator{
				{
					Expression: "[a-z",
					Fields:     []string{"var1"},
					Rule:       "regExp",
					Message:    "var1 failed to validate",
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkValidationRules(tt.placeholders, tt.rules); (err != nil) != tt.wantErr {
				t.Errorf("checkValidationRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
