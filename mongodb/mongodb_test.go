package mongodb

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMongodb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongodb Suite")
}
