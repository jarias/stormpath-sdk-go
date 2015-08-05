package stormpath_test

import (
	"encoding/json"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group", func() {
	Describe("Validate", func() {
		It("should return true if the group is valid", func() {
			ok, err := NewGroup("test").Validate()

			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(BeTrue())
		})
		It("should return false if group is invalid", func() {
			invalidGroups := []*Group{
				&Group{},
				&Group{Name: string256},
				&Group{Name: "name", Description: string1001},
			}

			for _, group := range invalidGroups {
				ok, err := group.Validate()

				Expect(err).To(HaveOccurred())
				Expect(ok).To(BeFalse())
			}
		})
	})
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the directory name", func() {
			group := NewGroup("name")

			jsonData, _ := json.Marshal(group)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})
})
