package scraper

import (
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/fikrimohammad/secret-scraper/usecase"
)

type useCase struct {
	configStaticRepository repository.ConfigStaticRepository
	githubClientRepository repository.GithubClientRepository
}

func New(configStaticRepository repository.ConfigStaticRepository, githubClientRepository repository.GithubClientRepository) usecase.Scraper {
	return &useCase{
		configStaticRepository: configStaticRepository,
		githubClientRepository: githubClientRepository,
	}
}
