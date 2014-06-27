package stormpath_test

import (
	. "github.com/jarias/stormpath"
	"github.com/jarias/stormpath/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	app     *Application
	cred    *Credentials
	account *Account
)

func TestStormpath(t *testing.T) {
	logger.InitInTestMode()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Stormpath Suite")
}

var _ = BeforeSuite(func() {
	var err error
	cred, err = NewDefaultCredentials()
	if err != nil {
		panic(err)
	}
	Client = NewStormpathClient(cred)

	app = NewApplication("test-app")

	err = app.Save()
	if err != nil {
		panic(err)
	}
	account = NewAccount("test@test.org", "1234567z!A89", "test", "test")
	app.RegisterAccount(account)
})

var _ = AfterSuite(func() {
	if app != nil {
		app.Purge()
	}
})
