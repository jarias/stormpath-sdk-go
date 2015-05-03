package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jarias/stormpath-sdk-go/web"

	"github.com/codegangsta/negroni"
	"github.com/jarias/stormpath-sdk-go"
	"github.com/julienschmidt/httprouter"

	"github.com/gorilla/sessions"
)

var indexHTML = `
<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title></title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    	Hello! <a href="/login">Login</a>
    </body>
</html>
`

var appHTML = `
<!doctype html>
<html class="no-js" lang="">
    <head>
        <meta charset="utf-8">
        <meta http-equiv="x-ua-compatible" content="ie=edge">
        <title></title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    	Cool your in! <a href="/logout">Logout</a>
    </body>
</html>
`

const sessionName = "go-sdk-demo"

var store = sessions.NewCookieStore([]byte("go-sdk-demo"))

func main() {
	credentials, _ := stormpath.NewDefaultCredentials()
	stormpath.Init(credentials, nil)

	n := negroni.Classic()

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if stormpathweb.IsAuthenticated(r) {
			http.Redirect(w, r, "/app", http.StatusFound)
			return
		}
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, indexHTML)
	})

	router.Handler("GET", "/login", loginHandler())
	router.Handler("GET", "/logout", logoutHandler())
	router.Handler("GET", "/callback", callbackHandler())

	authRouter := httprouter.New()

	authRouter.GET("/app", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprint(w, appHTML)
	})

	router.Handler("GET", "/app", negroni.New(
		negroni.HandlerFunc(authenticationMiddleware),
		negroni.Wrap(authRouter),
	))

	n.UseHandler(applicationMiddleware())
	n.UseHandler(accountMiddleware())
	n.UseHandler(router)

	n.Run(":9999")
}

func authenticationMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m := stormpathweb.AuthenticationMiddleware{
		Next:                next,
		SessionStore:        store,
		SessionName:         sessionName,
		UnauthorizedHandler: http.HandlerFunc(unauthorizedHandler),
	}
	m.ServeHTTP(rw, r)
}

func unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/?unauthorize", http.StatusFound)
}

func accountMiddleware() stormpathweb.AccountMiddleware {
	return stormpathweb.AccountMiddleware{
		SessionStore: store,
		SessionName:  sessionName,
	}
}

func applicationMiddleware() stormpathweb.ApplicationMiddleware {
	return stormpathweb.ApplicationMiddleware{
		ApplicationHref: os.Getenv("APPLICATION_HREF"),
	}
}

func loginHandler() stormpathweb.IDSiteLoginHandler {
	return stormpathweb.IDSiteLoginHandler{
		Options: map[string]string{"callbackURI": "/callback"},
	}
}

func logoutHandler() stormpathweb.IDSiteLogoutHandler {
	return stormpathweb.IDSiteLogoutHandler{
		Options: map[string]string{"callbackURI": "/callback"},
	}
}

func callbackHandler() stormpathweb.IDSiteAuthCallbackHandler {
	return stormpathweb.IDSiteAuthCallbackHandler{
		SessionStore:      store,
		SessionName:       sessionName,
		LoginRedirectURI:  "/app",
		LogoutRedirectURI: "/",
		ErrorHandler:      http.HandlerFunc(idSiteErrorHandler),
	}
}

func idSiteErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Ooops", http.StatusInternalServerError)
}
