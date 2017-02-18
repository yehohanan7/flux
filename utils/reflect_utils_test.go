package utils

import (
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestData struct {
}

func (t *TestData) HandleString(x string)  {}
func (t *TestData) HandleArray(x []string) {}
func (t *TestData) SomeMethod(x int)       {}
func (t *TestData) HandleInt(x int)        {}

var _ = Describe("Reflect Utils", func() {
	It("Should Find methods", func() {
		data := new(TestData)

		methods := FindMethods(data, func(m reflect.Method) bool {
			return strings.HasPrefix(m.Name, "Handle")
		})

		Expect(len(methods)).Should(Equal(3))
	})
})
