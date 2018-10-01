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

const orgID = 3

// taskStatus is a struct that imitates the anonymous struct within fossa.Build
type taskStatus struct {
	Status string
}

func TestRunCustomFetcherCustomProject(t *testing.T) {
	c := make(chan string)
	locator := "custom+" + strconv.Itoa(orgID) + "%2FtestRun$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("custom", "testRun", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func TestRunCustomFetcherGitSSHProject(t *testing.T) {
	c := make(chan string)
	locator := "custom+" + strconv.Itoa(orgID) + "%2Fgit@github.com:fossa%2Ffossa-cli.git$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("custom", "git@github.com:fossa/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func TestRunCustomFetcherGitHTTPProject(t *testing.T) {
	c := make(chan string)
	locator := "custom+" + strconv.Itoa(orgID) + "%2Fgithub.com%2Ffossa%2Ffossa-cli$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("custom", "http://github.com/fossa/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func TestRunGitFetcherGitSSHProject(t *testing.T) {
	c := make(chan string)
	locator := "git+github.com%2Ffossas%2Ffossa-cli$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("git", "git@github.com:fossas/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func TestRunGitFetcherGitHTTPProject(t *testing.T) {
	c := make(chan string)
	locator := "git+github.com%2Ffossas%2Ffossa-cli$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("git", "http://github.com/fossas/fossa-cli.git", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func TestRunGitFetcherCustomProject(t *testing.T) {
	c := make(chan string)
	locator := "git+testRun$1000"
	ts := testServer(locator, c)

	flagSet := testFlags("git", "testRun", ts.URL, 1000)
	context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

	go Run(context)
	msg := <-c
	if msg != "SUCCESS" {
		assert.Equal(t, locator, msg)
	}

	ts.Close()
}

func testServer(locator string, c chan string) *httptest.Server {
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
			c <- "SUCCESS"
			return
		}

		// Path for incorrect requests
		c <- r.URL.EscapedPath()
	}))

	return ts
}

func testFlags(fetcher, project, endpoint string, revision int) *flag.FlagSet {
	flagSet := &flag.FlagSet{}
	flagSet.Int("timeout", 5, "")
	flagSet.Int("revision", revision, "")
	flagSet.String("fetcher", fetcher, "")
	flagSet.String("project", project, "")
	flagSet.String("endpoint", endpoint, "")
	return flagSet
}
