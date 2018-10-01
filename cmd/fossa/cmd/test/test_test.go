package test

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fossas/fossa-cli/api/fossa"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

const (
	timeout = 3
	orgID   = 3
)

// taskStatus is a struct that imitates the anonymous struct within fossa.Build
type taskStatus struct {
	Status string
}

func TestRunCustomFetcherCustomProject(t *testing.T) {
	ts := testServer("custom+"+strconv.Itoa(orgID)+"%2FtestRun$1000", t)
	flagSet := testFlags("custom", "testRun", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func TestRunCustomFetcherGitSSHProject(t *testing.T) {
	ts := testServer("custom+"+strconv.Itoa(orgID)+"%2Fgit@github.com:fossa%2Ffossa-cli.git$1000", t)
	flagSet := testFlags("custom", "git@github.com:fossa/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func TestRunCustomFetcherGitHTTPProject(t *testing.T) {
	ts := testServer("custom+"+strconv.Itoa(orgID)+"%2Fgithub.com%2Ffossa%2Ffossa-cli$1000", t)
	flagSet := testFlags("custom", "http://github.com/fossa/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func TestRunGitFetcherGitSSHProject(t *testing.T) {
	ts := testServer("git+github.com%2Ffossas%2Ffossa-cli$1000", t)
	flagSet := testFlags("git", "git@github.com:fossas/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func TestRunGitFetcherGitHTTPProject(t *testing.T) {
	ts := testServer("git+github.com%2Ffossas%2Ffossa-cli$1000", t)
	flagSet := testFlags("git", "http://github.com/fossas/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func TestRunGitFetcherCustomProject(t *testing.T) {
	ts := testServer("git+testRun$1000", t)
	flagSet := testFlags("git", "testRun", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})
	Run(context)
	ts.Close()
}

func testServer(locator string, t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Path for OrganizationID
		if r.URL.EscapedPath() == "/api/cli/organization" {
			newResp := fossa.Organization{OrganizationID: orgID}
			request, _ := json.Marshal(newResp)
			fmt.Fprintf(w, string(request))
			return
		}

		// Path for check build
		if r.URL.EscapedPath() == "/api/cli/"+locator+"/latest_build" {
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			fmt.Fprintf(w, string(request))
			return
		}

		// Path for check issues
		if r.URL.EscapedPath() == "/api/cli/"+locator+"/issues" {
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			fmt.Fprintf(w, string(request))
			return
		}

		// Path for incorrect requests
		assert.Equal(t, "hsd", "sdfsadf")
		newResp := fossa.Build{Task: taskStatus{Status: "FAILED"}}
		request, _ := json.Marshal(newResp)
		fmt.Fprintf(w, string(request))
	}))

	return ts
}

func testFlags(fetcher, project, endpoint string, revision int) *flag.FlagSet {
	flagSet := &flag.FlagSet{}
	flagSet.Int("timeout", timeout, "")
	flagSet.Int("revision", revision, "")
	flagSet.String("fetcher", fetcher, "")
	flagSet.String("project", project, "")
	flagSet.String("endpoint", endpoint, "")
	return flagSet
}
