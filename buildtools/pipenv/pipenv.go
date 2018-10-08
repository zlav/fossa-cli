package pipenv

import (
	"encoding/json"

	"github.com/fossas/fossa-cli/errors"
	"github.com/fossas/fossa-cli/exec"
	"github.com/fossas/fossa-cli/pkg"
)

// TODO: Add fallback for Pipfile.Lock analysis in a situation where pipenv
// Is not present on the machine running fossa analyze

// DepTree is used to unmarshall the output from pipenv graph and store
// a object representing the dependcey tree
type DepTree struct {
	Package      string `json:"package_name"`
	Resolved     string `json:"installed_version"`
	Target       string `json:"required_version"`
	Dependencies []DepTree
}

// Deps returns the list of imports and associted package graph
// using the output of pipenv graph --json-tree
func Deps() ([]pkg.Import, map[pkg.ID]pkg.Package, error) {
	tree, err := getTree()
	if err != nil {
		return nil, nil, err
	}

	imports := importsFromTree(tree)
	graph := graphFromTree(tree)
	return imports, graph, nil
}

func getTree() ([]DepTree, error) {
	out, _, err := exec.Run(exec.Cmd{
		Name: "pipenv",
		Argv: []string{"graph", "--json-tree"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not run `pipenv graph`")
	}

	// Parse output.
	var tree []DepTree
	err = json.Unmarshal([]byte(out), &tree)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse `pipenv graph --json-tree` output")
	}
	return tree, nil
}

func importsFromTree(tree []DepTree) []pkg.Import {
	var imports []pkg.Import
	for _, dep := range tree {
		imports = append(imports, pkg.Import{
			Target: dep.Target,
			Resolved: pkg.ID{
				Type:     pkg.Python,
				Name:     dep.Package,
				Revision: dep.Resolved,
			},
		})
	}
	return imports
}

func graphFromTree(tree []DepTree) map[pkg.ID]pkg.Package {
	graph := make(map[pkg.ID]pkg.Package)
	for _, subtree := range tree {
		flattenTree(graph, subtree)
	}

	return graph
}

func flattenTree(graph map[pkg.ID]pkg.Package, tree DepTree) {
	for _, dep := range tree.Dependencies {
		// Construct ID.
		id := pkg.ID{
			Type:     pkg.Python,
			Name:     dep.Package,
			Revision: dep.Resolved,
		}
		// Don't process duplicates.
		_, ok := graph[id]
		if ok {
			continue
		}
		// Get direct imports.
		var imports []pkg.Import
		for _, i := range tree.Dependencies {
			imports = append(imports, pkg.Import{
				Resolved: pkg.ID{
					Type:     pkg.Python,
					Name:     i.Package,
					Revision: i.Resolved,
				},
			})
		}
		// Update map.
		graph[id] = pkg.Package{
			ID:      id,
			Imports: imports,
		}
		// Recurse in imports.
		flattenTree(graph, dep)
	}
}
