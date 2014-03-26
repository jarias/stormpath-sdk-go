package stormpath_test

import (
	. "github.com/jarias/stormpath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Credentials", func() {
	It("should load the API credentials from a valid file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/apiKeys.properties")
		Expect(err).NotTo(HaveOccurred())
		Expect(credentials.Id).To(Equal("APIKEY"))
		Expect(credentials.Secret).To(Equal("APISECRET"))
	})

	It("should return an error loading credentials from a unexisting file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/doesntexist.properties")
		Expect(err).To(HaveOccurred())
		Expect(credentials).To(BeNil())
	})

	It("should load empty credentials from an empty properties file", func() {
		credentials, err := NewCredentialsFromFile("./test_files/empty.properties")
		Expect(err).NotTo(HaveOccurred())
		Expect(credentials.Id).To(BeEmpty())
		Expect(credentials.Secret).To(BeEmpty())
	})

	/*
		Can't assume that ~/.config/stormpath/apiKeys.properties exists in the dev env

		It("should load the API credentials from a default file location", func() {
			credentials, err := NewDefaultCredentials()
			Expect(err).NotTo(HaveOccurred())
			Expect(credentials.Id).NotTo(BeEmpty())
			Expect(credentials.Secret).NotTo(BeEmpty())
		})
	*/
})
