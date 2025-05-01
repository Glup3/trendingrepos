package loader

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/glup3/trendingrepos/api"
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

func (l *Loader) LoadRepos(ctx context.Context, languages, ignoredLanguages []string, starsUpperBounds []int) []api.Repo {
	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		res  []api.Repo
		seen = make(map[string]struct{})
	)
	count := 0

	for _, maxStars := range starsUpperBounds {
		for _, cursor := range Cursors {
			wg.Add(1)
			count++

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

			if count%MaxConcurrentRequests == 0 {
				l.logger.Info("cooling down", slog.Int("count", count))
				time.Sleep(LoadingTimeout)
			}
		}
	}

	wg.Wait()

	return res
}
