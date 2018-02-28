// Package git provides GIT repo utilities.
package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// TODO: lookup root?... may not always have ./.git
// TODO: better IsDirty ... exit status or git describe --dirty=DIRTY... refactor tag only...?
// TODO: exchange Describe for tag...?
// TODO: cache lookup

// IsRepo returns true if dir is a git repo.
func IsRepo(dir string) bool {
	path := filepath.Join(dir, ".git")
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirty returns true if the dir is dirty.
func IsDirty(dir string) (bool, error) {
	bin, err := exec.LookPath("git")
	if err != nil {
		return false, errors.Wrap(err, "looking up git")
	}

	cmd := exec.Command(bin, "-C", dir, "status", "--porcelain")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Wrapf(err, "git/is_dirty: couldn't run command")
	}

	return len(stdout) > 0, nil
}

// Describe returns the git tag or sha.
func Describe(dir string) (string, error) {
	bin, err := exec.LookPath("git")
	if err != nil {
		return "", errors.Wrap(err, "looking up git")
	}

	cmd := exec.Command(bin, "-C", dir, "rev-parse", "--verify", "--short", "HEAD")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "executing")
	}

	return strings.TrimSpace(string(stdout)), nil
}
