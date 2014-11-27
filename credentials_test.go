package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Credentials", func() {
	It("should load the API credentials from a valid file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/apiKeys.properties")
		Expect(err).NotTo(HaveOccurred())
		Expect(credentials.ID).To(Equal("APIKEY"))
		Expect(credentials.Secret).To(Equal("APISECRET"))
	})

	It("should return an error loading credentials from a unexisting file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/doesntexist.properties")
		Expect(err).To(HaveOccurred())
		Expect(credentials).To(Equal(Credentials{}))
	})

	It("should load empty credentials from an empty properties file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/empty.properties")
		Expect(err).NotTo(HaveOccurred())
		Expect(credentials.ID).To(BeEmpty())
		Expect(credentials.Secret).To(BeEmpty())
	})

	It("should load the API credentials from a default file location", func() {
		credentials, err := NewDefaultCredentials()
		Expect(err).NotTo(HaveOccurred())
		Expect(credentials.ID).NotTo(BeEmpty())
		Expect(credentials.Secret).NotTo(BeEmpty())
	})
})
