package stormpath_test

import (
	. "github.com/jarias/stormpath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	It("should retrive the current tenant", func() {
		cred, _ := NewDefaultCredentials()
		tenant, err := CurrentTenant(cred)

		Expect(err).NotTo(HaveOccurred())
		Expect(tenant.Href).NotTo(BeEmpty())
		Expect(tenant.Name).NotTo(BeEmpty())
		Expect(tenant.Key).NotTo(BeEmpty())
		Expect(tenant.Applications.Href).NotTo(BeEmpty())
		Expect(tenant.Directories.Href).NotTo(BeEmpty())
	})
})
