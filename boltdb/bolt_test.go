package boltdb

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCqrs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Boltdb Suite")
}
