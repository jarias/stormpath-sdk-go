package stormpathweb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"github.com/jarias/stormpath-sdk-go"
)

//IDSiteLoginHandler is an http.Handler for Strompath's IDSite login
type IDSiteLoginHandler struct {
	Options map[string]string
}

//IDSiteLogoutHandler is an http.Handler for Strompath's IDSite logout
type IDSiteLogoutHandler struct {
	Options map[string]string
}

//ServeHTTP implements the http.Handler interface for IDSiteLoginHandler type
func (h IDSiteLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idSiteURLHandler(w, r, ensureOption("logout", "", h.Options))
}

//ServeHTTP implements the http.Handler interface for IDSiteLogoutHandler type
func (h IDSiteLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idSiteURLHandler(w, r, ensureOption("logout", "true", h.Options))
}

func idSiteURLHandler(w http.ResponseWriter, r *http.Request, options map[string]string) {
	if options["callbackURI"][0] == '/' {
		u, _ := url.Parse(r.Header.Get("Referer"))
		callbackURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, options["callbackURI"])
		options["callbackURI"] = callbackURL
	}

	idSiteURL, _ := GetApplication(r).CreateIDSiteURL(options)

	http.Redirect(w, r, idSiteURL, http.StatusFound)
}

func ensureOption(key string, value string, options map[string]string) map[string]string {
	if options == nil {
		options = make(map[string]string)
	}
	options[key] = value
	return options
}

//IDSiteAuthCallbackHandler is an http.Handler for the ID Site callback
type IDSiteAuthCallbackHandler struct {
	SessionStore      sessions.Store
	SessionName       string
	LoginRedirectURI  string
	LogoutRedirectURI string
	ErrorHandler      http.Handler
}

//ServeHTTP implements the http.Handler interface for the IDSiteAuthCallbackHandler type
func (h IDSiteAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app := GetApplication(r)
	result, err := app.HandleIDSiteCallback(r.URL.String())

	if err != nil {
		stormpath.Logger.Printf("[ERROR] IDSite %s", err)
		h.ErrorHandler.ServeHTTP(w, r)
		return
	}

	if result.Status == "AUTHENTICATED" {
		//Login succesfull
		h.storeAccountInSession(result.Account, w, r)
		http.Redirect(w, r, h.LoginRedirectURI, http.StatusFound)
	} else {
		//Logout
		h.clearAccountInSession(w, r)
		http.Redirect(w, r, h.LogoutRedirectURI, http.StatusFound)
	}
}

//StoreAccountInSession stores a given account in the session as the current account
func (h IDSiteAuthCallbackHandler) storeAccountInSession(account *stormpath.Account, w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, h.SessionName)

	jsonBody, _ := json.Marshal(account)

	session.Values[AccountKey] = jsonBody
	session.Save(r, w)
}

//ClearAccountInSession removes the current account form the session
func (h IDSiteAuthCallbackHandler) clearAccountInSession(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionStore.Get(r, h.SessionName)

	session.Values[AccountKey] = nil
	session.Save(r, w)
}
