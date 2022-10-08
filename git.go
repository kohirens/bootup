package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// gitClone Clone a repo from a path/URL to a local directory.
func gitClone(repoUri, repoDir, refName string) (string, string, error) {
	infof("branch to clone is %q", refName)
	infof("git clone %s", repoUri)

	var sco []byte
	var e1 error

	if isRemoteRepo(repoUri) {
		branchName := refName
		re := regexp.MustCompile("^refs/[^/]+/(.*)$")
		if re.MatchString(refName) {
			branchName = re.ReplaceAllString(refName, "${1}")
		}

		logf("cloning branch name: %v", branchName)

		// git clone --depth 1 --branch <tag_name> <repo_url>
		// NOTE: Branch cannot be a full ref but can be short ref name or a tag.
		sco, e1 = gitCmd(".", "clone", "--depth", "1", "--branch", branchName, repoUri, repoDir)
	} else {
		// git clone <repo_url>
		sco, e1 = gitCmd(".", "clone", repoUri, repoDir)
		if e1 != nil {
			return "", "", fmt.Errorf(errors.cloning, repoUri, e1.Error())
		}

		infof("clone output \n%s", sco)

		// git checkout <ref_name>
		sco3, e3 := gitCmd(repoDir, "checkout", "-b", refName)
		if e3 != nil {
			return "", "", fmt.Errorf(errors.cloning, repoUri, e1.Error(), sco3)
		}

		infof("clone output \n%s", sco3)
	}

	latestCommitHash, e2 := getLastCommitHash(repoDir)
	if e2 != nil {
		return "", "", fmt.Errorf(errors.gettingCommitHash, repoDir, e2.Error())
	}

	return repoDir, latestCommitHash, nil
}

// gitCheckout Open an existing repo and checkout commit by full ref-name
func gitCheckout(repoLocalPath, ref string) (string, string, error) {
	infof("pulling latest\n")
	_, e1 := gitCmd(repoLocalPath, "fetch", "--all", "-p")
	if e1 != nil {
		return "", "", fmt.Errorf(errors.gitFetchFailed, repoLocalPath, ref, e1.Error())
	}

	infof(messages.refInfo, ref)
	infof(messages.gitCheckout, ref)

	_, e2 := gitCmd(repoLocalPath, "checkout", ""+ref)
	if e2 != nil {
		return "", "", fmt.Errorf(errors.gitCheckoutFailed, e2.Error())
	}

	repoDir, e8 := filepath.Abs(repoLocalPath)
	if e8 != nil {
		return "", "", e8
	}

	latestCommitHash, e4 := getLastCommitHash(repoDir)
	if e4 != nil {
		return "", "", e4
	}

	return repoDir, latestCommitHash, nil
}

// gitCmd run a git command.
func gitCmd(repoPath string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Env = os.Environ()
	cmd.Dir = repoPath
	cmdStr := cmd.String()
	infof(messages.runningCommand, cmdStr)
	cmdOut, cmdErr := cmd.CombinedOutput()
	exitCode := cmd.ProcessState.ExitCode()

	if cmdErr != nil {
		return nil, fmt.Errorf(errors.runGitFailed, args, cmdErr.Error(), cmdOut)
	}

	if exitCode != 0 {
		return nil, fmt.Errorf(errors.gitExitErrCode, args, exitCode)
	}

	return cmdOut, nil
}

// getLastCommitHash Returns the HEAD commit hash.
func getLastCommitHash(repoDir string) (string, error) {
	latestCommitHash, e1 := gitCmd(repoDir, "rev-parse", "HEAD")
	if e1 != nil {
		return "", fmt.Errorf(errors.gettingCommitHash, repoDir, e1.Error())
	}

	return strings.Trim(string(latestCommitHash), "\n"), nil
}

// getLatestTag Will return the latest tag or an empty string from a repository.
func getLatestTag(repoDir string) (string, error) {
	tags, e1 := getRemoteTags(repoDir)
	if e1 != nil {
		return "", fmt.Errorf(errors.getLatestTag, repoDir, e1.Error())
	}

	return tags[0], nil
}

// getRemoteTags Get the remote tags on a repo using git ls-remote.
func getRemoteTags(repo string) ([]string, error) {
	// Even without cloning or fetching, you can check the list of tags on the upstream repo with git ls-remote:
	sco, e1 := gitCmd(repo, "ls-remote", "--sort=-version:refname", "--tags")
	if e1 != nil {
		return nil, fmt.Errorf(errors.getRemoteTags, e1.Error())
	}

	reTags := regexp.MustCompile("[a-f0-9]+\\s+refs/tags/(\\S+)")
	mat := reTags.FindAllSubmatch(sco, -1)
	if mat == nil {
		return nil, fmt.Errorf("%s", "no tags found")
	}

	ret := make([]string, len(mat))
	for i, v := range mat {
		dbugf(messages.remoteTagDbug1, string(v[1]))
		ret[i] = string(v[1])
	}

	return ret, nil
}

// getRepoDir extract a local dirname from a Git URL.
func getRepoDir(repoLocation, refName string) string {
	if len(repoLocation) < 1 {
		return repoLocation
	}

	baseName := filepath.Base(repoLocation)

	// trim .git from the end
	baseName = strings.TrimRight(baseName, ".git")

	// append ref, branch, or tag
	if refName != "" {
		baseName = baseName + "-" + strings.ReplaceAll(refName, "/", "-")
	}

	infof(messages.repoDir, baseName)

	return baseName
}

// getRepoDir extract a local dirname from a Git URL.
func isRemoteRepo(repoLocation string) bool {
	if len(repoLocation) < 1 {
		return false
	}
	// git@github.com:kohirens/tmpltoap.git
	// https://github.com/kohirens/tmpltoapp.git
	isGitUri := regexp.MustCompile("^(git|http|https)://.+$")
	if isGitUri.MatchString(repoLocation) {
		return true
	}

	return false
}
