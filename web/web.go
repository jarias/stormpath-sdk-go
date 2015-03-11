package stormpathweb

import (
	"encoding/json"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jarias/stormpath-sdk-go"

	"net/http"
)

//ApplicationKey is the key of the current application in the context
const ApplicationKey = "application"

//AccountKey is the key of the current account in the context and session
const AccountKey = "account"

//LoginLogoutHandler is an http.Handler use Stormpath ID Site login and logout
type LoginLogoutHandler struct {
	Options map[string]string
}

//AuthCallbackHandler is an http.Handler for the ID Site callback
type AuthCallbackHandler struct {
	SessionStore sessions.Store
	SessionName  string
	RedirectURI  string
}

//ApplicationMiddleware is an http.Handler that stores a given account in the request context
//to be use by any other handler in the chain
type ApplicationMiddleware struct {
	ApplicationHref string
}

//AccountMiddleware is an http.Handler that unmarshals the current account store in the session
//and stores it in the request context to be use by any other handler in the chain
type AccountMiddleware struct {
	SessionStore sessions.Store
	SessionName  string
}

//LoginHandler returns a http.Handler for ID Site login
func LoginHandler(options map[string]string) LoginLogoutHandler {
	//Ensure logout flag is not present
	if options["logout"] != "" {
		options["logout"] = ""
	}

	return LoginLogoutHandler{options}
}

//LogoutHandler returns a http.Handler for ID Site logout
func LogoutHandler(options map[string]string) LoginLogoutHandler {
	//Ensure logout flag is present
	if options["logout"] == "" {
		options["logout"] = "true"
	}

	return LoginLogoutHandler{options}
}

func (h LoginLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idSiteURL, _ := GetApplication(r).CreateIDSiteURL(h.Options)

	http.Redirect(w, r, idSiteURL, http.StatusFound)
}

//GetApplication returns the application for a web app from the context previouly set by the ApplicationMiddleware
func GetApplication(r *http.Request) *stormpath.Application {
	app := context.Get(r, ApplicationKey).(stormpath.Application)
	return &app
}

func (am ApplicationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Check if it the current app already exists
	app := context.Get(r, ApplicationKey)
	if app == nil {
		app, err := stormpath.NewApplicationRef(am.ApplicationHref).GetApplication()
		if err == nil {
			context.Set(r, ApplicationKey, *app)
		}
	}
}

func (accm AccountMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, _ := accm.SessionStore.Get(r, accm.SessionName)

	if session.Values[AccountKey] != nil {
		account := stormpath.Account{}

		json.Unmarshal([]byte(session.Values[AccountKey].([]uint8)), &account)
		context.Set(r, AccountKey, account)
	}
}

//GetCurrentAccount retrives the current account if any from the request context
func GetCurrentAccount(r *http.Request) *stormpath.Account {
	acc := context.Get(r, AccountKey)
	if acc == nil {
		return nil
	}
	account := acc.(stormpath.Account)
	return &account
}

//ServeHTTP AuthCallbackHandler http.Handler implementation
func (cbh AuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app := GetApplication(r)
	result, err := app.HandleIDSiteCallback(r.URL.String())

	if err != nil {
		//TODO log error or something
		http.Redirect(w, r, cbh.RedirectURI, 400)
		return
	}

	if result.Status == "AUTHENTICATED" {
		cbh.storeAccountInSession(result.Account, w, r)
	} else {
		cbh.clearAccountInSession(w, r)
	}
	http.Redirect(w, r, cbh.RedirectURI, http.StatusFound)
}

//StoreAccountInSession stores a given account in the session as the current account
func (cbh AuthCallbackHandler) storeAccountInSession(account *stormpath.Account, w http.ResponseWriter, r *http.Request) {
	session, _ := cbh.SessionStore.Get(r, cbh.SessionName)

	jsonBody, _ := json.Marshal(account)

	session.Values[AccountKey] = jsonBody
	session.Save(r, w)
}

//ClearAccountInSession removes the current account form the session
func (cbh AuthCallbackHandler) clearAccountInSession(w http.ResponseWriter, r *http.Request) {
	session, _ := cbh.SessionStore.Get(r, cbh.SessionName)

	session.Values[AccountKey] = nil
	session.Save(r, w)
}

//AuthenticationMiddleware handles authentication for a web application, it should only be apply to http.Handlers
//that require authentication it checks the session for current account if exists it calls handler else redirects with 401
//to the given UnauthorizedRedirectURL
type AuthenticationMiddleware struct {
	Next                    http.Handler
	SessionStore            sessions.Store
	SessionName             string
	UnauthorizedRedirectURL string
}

func (am AuthenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, _ := am.SessionStore.Get(r, am.SessionName)

	if session.Values[AccountKey] == nil {
		if am.UnauthorizedRedirectURL == "" {
			//Send 401
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		http.Redirect(w, r, am.UnauthorizedRedirectURL, http.StatusFound)
		return
	}

	am.Next.ServeHTTP(w, r)
}

//Authenticate is a convinices method to authenticate a single http.Handler
func (am AuthenticationMiddleware) Authenticate(next http.Handler) http.Handler {
	am.Next = next
	return am
}
