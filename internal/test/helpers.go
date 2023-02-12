package test

import (
	"fmt"
	"github.com/kohirens/stdlib"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const (
	FixturesDir = "testdata"
	SubCmdFlags = "SUB_CMD_FLAGS"
	TmpDir      = "tmp"
	Remotes     = TmpDir + PS + "remotes"
	PS          = string(os.PathSeparator)
)

// Silencer return a function that prevents output during a test run.
func Silencer() func() {
	// Abort in verbose mode.
	if testing.Verbose() {
		return func() {}
	}
	null, _ := os.Open(os.DevNull)
	sOut := os.Stdout
	sErr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sOut
		os.Stderr = sErr
		log.SetOutput(os.Stderr)
	}
}

func SetupARepository(bundleName, tmpDir, bundleDir, ps string) string {
	repoPath := tmpDir + ps + bundleName

	// It may have already been unbundled.
	fileInfo, err1 := os.Stat(repoPath)
	if (err1 == nil && fileInfo.IsDir()) || os.IsExist(err1) {
		absPath, e2 := filepath.Abs(repoPath)
		if e2 == nil {
			return absPath
		}
		return repoPath
	}

	wd, e := os.Getwd()
	if e != nil {
		panic(fmt.Sprintf("%v failed to get working directory", e.Error()))
	}

	srcRepo := wd + ps + bundleDir + ps + bundleName + ".bundle"
	// It may not exist.
	if !stdlib.PathExist(srcRepo) {
		panic(fmt.Sprintf("%v bundle not found", srcRepo))
	}

	cmd := exec.Command("git", "clone", "-b", "main", srcRepo, repoPath)
	_, _ = cmd.CombinedOutput()
	if ec := cmd.ProcessState.ExitCode(); ec != 0 {
		log.Panicf("error un-bundling %q to %q for a unit test", srcRepo, repoPath)
	}

	absPath, e2 := filepath.Abs(repoPath)
	if e2 != nil {
		panic(e2.Error())
	}

	fmt.Printf("\nabsPath = %v\n", absPath)
	return absPath
}

// GetTestBinCmd return a command to run the test binary in a sub-process, passing it flags as fixtures to produce expected output; `TestMain`, will be run automatically.
func GetTestBinCmd(args []string) *exec.Cmd {
	// call the generated test binary directly
	// Have it the function runAppMain.
	cmd := exec.Command(os.Args[0])
	// Run in the context of the source directory.
	_, filename, _, _ := runtime.Caller(0)
	cmd.Dir = path.Dir(filename)
	// Set an environment variable
	// 1. Only exist for the life of the test that calls this function.
	// 2. Passes arguments/flag to your app
	// 3. Lets TestMain know when to run the main function.
	subEnvVar := SubCmdFlags + "=" + strings.Join(args, " ")
	cmd.Env = append(os.Environ(), subEnvVar)

	return cmd
}

// TmpSetParentDataDir set the LOCALAPPDATA or HOME environment var for a unit test.
func TmpSetParentDataDir(d string) func() {
	dir, err := filepath.Abs(d)
	if err != nil {
		panic(fmt.Sprintf("failed to get path to %q for unit test", TmpDir))
	}

	// Set the app data dir to the local test tmp.
	if runtime.GOOS == "windows" {
		oldAppData, _ := os.LookupEnv("LOCALAPPDATA")

		if e := os.Setenv("LOCALAPPDATA", dir); e != nil {
			panic("failed to set LOCALAPPDATA for unit test")
		}

		return func() {
			_ = os.Setenv("LOCALAPPDATA", oldAppData)
		}
	} else {
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", dir)
		return func() {
			_ = os.Setenv("HOME", oldHome)
		}
	}
}
