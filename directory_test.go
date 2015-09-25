package stormpath_test

import (
	"encoding/json"

	. "github.com/jarias/stormpath-sdk-go"
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

	Describe("Delete", func() {
		It("should delete an existing directory", func() {
			directory := newTestDirectory()

			tenant.CreateDirectory(directory)
			err := directory.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetAccountCreationPolicy", func() {
		It("should retrive the directory account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, err := directory.GetAccountCreationPolicy()

			Expect(err).NotTo(HaveOccurred())
			Expect(policy).To(Equal(directory.AccountCreationPolicy))
			Expect(policy.VerificationEmailStatus).To(Equal("DISABLED"))
			Expect(policy.VerificationSuccessEmailStatus).To(Equal("DISABLED"))
			Expect(policy.WelcomeEmailStatus).To(Equal("DISABLED"))
		})
	})

	Describe("GetGroups", func() {
		It("should retrive all directory groups", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			groups, err := directory.GetGroups(MakeGroupCriteria())

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

			accounts, err := directory.GetAccounts(MakeAccountCriteria())

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
