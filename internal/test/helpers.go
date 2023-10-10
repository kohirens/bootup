package test

import (
	"fmt"
	path2 "github.com/kohirens/stdlib/path"
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

// SetupARepository set up a repository.
// Deprecated see stdlib.git.CloneFromBundle
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
	if !path2.Exist(srcRepo) {
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
			fmt.Println("executed now")
		}
	default:
		oldHome, _ := os.LookupEnv("HOME")
		_ = os.Setenv("HOME", dir)
		return func() {
			_ = os.Setenv("HOME", oldHome)
			fmt.Println("executed now")
		}
	}
}

func VerboseSubCmdOut(stdo []byte, stde error) ([]byte, error) {
	// Debug
	if testing.Verbose() {
		fmt.Print("\nBEGIN sub-command\n")
		fmt.Printf("stdout:\n%s\n", stdo)

		if stde != nil {
			fmt.Printf("stderr:\n%v\n", stde.Error())
		}

		fmt.Print("\nEND sub-command\n\n")
	}

	return stdo, stde
}
