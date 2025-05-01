package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

const pageSize = 100

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

func (r Repo) CSVRecord() []string {
	return []string{
		r.Id,
		r.NameWithOwner,
		strconv.Itoa(r.Stars),
		r.PrimaryLanguage,
		r.Description,
	}
}

func (r Repo) CSVHeader() []string {
	return []string{
		"Id",
		"NameWithOwner",
		"Stars",
		"PrimaryLanguage",
		"Description",
	}
}
