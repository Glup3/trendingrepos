package loader

import (
	"context"

	"github.com/glup3/trendingrepos/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoService struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewRepoService(pool *pgxpool.Pool) *RepoService {
	return &RepoService{
		pool: pool,
		q:    db.New(pool),
	}
}

func (s *RepoService) Insert(ctx context.Context, repos []Repo) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	params := make([]db.InsertTempRepositoriesParams, 0, len(repos))
	for _, repo := range repos {
		params = append(params, db.InsertTempRepositoriesParams{
			GithubID:        repo.Id,
			NameWithOwner:   repo.NameWithOwner,
			Description:     pgtype.Text{String: repo.Description, Valid: true},
			Stars:           int32(repo.Stars),
			PrimaryLanguage: pgtype.Text{String: repo.PrimaryLanguage, Valid: true},
			IsArchived:      repo.IsArchived,
		})
	}

	qtx := s.q.WithTx(tx)
	err = qtx.CreateTempRepositories(ctx)
	if err != nil {
		return err
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

func (s *RepoService) RefreshViews(ctx context.Context) error {
	_, err := s.pool.Exec(ctx, `CALL refresh_continuous_aggregate ('stars_daily', NULL, NULL)`)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `REFRESH MATERIALIZED VIEW stars_trend_monthly`)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `REFRESH MATERIALIZED VIEW stars_trend_weekly`)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `REFRESH MATERIALIZED VIEW stars_trend_daily`)
	if err != nil {
		return err
	}
	return nil
}
