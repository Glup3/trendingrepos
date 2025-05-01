package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"
)

const pageSize = 100

type APIClient struct {
	gClient graphql.Client
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

func NewAPIClient(apiKey string) *APIClient {
	c := graphql.NewClient("https://api.github.com/graphql", &http.Client{
		Transport: &authedTransport{
			apiKey:       apiKey,
			wrapped:      http.DefaultTransport,
			acceptHeader: "application/json",
		},
	})

	return &APIClient{
		gClient: c,
	}
}

type Repo struct {
	Id              string
	Stars           int
	Description     string
	NameWithOwner   string
	PrimaryLanguage string
}

type QueryArgs struct {
	MinStars         int
	MaxStars         int
	Languages        []string
	IgnoredLanguages []string
}

func (c *APIClient) SearchRepos(ctx context.Context, cursor string, queryArgs QueryArgs) ([]Repo, error) {
	repos := make([]Repo, 0, pageSize)

	resp, err := searchRepos(ctx, c.gClient, buildQuery(queryArgs), pageSize, cursor)
	if err != nil {
		return nil, err
	}

	for _, edge := range resp.Search.Edges {
		repo := edge.Node.(*searchReposSearchSearchResultItemConnectionEdgesSearchResultItemEdgeNodeRepository)
		repos = append(repos, Repo{
			Id:              repo.Id,
			Stars:           repo.StargazerCount,
			Description:     repo.Description,
			NameWithOwner:   repo.NameWithOwner,
			PrimaryLanguage: repo.PrimaryLanguage.Name,
		})
	}

	return repos, nil
}

func buildQuery(args QueryArgs) string {
	var b strings.Builder
	b.WriteString("is:public ")
	b.WriteString("fork:true ")
	fmt.Fprintf(&b, "stars:%d..%d ", args.MinStars, args.MaxStars)
	for _, lang := range args.Languages {
		fmt.Fprintf(&b, "language:%s ", lang)
	}
	for _, lang := range args.IgnoredLanguages {
		fmt.Fprintf(&b, "-language:%s ", lang)
	}
	return b.String()
}
