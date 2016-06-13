package stormpathweb

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/jarias/stormpath-sdk-go"
)

func GetTestServer() (*httptest.Server, string) {
	mux := http.NewServeMux()

	stormpathFilter := NewStormpathHandler(mux, []string{"/"})

	mux.Handle("/hello", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		fmt.Fprintf(w, "hello")
	}))

	return httptest.NewServer(stormpathFilter), stormpathFilter.Application.Href
}

func BenchmarkGETLoginHTML(b *testing.B) {
	ts, _ := GetTestServer()
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/login", nil)
		req.Header.Set(stormpath.AcceptHeader, stormpath.TextHTML)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGETLoginJSON(b *testing.B) {
	ts, _ := GetTestServer()
	defer ts.Close()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/login", nil)
		req.Header.Set(stormpath.AcceptHeader, stormpath.ApplicationJSON)
		_, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
	}
}

func TestTCK(t *testing.T) {
	ts, applicationHref := GetTestServer()
	defer ts.Close()

	url, _ := url.Parse(ts.URL)

	cmd := exec.Command("./tck.sh", url.Host[strings.Index(url.Host, ":")+1:], applicationHref)

	err := cmd.Start()
	if err != nil {
		t.Errorf("Failed to start tck.sh script: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		t.Errorf("tck.sh fail: %s", err)
	}
}
