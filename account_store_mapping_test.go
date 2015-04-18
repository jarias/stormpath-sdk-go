package stormpath_test

import (
	"encoding/json"

	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountStoreMapping", func() {
	Describe("JSON", func() {
		It("should marshall a minimum JSON with only the required fields", func() {
			accountStoreMapping := NewAccountStoreMapping("http://appurl", "http://storeUrl")

			jsonData, _ := json.Marshal(accountStoreMapping)

			Expect(string(jsonData)).To(Equal("{\"application\":{\"href\":\"http://appurl\"},\"accountStore\":{\"href\":\"http://storeUrl\"}}"))
		})
	})

	Describe("Save", func() {
		It("should create a new account store mapping", func() {
			dir := newTestDirectory()

			tenant.CreateDirectory(dir)

			asm := NewAccountStoreMapping(app.Href, dir.Href)
			err := asm.Save()

			Expect(err).NotTo(HaveOccurred())
			Expect(asm.Href).NotTo(BeEmpty())
		})
	})
})
