package main

import (
	"context"
	"fmt"
	"os"

	"github.com/glup3/trendingrepos/api"
)

var cursors = [10]string{
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

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	c := api.NewAPIClient(apiKey)

	maxStars := 1_000_000

	for {
		repos, err := c.SearchRepos(ctx, cursors[9], api.QueryArgs{
			MinStars:         200,
			MaxStars:         maxStars,
			Languages:        []string{"Java"},
			IgnoredLanguages: []string{},
		})
		if err != nil {
			fmt.Println("error happened", err)
		}

		if len(repos) == 0 {
			break
		}

		maxStars = repos[len(repos)-1].Stars
		fmt.Println(maxStars)

		if maxStars <= 200 {
			break
		}
	}
}
