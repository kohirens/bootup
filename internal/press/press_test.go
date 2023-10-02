package press

import (
	"github.com/kohirens/tmpltoapp/internal/cli"
	"github.com/kohirens/tmpltoapp/internal/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	_ = os.RemoveAll(test.TmpDir)
	// Set up a temporary dir for generate files
	_ = os.Mkdir(test.TmpDir, cli.DirMode) // set up a temporary dir for generate files
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}
