package boltdb

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
)

var _ = Describe("Bolt Offset store", func() {

	var store OffsetStore

	BeforeEach(func() {
		store = NewOffsetStore("offsetstore.db")
	})

	AfterEach(func() {
		os.Remove("offsetstore.db")
	})

	It("Should save offsets", func() {
		err := store.SaveOffset(1)
		Expect(err).To(BeNil())

		offset, err2 := store.GetLastOffset()
		Expect(err2).To(BeNil())
		Expect(offset).To(Equal(1))
	})
})
