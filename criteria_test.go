package stormpath_test

import (
	"net/url"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Criteria", func() {
	Describe("AccountCriteria", func() {
		It("should serialize to an empty string if criteria is empty", func() {
			str := MakeAccountCriteria().ToQueryString()

			Expect(str).To(Equal(""))
		})

		It("should serialize a basic paged criteria", func() {
			str := MakeAccountCriteria().Offset(0).Limit(25).ToQueryString()

			Expect(str).To(Equal("?limit=25&offset=0"))
		})

		Describe("filters", func() {
			It("should allow filter by GivenName", func() {
				str := MakeAccountCriteria().GivenNameEq("test").ToQueryString()

				Expect(str).To(Equal("?givenName=test"))
			})
			It("should allow filter by Surname", func() {
				str := MakeAccountCriteria().SurnameEq("test").ToQueryString()

				Expect(str).To(Equal("?surname=test"))
			})
			It("should allow filter by Email", func() {
				str := MakeAccountCriteria().EmailEq("test").ToQueryString()

				Expect(str).To(Equal("?email=test"))
			})
			It("should allow filter by Username", func() {
				str := MakeAccountCriteria().UsernameEq("test").ToQueryString()

				Expect(str).To(Equal("?username=test"))
			})
			It("should allow filter by MiddelName", func() {
				str := MakeAccountCriteria().MiddleNameEq("test").ToQueryString()

				Expect(str).To(Equal("?middleName=test"))
			})
			It("should allow filter by Status", func() {
				str := MakeAccountCriteria().StatusEq("test").ToQueryString()

				Expect(str).To(Equal("?status=test"))
			})
		})

		Describe("expansion", func() {
			It("should allow Directory expansion", func() {
				str := MakeAccountCriteria().WithDirectory().ToQueryString()

				Expect(str).To(Equal("?expand=directory"))
			})
			It("should allow CustomData expansion", func() {
				str := MakeAccountCriteria().WithCustomData().ToQueryString()

				Expect(str).To(Equal("?expand=customData"))
			})
			It("should allow Tenant expansion", func() {
				str := MakeAccountCriteria().WithTenant().ToQueryString()

				Expect(str).To(Equal("?expand=tenant"))
			})
			It("should allow Groups expansion", func() {
				str := MakeAccountCriteria().WithGroups(DefaultPageRequest).ToQueryString()

				Expect(str).To(Equal("?expand=groups" + url.QueryEscape("(offset:0,limit:25)")))
			})
			It("should allow GroupMemberships expansion", func() {
				str := MakeAccountCriteria().WithGroupMemberships(DefaultPageRequest).ToQueryString()

				Expect(str).To(Equal("?expand=groupMemberships" + url.QueryEscape("(offset:0,limit:25)")))
			})
		})

		It("should allow a multiple attribute expansion, filtering and paging", func() {
			c := MakeAccountCriteria()

			str := c.WithDirectory().WithTenant().UsernameEq("test").Offset(2).Limit(40).ToQueryString()

			Expect(str).To(Equal("?expand=directory%2Ctenant&limit=40&offset=2&username=test"))
		})
	})
})
