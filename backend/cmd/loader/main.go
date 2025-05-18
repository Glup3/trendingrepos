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
	"github.com/glup3/trendingrepos/internal/db"
	"github.com/glup3/trendingrepos/internal/loader"
	"github.com/jackc/pgx/v5/pgtype"
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

	maxStarss, err := loadMaxStarss()
	if err != nil {
		return err
	}

	c.AddFunc("0 * * * *", func() {
		repos := l.LoadMultipleRepos(ctx, maxStarss)
		err := logic(ctx, pool, repos)
		if err != nil {
			logger.Error("persisting data failed", slog.Any("erro", err))
		}
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

func logic(ctx context.Context, pool *pgxpool.Pool, repos []api.Repo) error {
	queries := db.New(pool)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)
	err = qtx.CreateTempRepositories(ctx)
	if err != nil {
		return err
	}
	params := make([]db.InsertTempRepositoriesParams, len(repos))
	for i, repo := range repos {
		params[i] = db.InsertTempRepositoriesParams{
			GithubID:        repo.Id,
			NameWithOwner:   repo.NameWithOwner,
			Description:     pgtype.Text{String: repo.Description, Valid: true},
			Stars:           int32(repo.Stars),
			PrimaryLanguage: pgtype.Text{String: repo.PrimaryLanguage, Valid: true},
		}
	}
	_, err = qtx.InsertTempRepositories(ctx, params)
	if err != nil {
		return err
	}
	err = qtx.InsertRepositories(ctx)
	if err != nil {
		return err
	}
	err = qtx.InsertStars(ctx)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
