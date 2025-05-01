package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/glup3/trendingrepos/api"
)

//go:embed languages.txt
var fileLanguages string

const (
	maxConcurrent = 100
	timeoutSec    = 20
	minStars      = 200
)

// These are 10 next page cursors for page size of 100
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
	languages := strings.Split(strings.TrimSpace(fileLanguages), "\n")

	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	c := api.NewAPIClient(apiKey)

	var wg sync.WaitGroup
	count := 0

	for _, language := range languages {
		maxStars := 1_000_000

		slog.Info("fetching by language", slog.String("language", language))

		for maxStars > minStars {
			for _, cursor := range cursors[:9] {
				wg.Add(1)
				count++

				slog.Info("count", slog.Int("count", count))

				if count%maxConcurrent == 0 {
					slog.Info("cooling down", slog.Int("count", count))
					time.Sleep(time.Second * time.Duration(timeoutSec))
				}

				go func(cursor, language string, maxStars int) {
					defer wg.Done()

					_, err := c.SearchRepos(ctx, cursor, api.QueryArgs{
						MinStars:         minStars,
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
				}(cursor, language, maxStars)
			}
			wg.Wait()

			count++
			slog.Info("count", slog.Int("count", count))
			repos, err := c.SearchRepos(ctx, cursors[9], api.QueryArgs{
				MinStars:         minStars,
				MaxStars:         maxStars,
				Languages:        []string{language},
				IgnoredLanguages: []string{},
			})
			if err != nil {
				slog.Error(
					"failed fetching last cursor",
					slog.String("cursor", cursors[9]),
					slog.String("language", language),
					slog.Int("maxStars", maxStars),
					slog.Any("error", err),
				)
				return // TODO: implement retry
			}

			if len(repos) == 0 {
				break
			}
			maxStars = repos[len(repos)-1].Stars
		}
	}

	slog.Info("Finished")
}
