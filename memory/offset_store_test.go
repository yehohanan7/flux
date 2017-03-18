package memory

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
)

var _ = Describe("Offset Store", func() {
	var store OffsetStore

	BeforeEach(func() {
		store = NewOffsetStore()
	})

	It("Should have 0 as default offset", func() {
		offset, err := store.GetLastOffset()

		Expect(err).Should(BeNil())
		Expect(offset).Should(Equal(0))
	})

	It("Should save offsets", func() {
		store.SaveOffset(10)

		offset, err := store.GetLastOffset()

		Expect(err).Should(BeNil())
		Expect(offset).Should(Equal(10))
	})
})
