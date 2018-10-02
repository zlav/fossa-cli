package test_test

import (
	"encoding/json"
	"flag"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"

	"github.com/fossas/fossa-cli/api/fossa"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/test"
)

var orgID = rand.Intn(1000)

// taskStatus is a struct that imitates the anonymous struct within fossa.Build
type taskStatus struct {
	Status string
}

var testConfigs = []struct {
	fetcher string
	project string
	locator string
}{
	{"custom", "testRun", "custom+" + strconv.Itoa(orgID) + "%2FtestRun$1000"},
	{"custom", "git@github.com:fossas/fossa-cli.git", "custom+" + strconv.Itoa(orgID) + "%2Fgit@github.com:fossas%2Ffossa-cli.git$1000"},
	{"custom", "http://github.com/fossas/fossa-cli.git", "custom+" + strconv.Itoa(orgID) + "%2Fgithub.com%2Ffossas%2Ffossa-cli$1000"},
	{"git", "git@github.com:fossas/fossa-cli.git", "git+github.com%2Ffossas%2Ffossa-cli$1000"},
	{"git", "http://github.com/fossas/fossa-cli.git", "git+github.com%2Ffossas%2Ffossa-cli$1000"},
	{"git", "testRun", "git+testRun$1000"},
}

// This function tests how the cli's test command constructs a locator and sends a request to the desired endpoint.
// Currently this function implements goroutines and channels to prevent test.Run() from logging fatal and killing
// the process. In the future, once errors are handled more gracefully the test server can be broken out
// into its own TestServer package for other commands to access.
func TestTestRunLocators(t *testing.T) {
	for _, testConfig := range testConfigs {
		c := make(chan string)
		ts := testServer(testConfig.fetcher, testConfig.locator, c)
		if ts == nil {
			t.Fail()
		}
		defer ts.Close()

		flagSet := testFlags(testConfig.fetcher, testConfig.project, ts.URL, 1000)
		context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

		go func() {
			err := test.Run(context)
			assert.NoError(t, err)
		}()

		msg := <-c
		if msg != "SUCCESS" {
			assert.Equal(t, testConfig.locator, msg)
		}
	}
}

func testServer(fetcher, locator string, c chan string) *httptest.Server {
	if fetcher == "custom" {
		return testCustomServer(locator, c)
	}

	if fetcher == "git" {
		return testGitServer(locator, c)
	}

	return nil
}

func testCustomServer(locator string, c chan string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {
		case "/api/cli/organization":
			newResp := fossa.Organization{OrganizationID: orgID}
			request, _ := json.Marshal(newResp)
			_, err := w.Write(request)
			log.Debugf("error writing message: %s\n", err)
			return
		case "/api/cli/" + locator + "/latest_build":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			_, err := w.Write(request)
			log.Debugf("error writing message: %s\n", err)
			return
		case "/api/cli/" + locator + "/issues":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			_, err := w.Write(request)
			log.Debugf("error writing message: %s\n", err)
			c <- "SUCCESS"
			return
		default:
			c <- locatorFromPath(r.URL.EscapedPath())
		}
	}))
	return ts
}

func testGitServer(locator string, c chan string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {
		case "/api/cli/" + locator + "/latest_build":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			_, err := w.Write(request)
			log.Debugf("error writing message: %s\n", err)
			return
		case "/api/cli/" + locator + "/issues":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			_, err := w.Write(request)
			log.Debugf("error writing message: %s\n", err)
			c <- "SUCCESS"
			return
		default:
			c <- locatorFromPath(r.URL.EscapedPath())
		}
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

func locatorFromPath(path string) string {
	path = strings.TrimPrefix(path, "/api/cli/")
	path = strings.TrimSuffix(path, "/latest_build")
	path = strings.TrimSuffix(path, "/issues")
	return path
}
