package press

import (
	"github.com/kohirens/stdlib/test"
	"os"
	"testing"
)

const tmpDir = "tmp"
const fixtureDir = "testdata"

func TestMain(m *testing.M) {
	test.ResetDir(test.TmpDir, dirMode)
	// Run all tests
	exitCode := m.Run()
	// Clean up
	os.Exit(exitCode)
}
