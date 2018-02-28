package git_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/apex/up/internal/git"
	"github.com/tj/assert"
)

func TestIsRepo(t *testing.T) {
	assert.True(t, git.IsRepo("."))
	assert.False(t, git.IsRepo("/tmp"))
}

func TestDescribe(t *testing.T) {
	s, err := git.Describe(filepath.Join("..", ".."))
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
	assert.True(t, strings.HasPrefix(s, "v"), "should have 'v' prefix")
}
