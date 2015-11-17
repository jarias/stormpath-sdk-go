package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OAuthPolicy", func() {
	Describe("Application.GetOAuthPolicy", func() {
		It("Should return the given policy", func() {
			application := newTestApplication()
			tenant.CreateApplication(application)

			policy, err := application.GetOAuthPolicy()

			Expect(err).NotTo(HaveOccurred())
			Expect(policy.AccessTokenTtl).To(Equal("PT1H"))
			Expect(policy.RefreshTokenTtl).To(Equal("P60D"))
		})
	})
	Describe("Update", func() {
		It("should update the policy if the new TTL values are valid", func() {
			application := newTestApplication()
			tenant.CreateApplication(application)

			policy, _ := application.GetOAuthPolicy()

			policy.AccessTokenTtl = "PT2H"
			policy.RefreshTokenTtl = "P50D"
			err := policy.Update()
			Expect(err).ToNot(HaveOccurred())

			policy, _ = application.GetOAuthPolicy()
			Expect("PT2H").To(Equal(policy.AccessTokenTtl))
			Expect("P50D").To(Equal(policy.RefreshTokenTtl))
		})
		It("should return an error if the access token TTL is invalid", func() {
			application := newTestApplication()
			tenant.CreateApplication(application)

			policy, _ := application.GetOAuthPolicy()

			policy.AccessTokenTtl = "hello i'm not a valid value deal with me"
			err := policy.Update()
			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(400))
			Expect(err.(Error).Code).To(Equal(2002))
		})
		It("should return an error if the refresh token TTL is invalid", func() {
			application := newTestApplication()
			tenant.CreateApplication(application)

			policy, _ := application.GetOAuthPolicy()

			policy.RefreshTokenTtl = "hello i'm not a valid value deal with me"
			err := policy.Update()
			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(400))
			Expect(err.(Error).Code).To(Equal(2002))
		})
	})
})
