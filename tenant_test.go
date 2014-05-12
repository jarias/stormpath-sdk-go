package stormpath_test

import (
	. "github.com/jarias/stormpath"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	var cred *Credentials

	BeforeEach(func() {
		var err error
		cred, err = NewDefaultCredentials()
		if err != nil {
			panic(err)
		}
	})

	Describe("CurrentTentant", func() {
		It("should retrive the current tenant", func() {
			tenant, err := CurrentTenant(cred)

			Expect(err).NotTo(HaveOccurred())
			Expect(tenant.Href).NotTo(BeEmpty())
			Expect(tenant.Name).NotTo(BeEmpty())
			Expect(tenant.Key).NotTo(BeEmpty())
			Expect(tenant.Applications.Href).NotTo(BeEmpty())
			Expect(tenant.Directories.Href).NotTo(BeEmpty())
		})
	})

	Describe("Tenant.GetApplications", func() {
		It("should retrive all the tenant applications", func() {
			tenant, _ := CurrentTenant(cred)

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), ApplicationFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(25))
			Expect(apps.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant applications by page", func() {
			tenant, _ := CurrentTenant(cred)

			apps, err := tenant.GetApplications(NewPageRequest(1, 0), ApplicationFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(1))
			Expect(apps.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant applications by page and filter", func() {
			tenant, _ := CurrentTenant(cred)

			f := ApplicationFilter{Name: "stormpath"}

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Items).To(HaveLen(1))
		})
	})
})
