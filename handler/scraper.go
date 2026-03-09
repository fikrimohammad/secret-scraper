package handler

import (
	"github.com/gofiber/fiber/v3"
)

type ScraperREST interface {
	ScrapeSecret(c fiber.Ctx) error
}

type ScrapeSecretRequestBody struct {
	SecretProvider        string `json:"secret_provider"`
	SecretType            string `json:"secret_type"`
	MaxIterations         int    `json:"max_iterations"`
	MaxLimitPerIterations int    `json:"max_limit_per_iterations"`
}
