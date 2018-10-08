package pipenv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDependencies(t *testing.T) {
	man := manifest{Ignored: []string{"apple", "orange*"}}

	// File listed in ignored list is ignored.
	valid := man.isIgnored("apple")
	assert.Equal(t, valid, true)

	// Wildcard entry properly ignores its own package.
	valid = man.isIgnored("orange")
	assert.Equal(t, valid, true)

	// Wildcard entry properly ignores other packages.
	valid = man.isIgnored("orange/blood")
	assert.Equal(t, valid, true)

	// File not listed in ignored list is not ignored.
	valid = man.isIgnored("apple/fuji")
	assert.Equal(t, valid, false)
}

func TestImportsFromDependencies(t *testing.T) {
}

func TestGraphFromDependencies(t *testing.T) {
}
