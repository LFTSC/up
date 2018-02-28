// Package git provides GIT repo utilities.
package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// TODO: lookup root?... may not always have ./.git
// TODO: message?

// Errors.
var (
	ErrDirty = errors.New("repo is dirty")
)

// IsRepo returns true if dir is a git repo.
func IsRepo(dir string) bool {
	path := filepath.Join(dir, ".git")
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Describe returns the git tag or sha.
func Describe(dir string) (string, error) {
	bin, err := exec.LookPath("git")
	if err != nil {
		return "", errors.Wrap(err, "looking up git")
	}

	cmd := exec.Command(bin, "-C", dir, "describe", "--abbrev=0", "--dirty=DIRTY")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "executing")
	}

	if isDirty(out) {
		return "", ErrDirty
	}

	return string(bytes.TrimSpace(out)), nil
}

// isDirty returns true if the DIRTY mark is present.
func isDirty(b []byte) bool {
	return bytes.Contains(b, []byte("DIRTY"))
}
