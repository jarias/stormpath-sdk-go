package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"

	stormpath "github.com/jarias/stormpath-sdk-go"
	"github.com/jarias/stormpath-sdk-go/web"
)

var helloTemplate = `
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

func main() {
	mux := http.NewServeMux()

	stormpathMiddleware := stormpathweb.NewStormpathMiddleware(mux, nil)

	stormpathMiddleware.SetPreLoginHandler(func(w http.ResponseWriter, r *http.Request, account *stormpath.Account) bool {
		fmt.Println("--> Pre Login")
		return true
	})

	stormpathMiddleware.SetPostLoginHandler(func(w http.ResponseWriter, r *http.Request, account *stormpath.Account) bool {
		fmt.Println("--> Post Login")
		return true
	})

	stormpathMiddleware.SetPreRegisterHandler(func(w http.ResponseWriter, r *http.Request, account *stormpath.Account) bool {
		fmt.Println("--> Pre Register")
		return true
	})

	stormpathMiddleware.SetPostRegisterHandler(func(w http.ResponseWriter, r *http.Request, account *stormpath.Account) bool {
		fmt.Println("--> Post Register")
		return true
	})

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		account := stormpathMiddleware.GetAuthenticatedAccount(w, r)

		w.Header().Add("Content-Type", "text/html")

		template, err := template.New("hello").Parse(helloTemplate)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		model := map[string]interface{}{
			"account":   account,
			"loginUri":  stormpathweb.Config.LoginURI,
			"logoutUri": stormpathweb.Config.LogoutURI,
		}

		if account != nil {
			model["name"] = account.GivenName
		}

		template.Execute(w, model)
	}))

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	fmt.Println("Starting example in port 8080 CTRL+C to stop")
	log.Fatal(http.ListenAndServe(":8080", stormpathMiddleware))
}
