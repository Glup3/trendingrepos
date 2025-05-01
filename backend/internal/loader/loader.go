package loader

import (
	"context"

	"github.com/glup3/trendingrepos/api"
)

type Loader struct {
	apiClient *api.APIClient
}

func NewLoader(apiClient *api.APIClient) *Loader {
	return &Loader{apiClient: apiClient}
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
