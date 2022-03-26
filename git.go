package main

import (
	"fmt"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/kohirens/stdlib"
	"os"
)

// gitClone Clone a repo from a path/URL to a local directory.
func gitClone(repoLocation, outPath, branchName string) (string, error) {
	verboseF(verboseLvlInfo, "git clone %s", repoLocation)

	// TODO: checkout by branch or tag
	options := &git.CloneOptions{
		URL:          repoLocation,
		Progress:     os.Stdout,
		SingleBranch: false,
		Depth:        1,
	}

	// Set an optional branch name, otherwise default to head.
	if branchName != "" {
		options.ReferenceName = plumbing.NewBranchReferenceName(branchName)
	}

	repo, e1 := git.PlainClone(outPath, false, options)

	if e1 != nil {
		return "", e1
	}

	logs, e2 := repo.Log(&git.LogOptions{
		All: true,
	})

	if e2 != nil {
		return "", e2
	}

	commit, e3 := logs.Next()

	if e3 != nil {
		return "", e3
	}

	if stdlib.DirExist(outPath) != true {
		fmt.Errorf(errs.tmplOutput)
	}

	return commit.Hash.String(), nil
}
