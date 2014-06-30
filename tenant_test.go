package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	Describe("CurrentTentant", func() {
		It("should retrive the current tenant", func() {
			tenant, err := CurrentTenant()

			Expect(err).NotTo(HaveOccurred())
			Expect(tenant.Href).NotTo(BeEmpty())
			Expect(tenant.Name).NotTo(BeEmpty())
			Expect(tenant.Key).NotTo(BeEmpty())
			Expect(tenant.Applications.Href).NotTo(BeEmpty())
			Expect(tenant.Directories.Href).NotTo(BeEmpty())
		})
	})

	Describe("GetDirectories", func() {
		It("should retrive all the tenant directories", func() {
			tenant, _ := CurrentTenant()

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(25))
			Expect(directories.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant directories by page", func() {
			tenant, _ := CurrentTenant()

			directories, err := tenant.GetDirectories(NewPageRequest(1, 0), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(1))
			Expect(directories.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant directories by page and filter", func() {
			tenant, _ := CurrentTenant()

			f := DefaultFilter{Name: "Stormpath Administrators"}

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Items).To(HaveLen(1))
		})

	})

	Describe("GetApplications", func() {
		It("should retrive all the tenant applications", func() {
			tenant, _ := CurrentTenant()

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(25))
			Expect(apps.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant applications by page", func() {
			tenant, _ := CurrentTenant()

			apps, err := tenant.GetApplications(NewPageRequest(1, 0), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(1))
			Expect(apps.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant applications by page and filter", func() {
			tenant, _ := CurrentTenant()

			f := DefaultFilter{Name: "stormpath"}

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Items).To(HaveLen(1))
		})
	})
})
