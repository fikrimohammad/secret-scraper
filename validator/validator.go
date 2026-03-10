package validator

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/fikrimohammad/secret-scraper/model"
)

var client = &http.Client{Timeout: 10 * time.Second}

type testFunc func(ctx context.Context, key string) model.SecretStatus

var providers = map[model.SecretProvider]testFunc{
	"anthropic": testAnthropic,
	"openai":    testOpenAI,
	"github":    testGitHub,
	"slack":     testSlack,
	"stripe":    testStripe,
	"sendgrid":  testSendGrid,
}

// TestSecrets validates a slice of secrets in-place, setting the Status field.
func TestSecrets(ctx context.Context, secrets []model.Secret) {
	for i := range secrets {
		fn, ok := providers[secrets[i].Provider]
		if !ok {
			secrets[i].Status = model.SecretStatusUntested
			continue
		}

		log.Printf("testing %s/%s key: %s...", secrets[i].Provider, secrets[i].Type, secrets[i].Value[:12])
		secrets[i].Status = fn(ctx, secrets[i].Value)
		log.Printf("  result: %s", secrets[i].Status)
	}
}

func doRequest(ctx context.Context, req *http.Request) (int, error) {
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

// testAnthropic checks against GET /v1/models.
// 401 = invalid, anything else (200, 400, 403, 429) = valid key.
func testAnthropic(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodGet, "https://api.anthropic.com/v1/models", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.Header.Set("x-api-key", key)
	req.Header.Set("anthropic-version", "2023-06-01")

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}

// testOpenAI checks against GET /v1/models.
func testOpenAI(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodGet, "https://api.openai.com/v1/models", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.Header.Set("Authorization", "Bearer "+key)

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}

// testGitHub checks against GET /user.
func testGitHub(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.Header.Set("Authorization", "Bearer "+key)

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}

// testSlack checks against POST /api/auth.test.
func testSlack(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/auth.test", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.Header.Set("Authorization", "Bearer "+key)

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}

// testStripe checks against GET /v1/balance.
func testStripe(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodGet, "https://api.stripe.com/v1/balance", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.SetBasicAuth(key, "")

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}

// testSendGrid checks against GET /v3/scopes.
func testSendGrid(ctx context.Context, key string) model.SecretStatus {
	req, err := http.NewRequest(http.MethodGet, "https://api.sendgrid.com/v3/scopes", nil)
	if err != nil {
		return model.SecretStatusError
	}
	req.Header.Set("Authorization", "Bearer "+key)

	code, err := doRequest(ctx, req)
	if err != nil {
		return model.SecretStatusError
	}
	if code == http.StatusUnauthorized {
		return model.SecretStatusInvalid
	}
	return model.SecretStatusValid
}
