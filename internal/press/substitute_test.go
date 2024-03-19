package press

import (
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/stdlib/test"
	"os"
	"strings"
	"testing"
)

func TestCopyDirToDir(runner *testing.T) {
	td := test.TmpDir + PS + runner.Name()
	_ = os.MkdirAll(td, 0774)

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
		want    []string
	}{
		{
			"copy all files 01",
			test.FixtureDir + PS + "dir-to-dir-01",
			td,
			false,
			[]string{td + PS + "README.md"},
		},
		{
			"copy all files 02",
			test.FixtureDir + PS + "dir-to-dir-02",
			td,
			false,
			[]string{
				td + PS + "README.md",
				td + PS + "sub-01" + PS + "README.md",
				td + PS + "sub-02" + PS + "file-01.txt",
				td + PS + "sub-02" + PS + "file-02.txt",
				td + PS + "sub-03" + PS + "README.md",
				td + PS + "sub-03" + PS + "sub-sub-04" + PS + "README-04.md",
			},
		},
	}

	for _, tt := range tests {
		runner.Run(tt.name, func(t *testing.T) {
			if err := Substitute(tt.src, tt.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyDirToDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, f := range tt.want {
				if !fsio.Exist(f) {
					t.Errorf("file not found %v", f)
				}
			}
		})
	}
}

// Verify that the Substitute does indeed overwrite any file in the destination.
func TestCopyDirToDirDstOverwrite(runner *testing.T) {
	td := test.TmpDir + PS + runner.Name()
	_ = os.MkdirAll(td, 0744)

	f := "dir-to-dir-03"

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			"success",
			test.FixtureDir + PS + f,
			td,
			false,
		},
	}
	for _, tt := range tests {
		runner.Run(tt.name, func(t *testing.T) {
			d := git.CloneFromBundle(f, td, test.FixtureDir, PS)

			if err := Substitute(tt.src, d); (err != nil) != tt.wantErr {
				t.Errorf("CopyDirToDir() error = %v, wantErr %v", err, tt.wantErr)
			}

			b, e1 := os.ReadFile(d + PS + ".config/config.yml")
			if e1 != nil {
				t.Errorf("failed to read file: %v", e1.Error())
			}

			if !strings.Contains(string(b), "2.1") {
				t.Errorf("did not get expected content from file.")
			}
		})
	}
}
