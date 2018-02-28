package git_test

import (
	"testing"

	"github.com/apex/up/internal/git"
	"github.com/tj/assert"
)

func TestIsRepo(t *testing.T) {
	assert.True(t, git.IsRepo("."))
	assert.False(t, git.IsRepo("/tmp"))
}
