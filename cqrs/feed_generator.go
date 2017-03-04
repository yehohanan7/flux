package cqrs

type FeedGenerator interface {
	Generate(string, string, []Event) string
}
