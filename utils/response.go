package utils

import (
	"github.com/gofiber/fiber/v2"
)

type MetaData struct {
	Status  int    `json:"status"`
	Page    int    `json:"page,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Total   int    `json:"total,omitempty"`
	Message string `json:"message"`
}

type APIResponse struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

func PaginatedResponse(c *fiber.Ctx, data interface{}, page, limit, total int, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusOK,
			Page:    page,
			Limit:   limit,
			Total:   total,
			Message: message,
		},
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func NonPaginatedResponse(c *fiber.Ctx, data interface{}, total int, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusOK,
			Total:   total,
			Message: message,
		},
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func ErrorResponse(c *fiber.Ctx, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusNotFound,
			Message: message,
		},
	}
	return c.Status(fiber.StatusNotFound).JSON(response)
}

func ValidationErrorResponse(c *fiber.Ctx, errors map[string][]string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusBadRequest,
			Message: "invalid user request",
		},
		Data: errors,
	}
	return c.Status(fiber.StatusBadRequest).JSON(response)
}

func InternalServerErrorResponse(c *fiber.Ctx, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusInternalServerError,
			Message: message,
		},
	}
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func GenericErrorResponse(c *fiber.Ctx, status int, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  status,
			Message: message,
		},
	}
	return c.Status(status).JSON(response)
}

func SuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusOK,
			Message: message,
		},
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func CreateResponse(c *fiber.Ctx, data interface{}, message string) error {
	response := APIResponse{
		Meta: MetaData{
			Status:  fiber.StatusCreated,
			Message: message,
		},
		Data: data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
