package test_test

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"

	"github.com/fossas/fossa-cli/api/fossa"
	"github.com/fossas/fossa-cli/cmd/fossa/cmd/test"
)

const orgID = 3

// taskStatus is a struct that imitates the anonymous struct within fossa.Build
type taskStatus struct {
	Status string
}

var locatorTypes = []struct {
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

func TestTestRunLocators(t *testing.T) {
	for _, locatorType := range locatorTypes {
		c := make(chan string)
		ts := testServer(locatorType.locator, c)
		defer ts.Close()

		flagSet := testFlags(locatorType.fetcher, locatorType.project, ts.URL, 1000)
		context := cli.NewContext(&cli.App{}, flagSet, &cli.Context{})

		go test.Run(context)
		msg := <-c
		if msg != "SUCCESS" {
			assert.Equal(t, locatorType.locator, msg)
		}
	}
}

func testServer(locator string, c chan string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.EscapedPath() {
		case "/api/cli/organization":
			newResp := fossa.Organization{OrganizationID: orgID}
			request, _ := json.Marshal(newResp)
			w.Write(request)
			return
		case "/api/cli/" + locator + "/latest_build":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			w.Write(request)
			return
		case "/api/cli/" + locator + "/issues":
			newResp := fossa.Build{Task: taskStatus{Status: "SUCCEEDED"}}
			request, _ := json.Marshal(newResp)
			w.Write(request)
			c <- "SUCCESS"
			return
		default:
			c <- r.URL.EscapedPath()
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
