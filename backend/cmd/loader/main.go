package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"

	"github.com/glup3/trendingrepos/internal/api"
	"github.com/glup3/trendingrepos/internal/loader"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

//go:embed stars.txt
var starsBounds string

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(ctx, logger); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	apiKey := os.Getenv("PAT_TOKEN")
	apiClient := api.NewAPIClient(apiKey)
	l := loader.NewLoader(apiClient, logger)
	c := cron.New()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	repoService := loader.NewRepoService(pool)

	maxStarss, err := loadMaxStarss()
	if err != nil {
		return err
	}

	c.AddFunc("0 * * * *", func() {
		repos := l.LoadMultipleRepos(ctx, maxStarss)
		logger.Info("finished loading repos - persisting now", slog.Int("repos", len(repos)))
		err := repoService.Insert(ctx, repos)
		if err != nil {
			logger.Error("persisting data failed", slog.Any("error", err))
		}
		logger.Info("finished persisting data - refreshing views")
		err = repoService.RefreshViews(ctx)
		if err != nil {
			logger.Error("failed refreshing views", slog.Any("error", err))
		}
		logger.Info("finished refreshing views")
	})
	c.Start()
	logger.Info("application started")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		c.Stop()
		logger.Info("stopped application")
	}()
	wg.Wait()
	return nil
}

func loadMaxStarss() ([]int, error) {
	starsBoundsString := strings.Split(strings.TrimSpace(starsBounds), "\n")
	starsBounds := make([]int, len(starsBoundsString))
	for i, v := range starsBoundsString {
		s, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		starsBounds[i] = s
	}
	return starsBounds, nil
}
