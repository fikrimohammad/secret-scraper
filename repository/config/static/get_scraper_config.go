package static

import (
	"context"

	"github.com/fikrimohammad/secret-scraper/repository"
)

func (ro *repositoryObject) GetScraperConfig(ctx context.Context, params repository.GetScraperConfigParams) (*repository.GetScraperConfigResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	scraperConfigKey := buildConfigKey(params.SecretProvider, params.SecretType)
	scraperConfig, ok := ro.config[scraperConfigKey]
	if !ok {
		return nil, repository.ScraperConfigNotFound
	}

	result := &repository.GetScraperConfigResult{
		SecretQueryKeyword: scraperConfig.SecretQueryKeyword,
		SecretRegexPattern: scraperConfig.SecretRegexPattern,
	}

	return result, nil
}
