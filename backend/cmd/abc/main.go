package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/glup3/trendingrepos/internal/api"
	"github.com/glup3/trendingrepos/internal/csv"
	"github.com/glup3/trendingrepos/internal/loader"
	"golang.org/x/sync/errgroup"
)

//go:embed stars2.txt
var starsBounds string

const (
	sleepTimeout       = 90 * time.Second
	concurrentRequests = 20 // multiplied by 10 cursors
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

	starsBoundsString := strings.Split(strings.TrimSpace(starsBounds), "\n")

	starsBounds := make([]int, len(starsBoundsString))
	for i, v := range starsBoundsString {
		s, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		starsBounds[i] = s
	}

	g := new(errgroup.Group)
	var allRepos []api.Repo
	var mu sync.Mutex
	i := 0
	for i < len(starsBounds) {
		batchSize := concurrentRequests
		if i+batchSize > len(starsBounds) {
			batchSize = len(starsBounds) - i
		}

		for j := range batchSize {
			g.Go(func() error {
				maxStars := starsBounds[i+j]
				logger.Info("fetching repos", slog.Int("maxStars", maxStars))
				repos, err := l.LoadRepos2(ctx, maxStars)
				if err != nil {
					return err
				}
				mu.Lock()
				allRepos = append(allRepos, repos...)
				mu.Unlock()
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			logger.Warn("failed fetching", slog.Any("error", err))
			time.Sleep(sleepTimeout)
			continue
		}

		time.Sleep(sleepTimeout)
		i += batchSize
	}

	return csv.ToCSV("repos", allRepos)
}
