package main

import (
	"context"
	"fmt"
	"os"

	"github.com/glup3/trendingrepos/internal/csv"
	"github.com/glup3/trendingrepos/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	repos, err := csv.ReadCsvFile("./repos.csv")
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer pool.Close()

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
