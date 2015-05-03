package stormpathweb

import (
	"github.com/gorilla/context"
	"github.com/jarias/stormpath-sdk-go"

	"net/http"
)

//ApplicationKey is the key of the current application in the context
const ApplicationKey = "application"

//AccountKey is the key of the current account in the context and session
const AccountKey = "account"

//GetApplication returns the application from the context previouly set by the ApplicationMiddleware
func GetApplication(r *http.Request) *stormpath.Application {
	app := context.Get(r, ApplicationKey)
	if app == nil {
		return nil
	}
	application := app.(stormpath.Application)
	return &application
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

//IsAuthenticated checks if there is an authenticated user
func IsAuthenticated(r *http.Request) bool {
	return GetCurrentAccount(r) != nil
}
