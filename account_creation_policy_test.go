package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountCreationPolicy", func() {
	Describe("GetVerificationEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{VerificationEmailTemplates: &EmailTemplates{}}
			policy.VerificationEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			_, err := policy.GetVerificationEmailTemplates()

			Expect(err).To(HaveOccurred())
		})

		It("should return the default verification email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()

			templates, err := policy.GetVerificationEmailTemplates()

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("GetVerificationSuccessEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{VerificationSuccessEmailTemplates: &EmailTemplates{}}
			policy.VerificationSuccessEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			_, err := policy.GetVerificationSuccessEmailTemplates()

			Expect(err).To(HaveOccurred())
		})

		It("should return the default verification success email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()

			templates, err := policy.GetVerificationSuccessEmailTemplates()

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("GetWelcomeEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{WelcomeEmailTemplates: &EmailTemplates{}}
			policy.WelcomeEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			_, err := policy.GetWelcomeEmailTemplates()

			Expect(err).To(HaveOccurred())
		})

		It("should return the default welcome email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()

			templates, err := policy.GetWelcomeEmailTemplates()

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("Update", func() {
		It("should update a given account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()
			policy.VerificationEmailStatus = Enabled
			err := policy.Update()

			Expect(err).NotTo(HaveOccurred())
			Expect(policy.VerificationEmailStatus).To(Equal(Enabled))
		})
	})
	Describe("Refresh", func() {
		It("should refresh a given account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()
			policy.VerificationEmailStatus = Enabled
			err := policy.Refresh()

			Expect(err).NotTo(HaveOccurred())
			Expect(policy.VerificationEmailStatus).To(Equal(Disabled))
		})
	})
})
