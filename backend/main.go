package main

import (
	"context"
	"encoding/csv"
	"log"
	"log/slog"
	"os"
	"strings"

	_ "embed"

	"github.com/glup3/trendingrepos/internal/api"
	"github.com/glup3/trendingrepos/internal/loader"
)

//go:embed languages.txt
var fileLanguages string

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)

	allLanguages := strings.Split(strings.TrimSpace(fileLanguages), "\n")
	timeoutCount := 0
	var allRepos []api.Repo

	type x struct {
		language         string
		starsUpperBounds []int
	}

	var ls []x

	for _, language := range allLanguages {
		languages := []string{language}
		ignoredLanguages := []string{}

		logger.Info("start collecting stars upper bounds", slog.String("lang", language))
		starsUpperBounds, err := l.CollectStarsUpperBounds(ctx, languages, ignoredLanguages)
		if err != nil {
			logger.Error("collecting stars upper bounds failed", slog.String("lang", language), slog.Any("error", err))
			return
		}
		logger.Info("finished collecting", slog.String("lang", language), slog.Any("starsUpperBounds", starsUpperBounds))
		ls = append(ls, x{language: language, starsUpperBounds: starsUpperBounds})
	}

	for _, ll := range ls {
		languages := []string{ll.language}
		ignoredLanguages := []string{}

		logger.Info("start loading repos", slog.String("lang", ll.language))
		repos := l.LoadRepos(ctx, languages, ignoredLanguages, ll.starsUpperBounds, &timeoutCount)
		logger.Info("finished loading repos", slog.String("lang", ll.language), slog.Int("timeoutCount", timeoutCount), slog.Int("repos", len(repos)))

		allRepos = append(allRepos, repos...)
	}

	f, err := os.Create("repos.csv")
	if err != nil {
		log.Fatalln(err)
	}

	w := csv.NewWriter(f)
	for i, repo := range allRepos {
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
