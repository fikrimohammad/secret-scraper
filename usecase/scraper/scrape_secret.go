package scraper

import (
	"bufio"
	"context"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/fikrimohammad/secret-scraper/usecase"
)

const (
	maxFileSize    int64 = 0 // 0 = no limit
	maxConcurrency       = 5
)

func (u *useCase) ScrapeSecret(ctx context.Context, params usecase.ScrapeSecretParams) (*usecase.ScrapeSecretResult, error) {
	scraperConfig, err := u.configStaticRepository.GetScraperConfig(ctx, repository.GetScraperConfigParams{
		SecretProvider: params.SecretProvider,
		SecretType:     params.SecretType,
	})
	if err != nil {
		return nil, err
	}

	secretPattern, err := regexp.Compile(scraperConfig.SecretRegexPattern)
	if err != nil {
		return nil, err
	}

	var (
		mu         sync.Mutex
		secretsMap = map[string]model.Secret{}
		secrets    = make([]model.Secret, 0)
		sem        = make(chan struct{}, maxConcurrency)
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

		var wg sync.WaitGroup

		for _, githubCode := range searchCodeResult.Codes {
			wg.Add(1)
			sem <- struct{}{} // acquire semaphore

			go func(code model.GithubCode) {
				defer wg.Done()
				defer func() { <-sem }() // release semaphore

				rawFileContentResult, err := u.githubClientRepository.GetFileRawContent(ctx, repository.GithubGetFileRawContentParams{
					HtmlUrl:     code.HtmlUrl,
					MaxFileSize: maxFileSize,
				})
				if err != nil {
					log.Printf("skipping file %s: %v", code.HtmlUrl, err)
					return
				}

				// Scan line-by-line to avoid loading entire content into regex engine
				scanner := bufio.NewScanner(strings.NewReader(rawFileContentResult.RawFileContent))
				mu.Lock()
				defer mu.Unlock()

				for scanner.Scan() {
					line := scanner.Text()
					for _, match := range secretPattern.FindAllString(line, -1) {
						if _, ok := secretsMap[match]; ok {
							continue
						}

						log.Printf("matched secret: %s", match)

						s := model.Secret{
							Provider: params.SecretProvider,
							Type:     params.SecretType,
							Value:    match,
						}

						secrets = append(secrets, s)
						secretsMap[match] = s
					}
				}
			}(githubCode)
		}

		wg.Wait()
	}

	result := &usecase.ScrapeSecretResult{
		Data: secrets,
	}

	return result, nil
}
