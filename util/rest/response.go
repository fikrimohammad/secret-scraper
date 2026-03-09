package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type BaseResponse struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code   string `json:"code"`
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func JsonApiError(c fiber.Ctx, status int, err error) error {
	responseBody := BaseResponse{
		Error: &ErrorResponse{
			Status: http.StatusText(status),
			Code:   "INTERNAL_ERROR",
			Title:  err.Error(),
			Detail: err.Error(),
		},
	}

	return c.Status(status).JSON(responseBody)
}

func JsonApiSuccess(c fiber.Ctx, status int, data interface{}) error {
	responseBody := BaseResponse{
		Data: data,
	}

	return c.Status(status).JSON(responseBody)
}
