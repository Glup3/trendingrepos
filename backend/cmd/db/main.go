package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/glup3/trendingrepos/internal/api"
	"github.com/glup3/trendingrepos/internal/csv"
	"github.com/glup3/trendingrepos/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	repos, err := csv.ReadCsvFile("./repos.csv")
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

	c := cron.New()
	c.AddFunc("@every 8s", func() {
		err := logic(ctx, pool, repos)
		if err != nil {
			fmt.Println(err)
		}
	})
	c.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		c.Stop()
	}()
	wg.Wait()
	return nil
}

func logic(ctx context.Context, pool *pgxpool.Pool, repos []api.Repo) error {
	fmt.Println("doing db stuff")
	defer fmt.Println("finished db stuff")
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
