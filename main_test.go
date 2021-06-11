package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	fixturesDir = "testdata"
	testTmp     = "tmp"
)

func TestMain(m *testing.M) {
	programName = filepath.Base(os.Args[0])
	// call flag.Parse() here if TestMain uses flags
	os.RemoveAll(testTmp)
	// Set up a temporary dir for generate files
	os.Mkdir(testTmp, DIR_MODE) // set up a temporary dir for generate files
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}

// SubCmdFlags space separated list of command line flags.
const SubCmdFlags = "SUB_CMD_FLAGS"

func TestCallingMain(tester *testing.T) {
	// This was adapted from https://golang.org/src/flag/flag_test.go; line 596-657 at the time.
	// This is called recursively, because we will have this test call itself
	// in a sub-command with the environment variable `GO_CHILD_FLAG` set.
	// Note that a call to `main()` MUST exit or you'll spin out of control.
	if os.Getenv(SubCmdFlags) != "" {
		// We're in the test binary, so test flags are set, lets reset it so
		// so that only the program is set
		// and whatever flags we want.
		args := strings.Split(os.Getenv(SubCmdFlags), " ")
		os.Args = append([]string{os.Args[0]}, args...)

		// Anything you print here will be passed back to the cmd.Stderr and
		// cmd.Stdout below, for example:
		fmt.Printf("os args = %v\n", os.Args)

		// Strange, I was expecting a need to manually call the conde in
		// `init()`,but that seem to happen automatically. So yet more I have learn.
		main()
	}

	var tests = []struct {
		name string
		want int
		args []string
	}{
		{"noArgs", 0, []string{"-v"}},
		{"noAppPath", 0, []string{"-h"}},
		{"allGood", 0, []string{"-a", fixturesDir + PS + "answers-parse-dir-02.json", "-t", fixturesDir + PS + "parse-dir-02", "-p", testTmp + PS + "app-parse-dir-02"}},
	}

	for _, test := range tests {
		tester.Run(test.name, func(t *testing.T) {
			cmd := runMain(tester.Name(), test.args)

			out, sce := cmd.CombinedOutput()

			// exec code.
			got := cmd.ProcessState.ExitCode()

			if got != test.want {
				t.Errorf("got %q, want %q", got, test.want)
			}

			if sce != nil {
				fmt.Printf("\nBEGIN sub-command\nstdout:\n%v\n\n", string(out))
				fmt.Printf("stderr:\n%v\n", sce.Error())
				fmt.Print("\nEND sub-command\n\n")
			}
		})
	}
}

func runMain(testFunc string, args []string) *exec.Cmd {
	subEnvVar := SubCmdFlags + "=" + strings.Join(args, " ")

	// Run the test binary and tell it to run just this test with environment set.
	cmd := exec.Command(os.Args[0], "-test.run", testFunc)
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}
