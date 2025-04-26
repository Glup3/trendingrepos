package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/glup3/trendingrepos/api"
)

var pagination_100_based_cursors = [...]string{
	"",
	"Y3Vyc29yOjEwMA==",
	"Y3Vyc29yOjIwMA==",
	"Y3Vyc29yOjMwMA==",
	"Y3Vyc29yOjQwMA==",
	"Y3Vyc29yOjUwMA==",
	"Y3Vyc29yOjYwMA==",
	"Y3Vyc29yOjcwMA==",
	"Y3Vyc29yOjgwMA==",
	"Y3Vyc29yOjkwMA==",
}

type authedTransport struct {
	wrapped      http.RoundTripper
	apiKey       string
	acceptHeader string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.apiKey)
	req.Header.Set("Accept", t.acceptHeader)
	req.Header.Set("X-Github-Next-Global-ID", "1")
	return t.wrapped.RoundTrip(req)
}

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("PAT_TOKEN")

	c := graphql.NewClient("https://api.github.com/graphql", &http.Client{
		Transport: &authedTransport{
			apiKey:       apiKey,
			wrapped:      http.DefaultTransport,
			acceptHeader: "application/json",
		},
	})

	search := func(start int, end int, cursor string) (*api.SearchReposResponse, error) {
		return api.SearchRepos(ctx, c, fmt.Sprintf("is:public stars:%d..%d", start, end), 100, cursor)
	}

	i := 0
	minStarCount := 200
	maxStarCount := 1_000_000

	for maxStarCount > minStarCount {
		i++
		resp, err := search(minStarCount, maxStarCount, "Y3Vyc29yOjkwMA==")
		if err != nil {
			slog.Error("failed fetching", slog.Any("error", err))
			return
		}

		lastRepo := resp.Search.Edges[len(resp.Search.Edges)-1].Node.(*api.SearchReposSearchSearchResultItemConnectionEdgesSearchResultItemEdgeNodeRepository)
		slog.Info(
			"updating max start count",
			slog.Int("count", i),
			slog.Int("prevMin", minStarCount),
			slog.Int("prevMax", maxStarCount),
			slog.Int("nextMax", lastRepo.StargazerCount),
		)
		maxStarCount = lastRepo.StargazerCount
	}

	slog.Info("finished", slog.Int("count", i))
}
