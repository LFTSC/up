// Package git provides GIT repo utilities.
package git

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

// Errors.
var (
	ErrDirty = errors.New("repo is dirty")
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
	bin, err := exec.LookPath("git")
	if err != nil {
		return "", errors.Wrap(err, "looking up git")
	}

	cmd := exec.Command(bin, "describe", "--abbrev=0", "--dirty=DIRTY")
	cmd.Dir = dir

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
