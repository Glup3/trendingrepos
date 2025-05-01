package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/glup3/trendingrepos/api"
	"github.com/glup3/trendingrepos/internal/loader"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)

	language := "Go"
	languages := []string{language}
	ignoredLanguages := []string{}

	logger.Info("start collecting stars upper bounds")
	starsUpperBounds, err := l.CollectStarsUpperBounds(ctx, languages, ignoredLanguages)
	if err != nil {
		log.Println("collecting stars upper bounds failed", err)
		return
	}
	logger.Info("finished collecting", slog.Any("starsUpperBounds", starsUpperBounds))

	logger.Info("start loading repos")
	timeoutCount := 0
	repos := l.LoadRepos(ctx, languages, ignoredLanguages, starsUpperBounds, &timeoutCount)
	logger.Info("finished loading repos", slog.Int("timeoutCount", timeoutCount), slog.Int("repos", len(repos)))

	f, err := os.Create(fmt.Sprintf("%s.csv", language))
	if err != nil {
		log.Fatalln(err)
	}

	w := csv.NewWriter(f)
	for i, repo := range repos {
		if i == 0 {
			if err := w.Write(repo.CSVHeader()); err != nil {
				log.Fatalln("error writing csv header:", err)
			}
		}

		if err := w.Write(repo.CSVRecord()); err != nil {
			log.Fatalln("error writing repo to csv:", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
