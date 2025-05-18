package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/glup3/trendingrepos/internal/api"
	"github.com/glup3/trendingrepos/internal/csv"
	"github.com/glup3/trendingrepos/internal/loader"
)

//go:embed stars2.txt
var starsBounds string

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

	starsBoundsString := strings.Split(strings.TrimSpace(starsBounds), "\n")

	starsBounds := make([]int, len(starsBoundsString))
	for i, v := range starsBoundsString {
		s, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		starsBounds[i] = s
	}

	allRepos := l.LoadMultipleRepos(ctx, starsBounds)
	return csv.ToCSV("repos", allRepos)
}
