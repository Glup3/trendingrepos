package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/glup3/trendingrepos/api"
	"github.com/glup3/trendingrepos/internal/loader"
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)

	languages := []string{"Python"}
	ignoredLanguages := []string{}

	starsUpperBounds, err := l.CollectStarsUpperBounds(ctx, languages, ignoredLanguages)
	if err != nil {
		log.Println("collecting stars upper bounds failed", err)
		return
	}

	log.Println(starsUpperBounds)

	l.LoadRepos(ctx, languages, ignoredLanguages, starsUpperBounds)
}
