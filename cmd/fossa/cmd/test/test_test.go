package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fossas/fossa-cli/cmd/fossa/cmd/test"
	"github.com/fossas/fossa-cli/config"
)

func TestPublishWrongResponseStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)

		if r.Method != "POST" {
			t.Errorf("Expected ‘POST’ request, got ‘%s’", r.Method)
		}

		if r.URL.EscapedPath() != "/pub" {
			t.Errorf("Expected request to ‘/pub’, got ‘%s’", r.URL.EscapedPath())
		}
	}))

	defer ts.Close()

	config.BackendEndpoint = ts.URL
	test.Do()
}
