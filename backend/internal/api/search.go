package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

type Repo struct {
	Id              string
	Stars           int
	Description     string
	NameWithOwner   string
	PrimaryLanguage string
	IsArchived      bool
}

type QueryArgs struct {
	PageSize int
	MinStars int
	MaxStars int
	Cursor   string
}

func (c *APIClient) SearchRepos(ctx context.Context, queryArgs QueryArgs) ([]Repo, error) {
	resp, err := searchRepos(ctx, c.gClient, buildQuery(queryArgs), queryArgs.PageSize, queryArgs.Cursor)
	if err != nil {
		return nil, err
	}
	repos := make([]Repo, 0, len(resp.Search.Edges))
	for _, edge := range resp.Search.Edges {
		repo := edge.Node.(*searchReposSearchSearchResultItemConnectionEdgesSearchResultItemEdgeNodeRepository)
		repos = append(repos, Repo{
			Id:              repo.Id,
			Stars:           repo.StargazerCount,
			Description:     repo.Description,
			NameWithOwner:   repo.NameWithOwner,
			PrimaryLanguage: repo.PrimaryLanguage.Name,
			IsArchived:      repo.IsArchived,
		})
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

func (r Repo) CSVRecord() []string {
	return []string{
		r.Id,
		r.NameWithOwner,
		strconv.Itoa(r.Stars),
		r.PrimaryLanguage,
		r.Description,
		strconv.FormatBool(r.IsArchived),
	}
}

func (r Repo) CSVHeader() []string {
	return []string{
		"Id",
		"NameWithOwner",
		"Stars",
		"PrimaryLanguage",
		"Description",
		"Archived",
	}
}
