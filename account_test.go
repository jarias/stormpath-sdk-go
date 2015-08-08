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
			acc := NewAccount("test@test.org", "123", "test@test.org", "test", "test")

			jsonData, _ := json.Marshal(acc)

			Expect(string(jsonData)).To(Equal("{\"username\":\"test@test.org\",\"email\":\"test@test.org\",\"password\":\"123\",\"givenName\":\"test\",\"surname\":\"test\"}"))
		})
	})
	Describe("Save", func() {
		It("should update an existing account", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			account.GivenName = "julio"
			err := account.Save()

			Expect(err).NotTo(HaveOccurred())
			Expect(account.GivenName).To(Equal("julio"))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing account", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			err := account.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AddToGroup", func() {
		It("should add an account to an existing group", func() {
			group := newTestGroup()
			app.CreateGroup(group)

			_, err := account.AddToGroup(group)
			gm, _ := account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(1))
			account.RemoveFromGroup(group)
			group.Delete()
		})
	})

	Describe("RemoveFromGroup", func() {
		It("should remove an account from an existing group", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			var groupCountBefore int
			group := newTestGroup()
			app.CreateGroup(group)

			gm, _ := account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))
			groupCountBefore = len(gm.Items)

			account.AddToGroup(group)

			err := account.RemoveFromGroup(group)
			gm, _ = account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(groupCountBefore))
			group.Delete()
		})
	})

	Describe("GetGroupMemberships", func() {
		It("should allow expanding the account", func() {
			acct := registerTestAccount()
			group := addAccountToGroup(acct)

			groupMemberships, err := acct.GetGroupMemberships(MakeGroupMemershipCriteria().WithAccount().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			for _, gm := range groupMemberships.Items {
				Expect(gm.Account).To(BeEquivalentTo(*acct))
				Expect(gm.Group).NotTo(BeEquivalentTo(*group))
			}
		})

		//It("should allow expanding the group", func() {
		//	account := registerTestAccount()
		//	group := addAccountToGroup(account)
		//
		//	groupMemberships, err := account.GetGroupMemberships(NewDefaultPageRequest(), "group")
		//
		//	Expect(err).NotTo(HaveOccurred())
		//	for _, gm := range groupMemberships.Items {
		//		Expect(gm.Account).To(BeEquivalentTo(account))
		//		Expect(gm.Group.Name).To(BeEquivalentTo(group.Name))
		//	}
		//})
	})

	Describe("GetCustomData", func() {
		It("should retrieve an account custom data", func() {
			customData, err := account.GetCustomData()

			Expect(err).NotTo(HaveOccurred())
			Expect(customData).NotTo(BeEmpty())
		})
	})

	Describe("UpdateCustomData", func() {
		It("should set an account custom data", func() {
			customData, err := account.UpdateCustomData(map[string]interface{}{"custom": "data"})

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("data"))
		})

		It("should update an account custom data", func() {
			account.UpdateCustomData(map[string]interface{}{"custom": "data"})
			customData, err := account.UpdateCustomData(map[string]interface{}{"custom": "nodata"})

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("nodata"))
		})
	})
})
