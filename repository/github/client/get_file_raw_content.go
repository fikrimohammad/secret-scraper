package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fikrimohammad/secret-scraper/repository"
)

func (ro *repositoryObject) GetFileRawContent(ctx context.Context, params repository.GithubGetFileRawContentParams) (*repository.GithubGetFileRawContentResult, error) {
	if err := ro.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

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

	// Skip files that exceed the max size based on Content-Length header
	if params.MaxFileSize > 0 && resp.ContentLength > params.MaxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d)", resp.ContentLength, params.MaxFileSize)
	}

	// Cap reading to MaxFileSize to prevent OOM on files without Content-Length
	reader := io.Reader(resp.Body)
	if params.MaxFileSize > 0 {
		reader = io.LimitReader(resp.Body, params.MaxFileSize)
	}

	rawFileContent, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	result := &repository.GithubGetFileRawContentResult{
		RawFileContent: string(rawFileContent),
	}

	return result, nil
}
