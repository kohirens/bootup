package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	FixturesDir = "testdata"
	TmpDir      = "tmp"
)

// TmpSetParentDataDir set the LOCALAPPDATA or HOME environment var for a unit test.
func TmpSetParentDataDir(d string) func() {
	dir, err := filepath.Abs(d)
	if err != nil {
		panic(fmt.Sprintf("failed to get path to %q for unit test", TmpDir))
	}

	// Set the app data dir to the local test tmp.
	switch runtime.GOOS {
	case "windows":
		appDataDir, _ := os.LookupEnv("APPDATA")
		cacheDataDir, _ := os.LookupEnv("LOCALAPPDATA")

		if e := os.Setenv("APPDATA", dir); e != nil {
			panic("failed to set LOCALAPPDATA for unit test")
		}

		if e := os.Setenv("LOCALAPPDATA", dir); e != nil {
			panic("failed to set LOCALAPPDATA for unit test")
		}

		return func() {
			_ = os.Setenv("APPDATA", appDataDir)
			_ = os.Setenv("LOCALAPPDATA", cacheDataDir)
		}
	default:
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", dir)
		return func() {
			_ = os.Setenv("HOME", oldHome)
		}
	}
}
