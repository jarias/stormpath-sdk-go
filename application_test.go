package stormpath_test

import (
	. "github.com/jarias/stormpath"
	"github.com/jarias/stormpath/logger"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	var (
		cred   *Credentials
		client *StormpathClient
		app    *Application
	)

	BeforeEach(func() {
		var err error
		cred, err = NewDefaultCredentials()
		if err != nil {
			panic(err)
		}
		client = NewStormpathClient(cred)

		logger.Init(false)
	})

	AfterEach(func() {
		if app != nil {
			app.Delete()
		}
	})

	Describe("JSON", func() {
		It("should marshall a minimum JSON with only the name", func() {
			application := NewApplication("name", client)

			jsonData, _ := json.Marshal(application)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})

	Describe("Save", func() {
		It("should create a new application", func() {
			app = NewApplication("create-test", client)

			err := app.Save()

			Expect(err).NotTo(HaveOccurred())
			Expect(app.Href).NotTo(BeEmpty())
			Expect(app.Status).To(Equal("ENABLED"))
		})
	})

	Describe("RegisterAccount", func() {
		It("should register a new account", func() {
			app = NewApplication("register-account-test", client)

			app.Save()
			account := NewAccount("test@test.org", "1234567z!A89", "test", "test")
			app.RegisterAccount(account)

			Expect(account.Href).NotTo(BeEmpty())
		})
	})
})
