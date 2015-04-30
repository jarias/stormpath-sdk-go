package stormpath_test

import (
	"net/url"
	"regexp"

	"encoding/json"

	"github.com/dgrijalva/jwt-go"
	. "github.com/jarias/stormpath-sdk-go"
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
		It("should update an existing application", func() {
			app.Name = "new-name" + randomName()
			err := app.Save()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RegisterAccount", func() {
		It("should register a new account", func() {
			account := newTestAccount()
			err := app.RegisterAccount(account)

			Expect(err).NotTo(HaveOccurred())
			Expect(account.Href).NotTo(BeEmpty())
		})
	})

	Describe("AuthenticateAccount", func() {
		It("should authenticate and return the account if the credentials are valid", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			a, err := app.AuthenticateAccount(account.Email, "1234567z!A89")
			Expect(err).NotTo(HaveOccurred())
			Expect(a.Account.Href).To(Equal(account.Href))
		})
	})

	Describe("groups", func() {
		Describe("CreateGroup", func() {
			It("should return error is group has no name", func() {
				err := app.CreateGroup(&Group{})

				Expect(err).To(HaveOccurred())
			})

			It("should create a new application group", func() {
				group := NewGroup("new-test-group")
				err := app.CreateGroup(group)

				Expect(err).NotTo(HaveOccurred())
				Expect(group.Href).NotTo(BeEmpty())
				Expect(group.Status).To(Equal("ENABLED"))
			})
		})

		Describe("GetApplicationGroups", func() {
			It("should return the paged list of application groups", func() {
				group := newTestGroup()
				app.CreateGroup(group)

				groups, err := app.GetGroups(NewDefaultPageRequest(), NewEmptyFilter())

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

		Describe("CreateIDSiteURL", func() {
			It("Should create valid ID Site URL", func() {
				idSiteURL, err := app.CreateIDSiteURL(map[string]string{"callbackURI": "http://localhost:8080"})

				u, _ := url.Parse(idSiteURL)

				Expect(err).NotTo(HaveOccurred())
				Expect(u.Path).To(Equal("/sso"))
				Expect(u.Query()).NotTo(BeEmpty())

				//Check Token
				jwtRequest := u.Query().Get("jwtRequest")

				token, _ := jwt.Parse(jwtRequest, func(token *jwt.Token) (interface{}, error) {
					return []byte(cred.Secret), nil
				})

				Expect(token.Valid).To(BeTrue())

				Expect(token.Claims["cb_uri"]).To(Equal("http://localhost:8080"))
				Expect(token.Claims["state"]).To(Equal(""))
				Expect(token.Claims["path"]).To(Equal("/"))
				Expect(token.Claims["iss"]).To(Equal(cred.ID))
				Expect(token.Claims["sub"]).To(Equal(app.Href))
				Expect(token.Claims["jti"]).NotTo(BeEmpty())
				Expect(token.Claims["iat"]).To(BeNumerically(">", 0))
			})

			It("Should create valid ID Site logout URL", func() {
				idSiteURL, err := app.CreateIDSiteURL(
					map[string]string{
						"callbackURI": "http://localhost:8080",
						"logout":      "true",
					})

				u, _ := url.Parse(idSiteURL)

				Expect(err).NotTo(HaveOccurred())
				Expect(u.Path).To(Equal("/sso/logout"))
				Expect(u.Query()).NotTo(BeEmpty())

				//Check Token
				jwtRequest := u.Query().Get("jwtRequest")

				token, _ := jwt.Parse(jwtRequest, func(token *jwt.Token) (interface{}, error) {
					return []byte(cred.Secret), nil
				})

				Expect(token.Valid).To(BeTrue())

				Expect(token.Claims["cb_uri"]).To(Equal("http://localhost:8080"))
				Expect(token.Claims["state"]).To(Equal(""))
				Expect(token.Claims["path"]).To(Equal("/"))
				Expect(token.Claims["iss"]).To(Equal(cred.ID))
				Expect(token.Claims["sub"]).To(Equal(app.Href))
				Expect(token.Claims["jti"]).NotTo(BeEmpty())
				Expect(token.Claims["iat"]).To(BeNumerically(">", 0))
			})
		})
	})
})
