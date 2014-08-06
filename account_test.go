package stormpath_test

import (
	"encoding/json"
	. "github.com/jarias/stormpath-sdk-go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the account required fields", func() {
			acc := NewAccount("test@test.org", "123", "test", "test")

			jsonData, _ := json.Marshal(acc)

			Expect(string(jsonData)).To(Equal("{\"email\":\"test@test.org\",\"password\":\"123\",\"givenName\":\"test\",\"surname\":\"test\"}"))
		})
	})

	Describe("Save", func() {
		It("should update an existing account", func() {
			account := NewAccount("u@test.org", "1234567z!A89", "teset", "test")
			app.RegisterAccount(account)

			account.GivenName = "julio"
			err := account.Save()

			Expect(err).NotTo(HaveOccurred())
			Expect(account.GivenName).To(Equal("julio"))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing account", func() {
			account := NewAccount("d@test.org", "1234567z!A89", "teset", "test")
			app.RegisterAccount(account)

			err := account.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AddToGroup", func() {
		It("should add an account to an existing group", func() {
			group := NewGroup("test-group-for-account")
			app.CreateApplicationGroup(group)

			_, err := account.AddToGroup(group)
			gm, _ := account.GetGroupMemberships(NewDefaultPageRequest())

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(1))
			group.Delete()
		})
	})

	Describe("RemoveFromGroup", func() {
		It("should remove an account from an existing group", func() {
			var groupCountBefore int
			group := NewGroup("test-group-for-account-remove")
			app.CreateApplicationGroup(group)

			account.AddToGroup(group)
			gm, _ := account.GetGroupMemberships(NewDefaultPageRequest())
			groupCountBefore = len(gm.Items)
			err := account.RemoveFromGroup(group)
			gm, _ = account.GetGroupMemberships(NewDefaultPageRequest())

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(groupCountBefore))
			group.Delete()
		})
	})

	Describe("GetCustomData", func() {
		It("should retrieve an account custom data", func() {
			customData, err := account.GetCustomData()

			Expect(err).NotTo(HaveOccurred())
			Expect(customData).NotTo(BeEmpty())
		})
	})

	Describe("SetCustomData", func() {
		It("should set an account custom data", func() {
			err := account.SetCustomData(map[string]string{"custom": "data"})

			customData, _ := account.GetCustomData()

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("data"))
		})

		It("should update an account custom data", func() {
			account.SetCustomData(map[string]string{"custom": "data"})
			err := account.SetCustomData(map[string]string{"custom": "nodata"})
			customData, _ := account.GetCustomData()

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("nodata"))
		})
	})
})
