package stormpathweb

import "github.com/jarias/stormpath-sdk-go"

type webContext struct {
	ContentType   string
	Error         *errorModel
	PostedData    map[string]string
	Account       *stormpath.Account
	Next          string
	originalError error
}

func newContext(contentType string, account *stormpath.Account) webContext {
	return webContext{ContentType: contentType, Account: account}
}

func (ctx webContext) withError(postedData map[string]string, err error) webContext {
	errorModel := buildErrorModel(err)
	return webContext{
		ContentType:   ctx.ContentType,
		Account:       ctx.Account,
		Next:          ctx.Next,
		PostedData:    sanitizePostedData(postedData),
		Error:         &errorModel,
		originalError: err,
	}
}

func (ctx webContext) getError() error {
	return ctx.originalError
}
