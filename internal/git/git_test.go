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

func TestDescribe(t *testing.T) {
	s, err := git.Describe("/tmp")
	assert.EqualError(t, err, `git repo not found`)
	assert.Empty(t, s)
}

func TestAuthor(t *testing.T) {
	s, err := git.Author(".")
	assert.NoError(t, err)
	assert.Equal(t, "TJ Holowaychuk", s)
}
