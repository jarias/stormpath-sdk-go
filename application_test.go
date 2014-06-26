package stormpath_test

import (
	"regexp"
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

		logger.InitInTestMode()
	})

	AfterEach(func() {
		if app != nil {
			app.Purge()
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
			err := app.RegisterAccount(account)

			Expect(err).NotTo(HaveOccurred())
			Expect(account.Href).NotTo(BeEmpty())
		})
	})

	Describe("AuthenticateAccount", func() {
		It("should authenticate and return the account if the credentials are valid", func() {
			app = NewApplication("authorize-account", client)
			app.Save()
			account := NewAccount("test@test.org", "1234567z!A89", "test", "test")
			app.RegisterAccount(account)

			a, err := app.AuthenticateAccount("test@test.org", "1234567z!A89")
			Expect(err).NotTo(HaveOccurred())
			Expect(a.Account.Href).To(Equal(account.Href))
		})
	})

	Describe("password reset", func() {
		Describe("SendPasswordResetEmail", func() {
			It("should create a new password reset token", func() {
				app = NewApplication("password-reset", client)
				app.Save()
				account := NewAccount("test@test.org", "1234567z!A89", "test", "test")
				app.RegisterAccount(account)
				token, err := app.SendPasswordResetEmail(account.Email)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.Href).NotTo(BeEmpty())
			})
		})

		Describe("ResetPassword", func() {
			It("should reset the account password", func() {
				app = NewApplication("password-reset", client)
				app.Save()
				account := NewAccount("test@test.org", "1234567z!A89", "test", "test")
				app.RegisterAccount(account)
				token, _ := app.SendPasswordResetEmail(account.Email)

				re := regexp.MustCompile("[^\\/]+$")

				a, err := app.ResetPassword(re.FindString(token.Href), "8787987!kJKJdfW")

				Expect(err).NotTo(HaveOccurred())
				Expect(a.Account.Href).To(Equal(account.Href))
			})
		})

		Describe("ValidatePasswordResetToken", func() {
			It("should return the reset token if its valid", func() {
				app = NewApplication("password-reset", client)
				app.Save()
				account := NewAccount("test@test.org", "1234567z!A89", "test", "test")
				app.RegisterAccount(account)
				token, _ := app.SendPasswordResetEmail(account.Email)

				re := regexp.MustCompile("[^\\/]+$")

				validatedToken, err := app.ValidatePasswordResetToken(re.FindString(token.Href))

				Expect(err).NotTo(HaveOccurred())
				Expect(validatedToken.Href).To(Equal(token.Href))
			})

			It("should return error if the token is invalid", func() {
				app = NewApplication("password-reset", client)
				app.Save()
				_, err := app.ValidatePasswordResetToken("invalid token")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
