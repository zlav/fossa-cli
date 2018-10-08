package pipenv

import (
	"encoding/json"

	"github.com/fossas/fossa-cli/errors"
	"github.com/fossas/fossa-cli/exec"
	"github.com/fossas/fossa-cli/pkg"
)

// dependency is used to unmarshall the output from pipenv graph and store
// an object representing an imported dependency as well as its
// child dependencies.
type dependency struct {
	Package      string `json:"package_name"`
	Resolved     string `json:"installed_version"`
	Target       string `json:"required_version"`
	Dependencies []dependency
}

// Deps returns the list of imports and associted package graph
// using the output of pipenv graph --json-tree.
func Deps() ([]pkg.Import, map[pkg.ID]pkg.Package, error) {
	deps, err := getDependencies()
	if err != nil {
		return nil, nil, err
	}

	imports := importsFromDependencies(deps)
	graph := graphFromDependencies(deps)
	return imports, graph, nil
}

func getDependencies() ([]dependency, error) {
	out, _, err := exec.Run(exec.Cmd{
		Name: "pipenv",
		Argv: []string{"graph", "--json-tree"},
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not run `pipenv graph`")
	}

	var depList []dependency
	err = json.Unmarshal([]byte(out), &depList)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse `pipenv graph --json-tree` output")
	}
	return depList, nil
}

func importsFromDependencies(depList []dependency) []pkg.Import {
	var imports []pkg.Import
	for _, dep := range depList {
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

func graphFromDependencies(depList []dependency) map[pkg.ID]pkg.Package {
	graph := make(map[pkg.ID]pkg.Package)
	for _, dep := range depList {
		id := pkg.ID{
			Type:     pkg.Python,
			Name:     dep.Package,
			Revision: dep.Resolved,
		}

		// Update map.
		graph[id] = pkg.Package{
			ID:      id,
			Imports: packageImports(dep.Dependencies),
		}

		flattenDeepDependencies(graph, dep)
	}

	return graph
}

func flattenDeepDependencies(graph map[pkg.ID]pkg.Package, dep dependency) {
	for _, dep := range dep.Dependencies {
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

		// Update map.
		graph[id] = pkg.Package{
			ID:      id,
			Imports: packageImports(dep.Dependencies),
		}
		// Recurse in imports.
		flattenDeepDependencies(graph, dep)
	}
}

func packageImports(packageDeps []dependency) []pkg.Import {
	var imports []pkg.Import
	for _, i := range packageDeps {
		imports = append(imports, pkg.Import{
			Resolved: pkg.ID{
				Type:     pkg.Python,
				Name:     i.Package,
				Revision: i.Resolved,
			},
		})
	}
	return imports
}
