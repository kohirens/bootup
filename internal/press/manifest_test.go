package press

import (
	"github.com/kohirens/tmplpress/internal/test"
	"testing"
)

func TestLoadAnswers(tester *testing.T) {
	var fixtures = []struct {
		name, file, want string
	}{
		{"goodJson", test.FixturesDir + PS + "answers-01.json", "value1"},
		{"badJson", test.FixturesDir + PS + "answers-02.json", ""},
	}

	fxtr := fixtures[0]
	tester.Run(fxtr.name, func(t *testing.T) {
		got, err := LoadAnswers(fxtr.file)

		if err != nil {
			t.Errorf("got error %v", err.Error())
		}

		if got.Placeholders["var1"] != fxtr.want {
			t.Errorf("got %q, want %q", got.Placeholders["var1"], fxtr.want)
		}
	})

	fxtr = fixtures[1]
	tester.Run(fxtr.name, func(t *testing.T) {
		_, err := LoadAnswers(fxtr.file)

		if err == nil {
			t.Error("got no error")
		}
	})
}
