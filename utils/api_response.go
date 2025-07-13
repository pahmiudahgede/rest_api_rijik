package utils

import (
	"github.com/gofiber/fiber/v2"
)

type Meta struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Page    *int   `json:"page,omitempty"`
	Limit   *int   `json:"limit,omitempty"`
	Total   *int   `json:"total,omitempty"`
}

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseMeta(c *fiber.Ctx, status int, message string) error {
	response := Response{
		Meta: Meta{
			Status:  status,
			Message: message,
		},
	}
	return c.Status(status).JSON(response)
}

func ResponseData(c *fiber.Ctx, status int, message string, data interface{}) error {
	response := Response{
		Meta: Meta{
			Status:  status,
			Message: message,
		},
		Data: data,
	}
	return c.Status(status).JSON(response)
}

func ResponseWithTotal(c *fiber.Ctx, status int, message string, data interface{}, total int) error {
	response := Response{
		Meta: Meta{
			Status:  status,
			Message: message,
			Total:   &total,
		},
		Data: data,
	}
	return c.Status(status).JSON(response)
}

func ResponsePagination(c *fiber.Ctx, status int, message string, data interface{}, page, limit int) error {
	response := Response{
		Meta: Meta{
			Status:  status,
			Message: message,
			Page:    &page,
			Limit:   &limit,
		},
		Data: data,
	}
	return c.Status(status).JSON(response)
}

func ResponsePaginationWithTotal(c *fiber.Ctx, status int, message string, data interface{}, page, limit, total int) error {
	response := Response{
		Meta: Meta{
			Status:  status,
			Message: message,
			Page:    &page,
			Limit:   &limit,
			Total:   &total,
		},
		Data: data,
	}
	return c.Status(status).JSON(response)
}

func ResponseErrorData(c *fiber.Ctx, status int, message string, errors interface{}) error {
	type ResponseWithErrors struct {
		Meta   Meta        `json:"meta"`
		Errors interface{} `json:"errors"`
	}
	response := ResponseWithErrors{
		Meta: Meta{
			Status:  status,
			Message: message,
		},
		Errors: errors,
	}
	return c.Status(status).JSON(response)
}

func Success(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusOK, message)
}

func SuccessWithData(c *fiber.Ctx, message string, data interface{}) error {
	return ResponseData(c, fiber.StatusOK, message, data)
}

func SuccessWithTotal(c *fiber.Ctx, message string, data interface{}, total int) error {
	return ResponseWithTotal(c, fiber.StatusOK, message, data, total)
}

func CreateSuccessWithData(c *fiber.Ctx, message string, data interface{}) error {
	return ResponseData(c, fiber.StatusCreated, message, data)
}

func SuccessWithPagination(c *fiber.Ctx, message string, data interface{}, page, limit int) error {
	return ResponsePagination(c, fiber.StatusOK, message, data, page, limit)
}

func SuccessWithPaginationAndTotal(c *fiber.Ctx, message string, data interface{}, page, limit, total int) error {
	return ResponsePaginationWithTotal(c, fiber.StatusOK, message, data, page, limit, total)
}

func BadRequest(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusBadRequest, message)
}

func NotFound(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusNotFound, message)
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusInternalServerError, message)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusUnauthorized, message)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return ResponseMeta(c, fiber.StatusForbidden, message)
}
