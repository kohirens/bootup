package press

import (
	"github.com/kohirens/stdlib/test"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	test.ResetDir(test.TmpDir, dirMode)
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}
