package api

import (
	"context"
	"fmt"
	"strings"
)

type GitHubRepo = searchReposSearchSearchResultItemConnectionEdgesSearchResultItemEdgeNodeRepository

type QueryArgs struct {
	PageSize int
	MinStars int
	MaxStars int
	Cursor   string
}

func (c *APIClient) SearchRepos(ctx context.Context, queryArgs QueryArgs) ([]GitHubRepo, error) {
	resp, err := searchRepos(ctx, c.gClient, buildQuery(queryArgs), queryArgs.PageSize, queryArgs.Cursor)
	if err != nil {
		return nil, err
	}
	repos := make([]GitHubRepo, 0, len(resp.Search.Edges))
	for _, edge := range resp.Search.Edges {
		repos = append(repos, *edge.Node.(*GitHubRepo))
	}
	return repos, nil
}

func buildQuery(args QueryArgs) string {
	var b strings.Builder
	b.WriteString("is:public ")
	b.WriteString("fork:true ")
	fmt.Fprintf(&b, "stars:%d..%d ", args.MinStars, args.MaxStars)
	return b.String()
}
