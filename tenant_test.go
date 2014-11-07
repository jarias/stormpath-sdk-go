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

	//Do to changes in the stormpath plans on a dev account only 1 app can be created so for now I'll comment this test
	//Describe("CreateApplication", func() {
	//	It("should create a new application", func() {
	//		application := NewApplication("create-app")
	//		err := tenant.CreateApplication(application)
	//		application.Purge()
	//
	//		Expect(err).NotTo(HaveOccurred())
	//		Expect(application.Href).NotTo(BeEmpty())
	//	})
	//})

	Describe("CreateDirectory", func() {
		It("should create a new directory", func() {
			dir := NewDirectory("create-dir")
			err := tenant.CreateDirectory(dir)
			dir.Delete()

			Expect(err).NotTo(HaveOccurred())
			Expect(dir.Href).NotTo(BeEmpty())
		})
	})

	Describe("GetDirectories", func() {
		It("should retrive all the tenant directories", func() {
			tenant, _ := CurrentTenant()

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), NewEmptyFilter())

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(25))
			Expect(directories.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant directories by page", func() {
			tenant, _ := CurrentTenant()

			directories, err := tenant.GetDirectories(NewPageRequest(1, 0), NewEmptyFilter())

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(1))
			Expect(directories.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant directories by page and filter", func() {
			tenant, _ := CurrentTenant()

			f := NewDefaultFilter("Stormpath Administrators", "", "")

			directories, err := tenant.GetDirectories(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Items).To(HaveLen(1))
		})

	})

	Describe("GetApplications", func() {
		It("should retrive all the tenant applications", func() {
			tenant, _ := CurrentTenant()

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), NewEmptyFilter())

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(25))
			Expect(apps.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant applications by page", func() {
			tenant, _ := CurrentTenant()

			apps, err := tenant.GetApplications(NewPageRequest(1, 0), NewEmptyFilter())

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(1))
			Expect(apps.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant applications by page and filter", func() {
			tenant, _ := CurrentTenant()

			f := NewDefaultFilter("stormpath", "", "")

			apps, err := tenant.GetApplications(NewDefaultPageRequest(), f)

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Items).To(HaveLen(1))
		})
	})
})
