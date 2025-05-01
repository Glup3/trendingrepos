package loader

import "time"

const (
	MaxStarsCount         = 1_000_000
	MinStarsCount         = 200
	MaxConcurrentRequests = 100
	LoadingTimeout        = time.Second * time.Duration(20)
)

// These are 10 next page cursors for page size of 100
var Cursors = [10]string{
	"",
	"Y3Vyc29yOjEwMA==",
	"Y3Vyc29yOjIwMA==",
	"Y3Vyc29yOjMwMA==",
	"Y3Vyc29yOjQwMA==",
	"Y3Vyc29yOjUwMA==",
	"Y3Vyc29yOjYwMA==",
	"Y3Vyc29yOjcwMA==",
	"Y3Vyc29yOjgwMA==",
	"Y3Vyc29yOjkwMA==",
}
