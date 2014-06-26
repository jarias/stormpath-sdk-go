package stormpath_test

import (
	. "github.com/jarias/stormpath"
	"github.com/jarias/stormpath/logger"
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

		logger.InitInTestMode()
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

	Describe("Tenant.GetDirectories", func() {
		It("should retrive all the tenant directories", func() {
			tenant, _ := CurrentTenant(cred)

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(25))
			Expect(directories.Items).NotTo(BeEmpty())
			for _, d := range directories.Items {
				Expect(d.Client).NotTo(BeNil())
			}
		})

		It("should retrive all the tenant directories by page", func() {
			tenant, _ := CurrentTenant(cred)

			directories, err := tenant.GetDirectories(NewPageRequest(1, 0), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(1))
			Expect(directories.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant directories by page and filter", func() {
			tenant, _ := CurrentTenant(cred)

			f := DefaultFilter{Name: "Stormpath Administrators"}

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Items).To(HaveLen(1))
		})

	})

	Describe("Tenant.GetApplications", func() {
		It("should retrive all the tenant applications", func() {
			tenant, _ := CurrentTenant(cred)

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(25))
			Expect(apps.Items).NotTo(BeEmpty())
			for _, a := range apps.Items {
				Expect(a.Client).NotTo(BeNil())
			}

		})

		It("should retrive all the tenant applications by page", func() {
			tenant, _ := CurrentTenant(cred)

			apps, err := tenant.GetApplications(NewPageRequest(1, 0), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(1))
			Expect(apps.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant applications by page and filter", func() {
			tenant, _ := CurrentTenant(cred)

			f := DefaultFilter{Name: "stormpath"}

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Items).To(HaveLen(1))
		})
	})
})
