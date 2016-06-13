package stormpathweb

import (
	"net/http"

	"github.com/jarias/stormpath-sdk-go"
	"golang.org/x/net/context"
)

type registerHandler struct {
	Parent      *StormpathMiddleware
	Application *stormpath.Application
	Form        form
}

func (h registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	if _, ok := isAuthenticated(w, r, ctx); ok {
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

func (h registerHandler) doGET(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	model := map[string]interface{}{
		"form": h.Form,
	}

	if contentType == stormpath.ApplicationJSON {
		model["accountStores"] = getApplicationAccountStores(h.Application)
		respondJSON(w, model, http.StatusOK)
		return
	}
	if contentType == stormpath.TextHTML {
		model["loginURI"] = Config.LoginURI
		model["postedData"] = ctx.Value("postedData")
		model["error"] = ctx.Value("error")

		respondHTML(w, model, Config.RegisterView)
	}
}

func (h registerHandler) doPOST(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	contentType := ctx.Value(ResolvedContentType)

	postedData, _ := getPostedData(r)

	err := validateForm(h.Form, postedData)
	if err != nil {
		h.handlePostError(w, r, ctx, err, postedData)
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
	if h.Parent.PreRegisterHandler != nil {
		pre := h.Parent.PreRegisterHandler(w, r, context.WithValue(ctx, AccountKey, account))
		if pre != nil {
			//User halted so we return
			return
		}
	}

	err = h.Application.RegisterAccount(account)
	if err != nil {
		h.handlePostError(w, r, ctx, err, postedData)
		return
	}

	if h.Parent.PostRegisterHandler != nil {
		post := h.Parent.PostRegisterHandler(w, r, context.WithValue(ctx, AccountKey, account))
		if post != nil {
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
				h.handlePostError(w, r, ctx, err, postedData)
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

func (h registerHandler) handlePostError(w http.ResponseWriter, r *http.Request, ctx context.Context, err error, postedData map[string]string) {
	contentType := ctx.Value(ResolvedContentType)

	if contentType == stormpath.TextHTML {
		//Sanitize postedData
		postedData["password"] = ""
		postedData["confirmPassword"] = ""

		h.doGET(w, r, contextWithError(ctx, err, postedData))
		return
	}
	if contentType == stormpath.ApplicationJSON {
		badRequest(w, r, err)
	}
}
