package stormpathweb

import "github.com/jarias/stormpath-sdk-go"

type webContext struct {
	contentType   string
	webError      *errorModel
	postedData    map[string]string
	account       *stormpath.Account
	next          string
	originalError error
}

func newContext(contentType string, account *stormpath.Account) webContext {
	return webContext{contentType: contentType, account: account}
}

func (ctx webContext) withError(postedData map[string]string, err error) webContext {
	errorModel := buildErrorModel(err)
	return webContext{
		contentType:   ctx.contentType,
		account:       ctx.account,
		next:          ctx.next,
		postedData:    sanitizePostedData(postedData),
		webError:      &errorModel,
		originalError: err,
	}
}
