package cqrs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Offset Store", func() {
	var store OffsetStore

	BeforeEach(func() {
		store = NewInMemoryOffsetStore()
	})

	It("Should have -1 as default offset", func() {
		offset, err := store.GetLastOffset()

		Expect(err).Should(BeNil())
		Expect(offset).Should(Equal(-1))
	})

	It("Should save offsets", func() {
		store.SaveOffset(10)

		offset, err := store.GetLastOffset()

		Expect(err).Should(BeNil())
		Expect(offset).Should(Equal(10))
	})
})
