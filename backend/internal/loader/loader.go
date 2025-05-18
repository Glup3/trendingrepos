package loader

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/glup3/trendingrepos/internal/api"
	"golang.org/x/sync/errgroup"
)

type Loader struct {
	apiClient *api.APIClient
	logger    *slog.Logger
}

func NewLoader(apiClient *api.APIClient, logger *slog.Logger) *Loader {
	return &Loader{
		apiClient: apiClient,
		logger:    logger,
	}
}

func (l *Loader) CollectStarsUpperBounds(ctx context.Context, languages, ignoredLanguages []string) ([]int, error) {
	currStars := MaxStarsCount
	starCounts := []int{currStars}
	for currStars > MinStarsCount {
		repos, err := l.apiClient.SearchRepos(ctx, Cursors[9], api.QueryArgs{
			MinStars:         MinStarsCount,
			MaxStars:         currStars,
			Languages:        languages,
			IgnoredLanguages: ignoredLanguages,
		})
		if err != nil {
			return nil, err
		}
		if len(repos) == 0 {
			break
		}
		currStars = repos[len(repos)-1].Stars
		starCounts = append(starCounts, currStars)
	}

	return starCounts, nil
}

func (l *Loader) LoadRepos(ctx context.Context, languages, ignoredLanguages []string, starsUpperBounds []int, timeoutCount *int) []api.Repo {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		res  []api.Repo
		seen = make(map[string]struct{})
	)

	for _, maxStars := range starsUpperBounds {
		for _, cursor := range Cursors {
			wg.Add(1)
			*timeoutCount++

			go func(cursor string, maxStars int) {
				defer wg.Done()

				repos, err := l.apiClient.SearchRepos(ctx, cursor, api.QueryArgs{
					MinStars:         200,
					MaxStars:         maxStars,
					Languages:        languages,
					IgnoredLanguages: ignoredLanguages,
				})
				if err != nil {
					l.logger.Error(
						"failed fetching",
						slog.String("cursor", cursor),
						slog.Any("language", languages),
						slog.Any("ignoredLanguages", ignoredLanguages),
						slog.Int("maxStars", maxStars),
						slog.Any("error", err),
					)
					return
				}

				mu.Lock()
				for _, repo := range repos {
					if _, exists := seen[repo.Id]; !exists {
						seen[repo.Id] = struct{}{}
						res = append(res, repo)
					}
				}
				mu.Unlock()
			}(cursor, maxStars)
		}

		if *timeoutCount%MaxConcurrentRequests == 0 {
			l.logger.Info("cooling down", slog.Int("count", *timeoutCount))
			time.Sleep(LoadingTimeout)
		}
		l.logger.Info("fetching for star range", slog.Int("maxStars", maxStars))
	}

	wg.Wait()

	return res
}

func (l *Loader) LoadRepos2(ctx context.Context, maxStars int) ([]api.Repo, error) {
	g := new(errgroup.Group)
	var mu sync.Mutex
	var res []api.Repo

	for _, cursor := range Cursors {
		g.Go(func() error {
			repos, err := l.apiClient.SearchRepos(ctx, cursor, api.QueryArgs{
				MinStars:         200,
				MaxStars:         maxStars,
				Languages:        []string{},
				IgnoredLanguages: []string{},
			})
			if err != nil {
				l.logger.Error(
					"failed fetching",
					slog.String("cursor", cursor),
					slog.Int("maxStars", maxStars),
					slog.Any("error", err),
				)
				return err
			}

			mu.Lock()
			res = append(res, repos...)
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return res, nil
}
