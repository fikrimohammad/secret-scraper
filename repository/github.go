package repository

import (
	"context"

	"github.com/fikrimohammad/secret-scraper/model"
)

type GithubClientRepository interface {
	GetFileRawContent(ctx context.Context, params GithubGetFileRawContentParams) (*GithubGetFileRawContentResult, error)
	SearchCode(ctx context.Context, params GithubSearchCodeParams) (*GithubSearchCodeResult, error)
}

type GithubSearchCodeParams struct {
	Query string
	Limit int
	Page  int
}

type GithubSearchCodeResult struct {
	Codes []model.GithubCode
}

type GithubGetFileRawContentParams struct {
	HtmlUrl string
}

type GithubGetFileRawContentResult struct {
	RawFileContent string
}
