package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type registerHandler struct {
	preRegisterHandler  UserHandler
	postRegisterHandler UserHandler
	application         *stormpath.Application
}

func (h registerHandler) serveHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if ctx.account != nil {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteRegisterURI,
			CallbackURL: baseURL(r) + Config.CallbackURI,
		}

		idSiteURL, err := h.application.CreateIDSiteURL(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, idSiteURL, http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		h.doPOST(w, r, ctx)
		return
	}
	if r.Method == http.MethodGet {
		h.doGET(w, r, ctx)
		return
	}

	methodNotAllowed(w, r, ctx)
}

func (h registerHandler) doGET(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.contentType

	model := map[string]interface{}{
		"form": Config.RegisterForm,
	}

	if contentType == stormpath.ApplicationJSON {
		model["accountStores"] = getApplicationAccountStores(h.application)
		respondJSON(w, model, http.StatusOK)
		return
	}
	if contentType == stormpath.TextHTML {
		model["loginURI"] = Config.LoginURI
		model["postedData"] = ctx.postedData
		model["error"] = ctx.webError

		respondHTML(w, model, Config.RegisterView)
	}
}

func (h registerHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.contentType

	postedData, _ := getPostedData(r)

	err := validateForm(Config.RegisterForm, postedData)
	if err != nil {
		handleError(w, r, ctx.withError(postedData, err), h.doGET)
		return
	}

	account := &stormpath.Account{}

	account.Email = postedData["email"]
	account.Username = postedData["username"]
	account.Password = postedData["password"]
	account.GivenName = postedData["givenName"]
	if account.GivenName == "" {
		account.GivenName = "UNKNOWN"
	}
	account.Surname = postedData["surname"]
	if account.Surname == "" {
		account.Surname = "UNKNOWN"
	}

	//TODO custom data
	if h.preRegisterHandler != nil {
		pre := h.preRegisterHandler(w, r, account)
		if !pre {
			//User halted so we return
			return
		}
	}

	err = h.application.RegisterAccount(account)
	if err != nil {
		handleError(w, r, ctx.withError(postedData, err), h.doGET)
		return
	}

	if h.postRegisterHandler != nil {
		post := h.postRegisterHandler(w, r, account)
		if !post {
			//user halted so we return
			return
		}
	}

	if contentType == stormpath.ApplicationJSON {
		respondJSON(w, accountModel(account), http.StatusOK)
		return
	}

	accountStatus := account.Status

	if accountStatus == stormpath.Enabled {
		if Config.RegisterAutoLoginEnabled {
			err := saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.application)
			if err != nil {
				handleError(w, r, ctx.withError(postedData, err), h.doGET)
				return
			}
		}
		http.Redirect(w, r, Config.LoginURI+"?status=created", http.StatusFound)
		return
	} else if accountStatus == stormpath.Unverified {
		http.Redirect(w, r, Config.LoginURI+"?status=unverified", http.StatusFound)
		return
	}
	http.Redirect(w, r, Config.RegisterNextURI, http.StatusFound)
}
