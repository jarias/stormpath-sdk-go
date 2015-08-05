package stormpath_test

import (
	"encoding/json"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Directory", func() {
	Describe("Validate", func() {
		It("should return true if the directory is valid", func() {
			ok, err := NewDirectory("test").Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})
		It("should return false if directory is invalid", func() {
			invalidDirs := []*Directory{
				&Directory{},
				&Directory{Name: string256},
				&Directory{Name: "name", Description: string1001},
			}

			for _, dir := range invalidDirs {
				ok, err := dir.Validate()

				Expect(err).To(HaveOccurred())
				Expect(ok).To(BeFalse())
			}
		})
	})
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the directory name", func() {
			directory := NewDirectory("name")

			jsonData, _ := json.Marshal(directory)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing directory", func() {
			directory := newTestDirectory()

			tenant.CreateDirectory(directory)
			err := directory.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetGroups", func() {
		It("should retrive all directory groups", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			groups, err := directory.GetGroups(NewDefaultPageRequest(), NewEmptyFilter())

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
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			accounts, err := directory.GetAccounts(NewDefaultPageRequest(), NewEmptyFilter())

			Expect(err).NotTo(HaveOccurred())
			Expect(accounts.Href).NotTo(BeEmpty())
			Expect(accounts.Offset).To(Equal(0))
			Expect(accounts.Limit).To(Equal(25))
			Expect(accounts.Items).To(BeEmpty())
			directory.Delete()
		})
	})

	Describe("CreateGroup", func() {
		It("should create new group", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			group := NewGroup("new-group")
			err := directory.CreateGroup(group)

			Expect(err).NotTo(HaveOccurred())
			Expect(group.Href).NotTo(BeEmpty())
			directory.Delete()
		})
	})

	Describe("RegisterAccount", func() {
		It("should create a new accout for the group", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			account := newTestAccount()
			err := directory.RegisterAccount(account)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.Href).NotTo(BeEmpty())
			directory.Delete()
		})
	})
})
