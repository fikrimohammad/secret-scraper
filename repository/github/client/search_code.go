package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/repository"
	"github.com/google/go-github/v84/github"
)

func (ro *repositoryObject) SearchCode(ctx context.Context, params repository.GithubSearchCodeParams) (*repository.GithubSearchCodeResult, error) {
	data, resp, err := ro.cli.Search.Code(ctx, params.Query, &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: params.Limit,
			Page:    params.Page,
		},
	})
	if err != nil {
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
