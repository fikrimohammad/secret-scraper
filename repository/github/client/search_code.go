package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/google/go-github/v84/github"
)

const maxRetries = 3

func (ro *repositoryObject) SearchCode(ctx context.Context, params repository.GithubSearchCodeParams) (*repository.GithubSearchCodeResult, error) {
	for attempt := range maxRetries + 1 {
		if err := ro.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		data, resp, err := ro.cli.Search.Code(ctx, params.Query, &github.SearchOptions{
			ListOptions: github.ListOptions{
				PerPage: params.Limit,
				Page:    params.Page,
			},
		})
		if err != nil {
			waitDuration, shouldRetry := rateLimitWait(err)
			if shouldRetry && attempt < maxRetries {
				log.Printf("rate limited, waiting %s before retry (%d/%d)", waitDuration, attempt+1, maxRetries)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(waitDuration):
					continue
				}
			}
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("[%d] %s", resp.StatusCode, resp.Status)
		}

		githubCodes := make([]model.GithubCode, 0)
		for _, cr := range data.CodeResults {
			gc := model.GithubCode{}
			if cr.Name != nil {
				gc.Name = *cr.Name
			}

			if cr.Path != nil {
				gc.Path = *cr.Path
			}

			if cr.HTMLURL != nil {
				gc.HtmlUrl = *cr.HTMLURL
			}

			githubCodes = append(githubCodes, gc)
		}

		result := &repository.GithubSearchCodeResult{
			Codes: githubCodes,
		}

		return result, nil
	}

	return nil, fmt.Errorf("search code: max retries exceeded")
}

// rateLimitWait checks if an error is a GitHub rate limit error and returns
// how long to wait before retrying.
func rateLimitWait(err error) (time.Duration, bool) {
	var rateLimitErr *github.RateLimitError
	if errors.As(err, &rateLimitErr) {
		wait := time.Until(rateLimitErr.Rate.Reset.Time) + time.Second // add 1s buffer
		if wait < time.Second {
			wait = time.Second
		}
		return wait, true
	}

	var abuseErr *github.AbuseRateLimitError
	if errors.As(err, &abuseErr) {
		wait := 30 * time.Second
		if abuseErr.RetryAfter != nil {
			wait = *abuseErr.RetryAfter + time.Second
		}
		return wait, true
	}

	return 0, false
}
