package stormpath_test

import (
	"regexp"
	. "github.com/jarias/stormpath"

	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the name", func() {
			application := NewApplication("name")

			jsonData, _ := json.Marshal(application)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})

	Describe("Save", func() {
		It("should create a new application", func() {
			application := NewApplication("create-test")

			err := application.Save()
			application.Purge()

			Expect(err).NotTo(HaveOccurred())
			Expect(application.Href).NotTo(BeEmpty())
			Expect(application.Status).To(Equal("ENABLED"))
		})
	})

	Describe("RegisterAccount", func() {
		It("should register a new account", func() {
			account := NewAccount("newaccount@test.org", "1234567z!A89", "test", "test")
			err := app.RegisterAccount(account)

			Expect(err).NotTo(HaveOccurred())
			Expect(account.Href).NotTo(BeEmpty())
		})
	})

	Describe("AuthenticateAccount", func() {
		It("should authenticate and return the account if the credentials are valid", func() {
			a, err := app.AuthenticateAccount("test@test.org", "1234567z!A89")
			Expect(err).NotTo(HaveOccurred())
			Expect(a.Account.Href).To(Equal(account.Href))
		})
	})

	Describe("groups", func() {
		Describe("CreateApplicationGroup", func() {
			It("should return error is group has no name", func() {
				err := app.CreateApplicationGroup(&Group{})

				Expect(err).To(HaveOccurred())
			})

			It("should create a new application group", func() {
				group := NewGroup("new-test-group")
				err := app.CreateApplicationGroup(group)

				Expect(err).NotTo(HaveOccurred())
				Expect(group.Href).NotTo(BeEmpty())
				Expect(group.Status).To(Equal("ENABLED"))
			})
		})

		Describe("GetApplicationGroups", func() {
			It("should return the paged list of application groups", func() {
				group := NewGroup("another-test-group")
				app.CreateApplicationGroup(group)

				groups, err := app.GetApplicationGroups(NewDefaultPageRequest(), DefaultFilter{})

				Expect(err).NotTo(HaveOccurred())

				Expect(groups.Href).NotTo(BeEmpty())
				Expect(groups.Offset).To(Equal(0))
				Expect(groups.Limit).To(Equal(25))
				Expect(groups.Items).NotTo(BeEmpty())
			})
		})
	})
	Describe("password reset", func() {
		Describe("SendPasswordResetEmail", func() {
			It("should create a new password reset token", func() {
				token, err := app.SendPasswordResetEmail(account.Email)

				Expect(err).NotTo(HaveOccurred())
				Expect(token.Href).NotTo(BeEmpty())
			})
		})

		Describe("ResetPassword", func() {
			It("should reset the account password", func() {
				token, _ := app.SendPasswordResetEmail(account.Email)

				re := regexp.MustCompile("[^\\/]+$")

				a, err := app.ResetPassword(re.FindString(token.Href), "8787987!kJKJdfW")

				Expect(err).NotTo(HaveOccurred())
				Expect(a.Account.Href).To(Equal(account.Href))
			})
		})

		Describe("ValidatePasswordResetToken", func() {
			It("should return the reset token if its valid", func() {
				token, _ := app.SendPasswordResetEmail(account.Email)

				re := regexp.MustCompile("[^\\/]+$")

				validatedToken, err := app.ValidatePasswordResetToken(re.FindString(token.Href))

				Expect(err).NotTo(HaveOccurred())
				Expect(validatedToken.Href).To(Equal(token.Href))
			})

			It("should return error if the token is invalid", func() {
				_, err := app.ValidatePasswordResetToken("invalid token")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
