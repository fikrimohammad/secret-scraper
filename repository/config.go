package repository

import (
	"context"
	"errors"

	"github.com/fikrimohammad/secret-scraper/model"
)

var (
	ScraperConfigNotFound = errors.New("scraper config not found")
)

type ConfigStaticRepository interface {
	GetScraperConfig(ctx context.Context, params GetScraperConfigParams) (*GetScraperConfigResult, error)
}

type GetScraperConfigParams struct {
	SecretProvider model.SecretProvider
	SecretType     model.SecretType
}

type GetScraperConfigResult struct {
	SecretQueryKeyword string
	SecretRegexPattern string
}
