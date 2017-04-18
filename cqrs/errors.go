package cqrs

import "errors"

var (
	Conflict = errors.New("There was a conflict with your update, try again")
)
