package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
)

type registerHandler struct {
	PreRegisterHandler  UserHandler
	PostRegisterHandler UserHandler
	Application         *stormpath.Application
}

func (h registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx webContext) {
	if ctx.Account != nil {
		http.Redirect(w, r, Config.LoginNextURI, http.StatusFound)
		return
	}

	if Config.IDSiteEnabled {
		options := stormpath.IDSiteOptions{
			Path:        Config.IDSiteRegisterURI,
			CallbackURL: baseURL(r) + Config.CallbackURI,
		}

		idSiteURL, err := h.Application.CreateIDSiteURL(options)
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
	contentType := ctx.ContentType

	model := map[string]interface{}{
		"form": Config.RegisterForm,
	}

	if contentType == stormpath.ApplicationJSON {
		model["accountStores"] = getApplicationAccountStores(h.Application)
		respondJSON(w, model, http.StatusOK)
		return
	}
	if contentType == stormpath.TextHTML {
		model["loginURI"] = Config.LoginURI
		model["postedData"] = ctx.PostedData
		model["error"] = ctx.Error

		respondHTML(w, model, Config.RegisterView)
	}
}

func (h registerHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx webContext) {
	contentType := ctx.ContentType

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
	if h.PreRegisterHandler != nil {
		pre := h.PreRegisterHandler(w, r, account)
		if !pre {
			//User halted so we return
			return
		}
	}

	err = h.Application.RegisterAccount(account)
	if err != nil {
		handleError(w, r, ctx.withError(postedData, err), h.doGET)
		return
	}

	if h.PostRegisterHandler != nil {
		post := h.PostRegisterHandler(w, r, account)
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
			err := saveAuthenticationResult(w, r, transientAuthenticationResult(account), h.Application)
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
