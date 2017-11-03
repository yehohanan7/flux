package mongodb

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
	mgo "gopkg.in/mgo.v2"
)

var _ = Describe("Mongo offset Store", func() {
	var store OffsetStore
	var session *mgo.Session
	storeID := "test-consumer"
	databaseName := "test"
	collectionName := "offset"

	BeforeSuite(func() {
		sess, err := mgo.Dial("localhost")
		Expect(err).Should(BeNil())
		session = sess
	})

	AfterSuite(func() {
		if session != nil {
			session.Close()
		}
	})

	BeforeEach(func() {
		options := DefaultMongoOffsetStoreOptions()
		options.DatabaseName = databaseName
		options.CollectionName = collectionName
		options.Session = session
		options.StoreID = storeID
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
