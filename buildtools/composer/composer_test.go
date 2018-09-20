package composer_test

import (
	"path/filepath"
	"testing"

	"github.com/fossas/fossa-cli/analyzers"
	"github.com/fossas/fossa-cli/analyzers/php"
	"github.com/fossas/fossa-cli/module"
	"github.com/fossas/fossa-cli/pkg"
	"github.com/stretchr/testify/assert"
)

func TestNoPHP(t *testing.T) {
	// Run analysis.
	m := module.Module{
		Name:        "fixture",
		Type:        pkg.Composer,
		BuildTarget: filepath.Join("testdata", "composer.json"),
		Dir:         "testdata",
	}

	a, err := analyzers.New(m)
	assert.NoError(t, err)
	assert.IsType(t, &php.Analyzer{}, a)

	deps, err := a.Analyze()
	assert.NoError(t, err)

	// Ensure no PHP dependencies.
	for _, dep := range deps.Direct {
		assert.NotEqual(t, "php", dep.Resolved.Name)
	}

	for id, dep := range deps.Transitive {
		assert.NotEqual(t, "php", id.Name)
		for _, i := range dep.Imports {
			assert.NotEqual(t, "php", i.Resolved.Name)
		}
	}
}
