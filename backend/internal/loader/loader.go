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
		repos, err := l.apiClient.SearchRepos(ctx, api.QueryArgs{
			PageSize: PageSize,
			Cursor:   Cursors[9],
			MinStars: MinStarsCount,
			MaxStars: currStars,
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

func (l *Loader) LoadRepos(ctx context.Context, maxStars int) ([]api.Repo, error) {
	g := new(errgroup.Group)
	var mu sync.Mutex
	var res []api.Repo

	for _, cursor := range Cursors {
		g.Go(func() error {
			repos, err := l.apiClient.SearchRepos(ctx, api.QueryArgs{
				PageSize: PageSize,
				Cursor:   cursor,
				MinStars: MinStarsCount,
				MaxStars: maxStars,
			})
			if err != nil {
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

func (l *Loader) LoadMultipleRepos(ctx context.Context, maxStarss []int) []api.Repo {
	g := new(errgroup.Group)
	var allRepos []api.Repo
	var mu sync.Mutex
	seen := make(map[string]struct{})

	i := 0
	for i < len(maxStarss) {
		batchSize := MaxConcurrentRequests
		if i+batchSize > len(maxStarss) {
			batchSize = len(maxStarss) - i
		}

		l.logger.Info("fetching repos batch", slog.Int("i", i), slog.Any("maxStarss", maxStarss[i:i+batchSize]))
		for j := range batchSize {
			g.Go(func() error {
				maxStars := maxStarss[i+j]
				repos, err := l.LoadRepos(ctx, maxStars)
				if err != nil {
					return err
				}
				mu.Lock()
				for _, repo := range repos {
					if _, exists := seen[repo.Id]; !exists {
						seen[repo.Id] = struct{}{}
						allRepos = append(allRepos, repo)
					}
				}
				mu.Unlock()
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			l.logger.Warn("failed fetching - will sleep", slog.Any("error", err))
			time.Sleep(ErrorSleepTimeout)
			continue
		}

		time.Sleep(SleepTimeout)
		i += batchSize
	}
	return allRepos
}
