package mongodb

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
)

var _ = Describe("Mongo offset Store", func() {
	var store OffsetStore
	storeId := "test-consumer"
	databaseName := "test"
	collectionName := "offset"

	BeforeEach(func() {
		options := DefaultMongoOffsetStoreOptions()
		options.Database = databaseName
		options.Collection = collectionName
		options.Session = session
		options.StoreId = storeId
		store = NewOffsetStore(options)
	})

	AfterEach(func() {
		session.DB(databaseName).C(collectionName).DropCollection()
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
