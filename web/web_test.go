package stormpathweb

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"strings"
	"testing"

	"github.com/jarias/stormpath-sdk-go"
)

var mainTemplate = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <title>Example</title>

    <!-- Bootstrap -->
    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7"
        crossorigin="anonymous">
    <link rel="stylesheet" href="/stormpath/assets/css/stormpath.css">

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
</head>

<body>
    <div class="container">		
		{{ if .account }}
		<h1>Hello {{ .account.FullName }}</h1>
		<h4>Provider: {{ .account.ProviderData.ProviderID }}</h4>     
		<form id="logoutForm" action="{{ .logoutUri }}" method="post">
        	<input type="submit" class="btn btn-danger" value="Logout"/>
        </form>
		{{ else }}
		<h1>Hello World</h1>
		<a href="{{ .loginUri }}" class="btn btn-primary">Login</a>
		{{ end }}
    </div>
</body>

</html>
`

func GetTestServer() (*httptest.Server, string) {
	mux := http.NewServeMux()

	stormpathMiddleware := NewStormpathMiddleware(mux)

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account := stormpathMiddleware.GetAuthenticatedAccount(w, r)

		w.Header().Add("Content-Type", "text/html")

		template, err := template.New("main").Parse(mainTemplate)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		model := map[string]interface{}{
			"account":   account,
			"loginUri":  Config.LoginURI,
			"logoutUri": Config.LogoutURI,
		}

		if account != nil {
			model["name"] = account.GivenName
		}

		template.Execute(w, model)
	}))

	return httptest.NewServer(stormpathMiddleware), stormpathMiddleware.Application.Href
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
