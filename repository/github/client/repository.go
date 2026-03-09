package client

import (
	"github.com/fikrimohammad/secret-scraper/config"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/google/go-github/v84/github"
)

type repositoryObject struct {
	cli *github.Client
}

func New(cfg *config.Config) repository.GithubClientRepository {
	cli := github.NewClient(nil).WithAuthToken(cfg.Github.AccessToken)
	return &repositoryObject{
		cli: cli,
	}
}
