package client

import (
	"time"

	"github.com/fikrimohammad/secret-scraper/config"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/google/go-github/v84/github"
	"golang.org/x/time/rate"
)

type repositoryObject struct {
	cli         *github.Client
	rateLimiter *rate.Limiter
}

func New(cfg *config.Config) repository.GithubClientRepository {
	cli := github.NewClient(nil).WithAuthToken(cfg.Github.AccessToken)
	// GitHub allows 5000 requests/hour authenticated (~1.4/sec).
	// Use 1 request/sec with burst of 5 to stay under the limit.
	limiter := rate.NewLimiter(rate.Every(time.Second), 5)
	return &repositoryObject{
		cli:         cli,
		rateLimiter: limiter,
	}
}
