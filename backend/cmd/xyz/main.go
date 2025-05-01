package main

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/glup3/trendingrepos/api"
)

// Go
// var y = [...]int{
//  1_000_000,
// 	4043,
// 	1924,
// 	1177,
// 	811,
// 	616,
// 	482,
// 	396,
// 	331,
// 	283,
// 	246,
// 	215,
// }

var y = [...]int{
	1_000_000,
	3359,
	1863,
	1262,
	924,
	724,
	603,
	504,
	429,
	374,
	330,
	294,
	266,
	241,
	220,
	202,
}

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
	var wg sync.WaitGroup
	count := 0

	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	c := api.NewAPIClient(apiKey)

	for _, maxStars := range y {
		for _, cursor := range cursors {
			wg.Add(1)
			count++

			go func(cursor, language string, maxStars int) {
				defer wg.Done()

				_, err := c.SearchRepos(ctx, cursor, api.QueryArgs{
					MinStars:         200,
					MaxStars:         maxStars,
					Languages:        []string{language},
					IgnoredLanguages: []string{},
				})
				if err != nil {
					slog.Error(
						"failed fetching",
						slog.String("cursor", cursor),
						slog.String("language", language),
						slog.Int("maxStars", maxStars),
						slog.Any("error", err),
					)
				}
			}(cursor, "Java", maxStars)

			if count%100 == 0 {
				slog.Info("cooling down", slog.Int("count", count))
				time.Sleep(time.Second * time.Duration(20))
			}

		}
	}

	wg.Wait()
}
