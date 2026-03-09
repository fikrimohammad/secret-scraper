package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fikrimohammad/secret-scraper/handler"
	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/usecase"
	restutil "github.com/fikrimohammad/secret-scraper/util/rest"
	"github.com/gofiber/fiber/v3"
)

func (ho *handlerObject) ScrapeSecret(c fiber.Ctx) error {
	var requestBody handler.ScrapeSecretRequestBody
	if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
		return restutil.JsonApiError(c, http.StatusBadRequest, err)
	}

	if requestBody.SecretProvider == "" {
		return restutil.JsonApiError(c, http.StatusBadRequest, errors.New("secret provider is required"))
	}

	if requestBody.SecretType == "" {
		return restutil.JsonApiError(c, http.StatusBadRequest, errors.New("secret type is required"))
	}

	if requestBody.MaxLimitPerIterations == 0 {
		requestBody.MaxLimitPerIterations = 10
	}

	if requestBody.MaxIterations == 0 {
		requestBody.MaxIterations = 10
	}

	data, err := ho.scraperUseCase.ScrapeSecret(c.Context(), usecase.ScrapeSecretParams{
		SecretProvider:  model.SecretProvider(requestBody.SecretProvider),
		SecretType:      model.SecretType(requestBody.SecretType),
		MaxIterations:   requestBody.MaxIterations,
		MaxLimitPerIter: requestBody.MaxLimitPerIterations,
	})
	if err != nil {
		return restutil.JsonApiError(c, http.StatusInternalServerError, err)
	}

	return c.Status(http.StatusOK).JSON(data)
}
