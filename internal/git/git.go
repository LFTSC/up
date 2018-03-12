// Package git provides GIT repo utilities.
package git

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

// TODO: support git SHA... currently errors unless there's a tag

// Errors.
var (
	ErrDirty  = errors.New("git repo is dirty")
	ErrNoRepo = errors.New("git repo not found")
	ErrLookup = errors.New("git is not installed")
)

// IsRepo returns true if dir is a git repo.
func IsRepo(dir string) bool {
	bin, err := exec.LookPath("git")
	if err != nil {
		return false
	}

	cmd := exec.Command(bin, "status")
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return false
	}

	return cmd.ProcessState.Success()
}

// Describe returns the git tag or sha.
func Describe(dir string) (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0", "--dirty=DIRTY")
	cmd.Dir = dir
	return output(cmd)
}

// Author returns the author of HEAD.
func Author(dir string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--pretty=format:%an")
	cmd.Dir = dir
	return output(cmd)
}

// output returns GIT command output with error normalization.
func output(cmd *exec.Cmd) (string, error) {
	switch out, err := cmd.CombinedOutput(); {
	case err == exec.ErrNotFound:
		return "", ErrLookup
	case bytes.Contains(out, []byte("Not a git repository")):
		return "", ErrNoRepo
	case bytes.Contains(out, []byte("DIRTY")):
		return "", ErrDirty
	case err != nil:
		return "", errors.New(string(out))
	default:
		return string(bytes.TrimSpace(out)), nil
	}
}
