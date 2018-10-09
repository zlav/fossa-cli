package fixtures

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/fossas/fossa-cli/files"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Project describes a git repo for a test fixture and env vars for when it is to be built
type Project struct {
	Name    string
	URL     string
	Commit  string
	Options map[string]interface{}
	Env     map[string]string
}

// Directory returns the directory under which tests fixtures should be placed
func Directory() string {
	return filepath.Join(os.TempDir(), "fossa-cli-fixtures")
}

var mutex = sync.Mutex{}

func ensureBaseDirectoryExists() error {
	_, err := createFixtureFolder(Directory())

	return err
}

// createFixtureFolder threadsafe creation of fixture folders. Returns true if they already exist
func createFixtureFolder(baseDir string) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()
	baseDirExists, err := files.ExistsFolder(baseDir)
	if err != nil {
		return false, err
	}
	if baseDirExists {
		println(baseDir + "already exists, assuming that clone has already been executed")
		return false, nil
	}

	err = os.MkdirAll(baseDir, os.FileMode(0700))
	if err != nil {
		return false, err
	}

	return false, nil
}

// ProjectInitializerFunction defines how a single project should be initialized *after* it has already been cloned
type ProjectInitializerFunction func(proj Project, projectDir string) error

// Initialize executes git clone in target directory and checksout the provided commit, then runts the initializerFn. This is done asynchronously for each provided project
func Initialize(baseDir string, projects []Project, initializerFn ProjectInitializerFunction) {
	ensureBaseDirectoryExists()
	fixturesExist, err := createFixtureFolder(baseDir)
	if err != nil {
		panic(err)
	}
	if fixturesExist {
		return
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(projects))

	for _, project := range projects {
		go func(proj Project) {
			defer waitGroup.Done()
			projectDir := filepath.Join(baseDir, proj.Name)

			repo, err := git.PlainClone(projectDir, false, &git.CloneOptions{
				URL:               proj.URL,
				RecurseSubmodules: 1,
			})
			if err != nil {
				panic(err)
			}

			worktree, err := repo.Worktree()
			if err != nil {
				panic(err)
			}
			err = worktree.Checkout(&git.CheckoutOptions{
				Hash: plumbing.NewHash(proj.Commit),
			})
			if err != nil {
				panic(err)
			}

			err = initializerFn(proj, projectDir)
			if err != nil {
				panic(err)
			}
		}(project)
	}

	waitGroup.Wait()
}
