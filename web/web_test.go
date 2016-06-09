package stormpathweb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/jarias/stormpath-sdk-go/web"
)

func TestTCK(t *testing.T) {
	mux := http.NewServeMux()

	stormpathFilter := stormpathweb.NewStormpathHandler(mux, []string{"/"})

	mux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		fmt.Fprintf(w, "hello")
	}))

	ts := httptest.NewServer(stormpathFilter)
	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	cmd := exec.Command("./tck.sh", url.Host[strings.Index(url.Host, ":")+1:], stormpathFilter.Application.Href)

	err := cmd.Start()
	if err != nil {
		t.Errorf("Failed to start tck.sh script: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		t.Errorf("tck.sh fail: %s", err)
	}
}
