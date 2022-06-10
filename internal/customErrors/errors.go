package customErrors

import "errors"

var (
	NoAuthor          = errors.New("author of post not found")
	AnotherThreadPost = errors.New("parent post was created in another thread")
)
