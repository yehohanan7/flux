package feed

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/yehohanan7/flux/cqrs"
)

var _ = Describe("Json Feed Generator", func() {
	It("Should generate json feed", func() {
		var actualFeed JsonEventFeed
		generator := JsonFeedGenerator{}
		event1 := NewEvent("AggregateName", 0, "event payload 1")
		event2 := NewEvent("AggregateName", 1, "event payload 2")

		b := generator.Generate("someurl", "some description", []Event{event1, event2})

		json.Unmarshal(b, &actualFeed)
		Expect(actualFeed.Description).To(Equal("some description"))
		Expect(actualFeed.Events).To(HaveLen(2))
	})
})
