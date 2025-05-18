package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/glup3/trendingrepos/api"
	"github.com/glup3/trendingrepos/internal/loader"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	apiKey := os.Getenv("PAT_TOKEN")

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)

	var languages, ignoredLanguages []string
	starsUpperBounds, err := l.CollectStarsUpperBounds(ctx, languages, ignoredLanguages)
	if err != nil {
		return fmt.Errorf("collecting stars upper bounds failed: %v", err)
	}
	logger.Info("upperBounds", slog.Int("len", len(starsUpperBounds)))

	return nil
}
