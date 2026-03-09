package usecase

import (
	"context"

	"github.com/fikrimohammad/secret-scraper/model"
)

type Scraper interface {
	ScrapeSecret(ctx context.Context, params ScrapeSecretParams) (*ScrapeSecretResult, error)
}

type ScrapeSecretParams struct {
	SecretProvider  model.SecretProvider
	SecretType      model.SecretType
	MaxIterations   int
	MaxLimitPerIter int
}

type ScrapeSecretResult struct {
	Data []model.Secret
}
