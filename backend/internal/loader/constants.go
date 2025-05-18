package loader

import "time"

const (
	PageSize              = 100
	MaxStarsCount         = 1_000_000
	MinStarsCount         = 200
	MaxConcurrentRequests = 20
	SleepTimeout          = time.Second * time.Duration(90)
)

// These are 10 next page cursors for a 100 pageSize
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
