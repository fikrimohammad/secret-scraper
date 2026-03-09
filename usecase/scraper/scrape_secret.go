package scraper

import (
	"context"
	"log"
	"regexp"

	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/fikrimohammad/secret-scraper/usecase"
)

func (u *useCase) ScrapeSecret(ctx context.Context, params usecase.ScrapeSecretParams) (*usecase.ScrapeSecretResult, error) {
	scraperConfig, err := u.configStaticRepository.GetScraperConfig(ctx, repository.GetScraperConfigParams{
		SecretProvider: params.SecretProvider,
		SecretType:     params.SecretType,
	})
	if err != nil {
		return nil, err
	}

	var (
		secretsMap = map[string]model.Secret{}
		secrets    = make([]model.Secret, 0)
	)

	for currentIter := range params.MaxIterations {
		searchCodeResult, err := u.githubClientRepository.SearchCode(ctx, repository.GithubSearchCodeParams{
			Query: scraperConfig.SecretQueryKeyword,
			Limit: params.MaxLimitPerIter,
			Page:  currentIter + 1,
		})
		if err != nil {
			return nil, err
		}

		if len(searchCodeResult.Codes) == 0 {
			break
		}

		for _, githubCode := range searchCodeResult.Codes {
			rawFileContentResult, err := u.githubClientRepository.GetFileRawContent(ctx, repository.GithubGetFileRawContentParams{
				HtmlUrl: githubCode.HtmlUrl,
			})
			if err != nil {
				return nil, err
			}

			secretPattern, err := regexp.Compile(scraperConfig.SecretRegexPattern)
			if err != nil {
				return nil, err
			}

			for _, match := range secretPattern.FindAll([]byte(rawFileContentResult.RawFileContent), -1) {
				rawSecret := string(match)
				if _, ok := secretsMap[rawSecret]; ok {
					continue
				}

				log.Printf("matched secret: %s", rawSecret)

				s := model.Secret{
					Provider: params.SecretProvider,
					Type:     params.SecretType,
					Value:    rawSecret,
				}

				secrets = append(secrets, s)
				secretsMap[rawSecret] = s
			}
		}
	}

	result := &usecase.ScrapeSecretResult{
		Data: secrets,
	}

	return result, nil
}
