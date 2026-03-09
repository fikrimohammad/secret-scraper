package client

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/fikrimohammad/secret-scraper/repository"
)

func (ro *repositoryObject) GetFileRawContent(ctx context.Context, params repository.GithubGetFileRawContentParams) (*repository.GithubGetFileRawContentResult, error) {
	rawFileUrl := strings.Replace(params.HtmlUrl, "github.com", "raw.githubusercontent.com", 1)
	rawFileUrl = strings.Replace(rawFileUrl, "/blob/", "/", 1)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawFileUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := ro.cli.Client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawFileContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &repository.GithubGetFileRawContentResult{
		RawFileContent: string(rawFileContent),
	}

	return result, nil
}
