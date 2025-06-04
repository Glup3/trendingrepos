package loader

import "github.com/glup3/trendingrepos/internal/api"

type Repo struct {
	Id              string
	Stars           int
	Description     string
	NameWithOwner   string
	PrimaryLanguage string
	IsArchived      bool
}

func repoFromGitHubRepo(ghRepo api.GitHubRepo) Repo {
	return Repo{
		Id:              ghRepo.Id,
		Stars:           ghRepo.StargazerCount,
		Description:     ghRepo.Description,
		NameWithOwner:   ghRepo.NameWithOwner,
		PrimaryLanguage: ghRepo.PrimaryLanguage.Name,
		IsArchived:      ghRepo.IsArchived,
	}
}
