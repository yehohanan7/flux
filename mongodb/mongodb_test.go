package mongodb

import (
	"encoding/gob"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	mgo "gopkg.in/mgo.v2"

	"testing"
)

var session *mgo.Session

func TestMongodb(t *testing.T) {
	BeforeSuite(func() {
		gob.Register(EventPayload{})
		sess, err := mgo.Dial("localhost")
		Expect(err).Should(BeNil())
		session = sess
	})

	AfterSuite(func() {
		if session != nil {
			session.Close()
		}
	})
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}
