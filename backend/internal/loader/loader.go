package loader

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/glup3/trendingrepos/internal/api"
	"golang.org/x/sync/errgroup"
)

type apiClient interface {
	SearchRepos(ctx context.Context, queryArgs api.QueryArgs) ([]api.GitHubRepo, error)
}

type Loader struct {
	apiClient apiClient
	logger    *slog.Logger
}

func NewLoader(apiClient apiClient, logger *slog.Logger) *Loader {
	return &Loader{
		apiClient: apiClient,
		logger:    logger,
	}
}

func (l *Loader) CollectStarsUpperBounds(ctx context.Context, languages, ignoredLanguages []string) ([]int, error) {
	currStars := MaxStarsCount
	starCounts := []int{currStars}
	fmt.Printf("%d,\n", currStars)
	for currStars > MinStarsCount {
		ghRepos, err := l.apiClient.SearchRepos(ctx, api.QueryArgs{
			PageSize: PageSize,
			Cursor:   Cursors[9],
			MinStars: MinStarsCount,
			MaxStars: currStars,
		})
		if err != nil {
			return nil, err
		}
		if len(ghRepos) == 0 {
			break
		}

		repos := make([]Repo, 0, len(ghRepos))
		for _, ghRepo := range ghRepos {
			repos = append(repos, repoFromGitHubRepo(ghRepo))
		}

		currStars = repos[len(repos)-1].Stars
		fmt.Printf("%d,\n", currStars)
		starCounts = append(starCounts, currStars)
	}

	return starCounts, nil
}

func (l *Loader) LoadRepos(ctx context.Context, maxStars int) ([]Repo, error) {
	var mu sync.Mutex
	g := new(errgroup.Group)
	repos := make([]Repo, 0, len(Cursors)*PageSize)

	for _, cursor := range Cursors {
		g.Go(func() error {
			ghRepos, err := l.apiClient.SearchRepos(ctx, api.QueryArgs{
				PageSize: PageSize,
				Cursor:   cursor,
				MinStars: MinStarsCount,
				MaxStars: maxStars,
			})
			if err != nil {
				return err
			}

			mu.Lock()
			for _, ghRepo := range ghRepos {
				repos = append(repos, repoFromGitHubRepo(ghRepo))
			}
			mu.Unlock()
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return repos, nil
}

func (l *Loader) LoadMultipleRepos(ctx context.Context, maxStarss []int) []Repo {
	g := new(errgroup.Group)
	var allRepos []Repo
	var mu sync.Mutex
	seen := make(map[string]struct{})
	retryCount := 0

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
					return fmt.Errorf("batchSize %d failed: %w", i, err)
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
			if retryCount >= MaxRetries {
				l.logger.Error("exceeded max retries - will stop fetching", slog.Any("error", err))
				break
			}

			retryCount++
			time.Sleep(ErrorSleepTimeout)
			continue
		}

		i += batchSize
		if i < len(maxStarss) {
			time.Sleep(SleepTimeout)
		}
	}
	return allRepos
}
