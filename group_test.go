package stormpath_test

import (
	"encoding/json"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the directory name", func() {
			group := NewGroup("name")

			jsonData, _ := json.Marshal(group)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})
})
