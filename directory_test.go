package stormpath_test

import (
	"encoding/json"
	. "github.com/jarias/stormpath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Directory", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the directory name", func() {
			directory := NewDirectory("name")

			jsonData, _ := json.Marshal(directory)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})

	Describe("Save", func() {
		It("should create a new directory", func() {
			directory := NewDirectory("new-directory-test")

			err := directory.Save()
			directory.Delete()

			Expect(err).NotTo(HaveOccurred())
			Expect(directory.Href).NotTo(BeEmpty())
			Expect(directory.Status).To(Equal("ENABLED"))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing directory", func() {
			directory := NewDirectory("new-directory-test")
			directory.Save()
			err := directory.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetGroups", func() {
		It("should retrive all directory groups", func() {
			directory := NewDirectory("new-directory-test")
			directory.Save()

			groups, err := directory.GetGroups(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(groups.Href).NotTo(BeEmpty())
			Expect(groups.Offset).To(Equal(0))
			Expect(groups.Limit).To(Equal(25))
			Expect(groups.Items).To(BeEmpty())
			directory.Delete()
		})
	})

	Describe("GetAccounts", func() {
		It("should retrieve all directory accounts", func() {
			directory := NewDirectory("new-directory-test")
			directory.Save()

			accounts, err := directory.GetAccounts(NewDefaultPageRequest(), DefaultFilter{})

			Expect(err).NotTo(HaveOccurred())
			Expect(accounts.Href).NotTo(BeEmpty())
			Expect(accounts.Offset).To(Equal(0))
			Expect(accounts.Limit).To(Equal(25))
			Expect(accounts.Items).To(BeEmpty())
			directory.Delete()
		})
	})
})
