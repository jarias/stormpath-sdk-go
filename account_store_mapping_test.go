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
		It("should return a not found error if the application doesn't exists", func() {
			dir := newTestDirectory()

			tenant.CreateDirectory(dir)

			asm := NewAccountStoreMapping(BaseURL+"applications/7ZSGIObcwO8UosxFGfUaxx", dir.Href)
			err := asm.Save()

			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(400))
			Expect(err.(Error).Code).To(Equal(2014))
		})
		It("should return a not found error if the application doesn't exists", func() {
			asm := NewAccountStoreMapping(app.Href, BaseURL+"directories/7ZSGIObcwO8UosxFGfUaxx")
			err := asm.Save()

			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(400))
			Expect(err.(Error).Code).To(Equal(2014))
		})
	})
})
