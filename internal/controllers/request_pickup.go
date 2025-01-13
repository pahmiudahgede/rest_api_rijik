package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetRequestPickupsByUser(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"User not authenticated",
			nil,
		))
	}

	service := services.NewRequestPickupService(repositories.NewRequestPickupRepository())

	requestPickups, err := service.GetRequestPickupsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch request pickups",
			nil,
		))
	}

	var requestPickupResponses []dto.RequestPickupResponse
	for _, requestPickup := range requestPickups {
		userAddress := dto.UserAddressDTO{
			Province:    requestPickup.UserAddress.Province,
			District:    requestPickup.UserAddress.District,
			Subdistrict: requestPickup.UserAddress.Subdistrict,
			PostalCode:  requestPickup.UserAddress.PostalCode,
			Village:     requestPickup.UserAddress.Village,
			Detail:      requestPickup.UserAddress.Detail,
			Geography:   requestPickup.UserAddress.Geography,
		}

		var requestItems []dto.RequestItemDTO
		for _, item := range requestPickup.Request {
			requestItems = append(requestItems, dto.RequestItemDTO{
				TrashCategory:   item.TrashCategory.Name,
				EstimatedAmount: item.EstimatedAmount,
			})
		}

		requestPickupResponses = append(requestPickupResponses, dto.NewRequestPickupResponse(
			requestPickup.ID,
			requestPickup.UserID,
			requestPickup.RequestTime,
			requestPickup.StatusRequest,
			requestItems,
			userAddress,
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Request pickup by user has been fetched",
		requestPickupResponses,
	))
}

func CreateRequestPickup(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"User not authenticated",
			nil,
		))
	}

	var req dto.RequestPickupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			nil,
		))
	}

	if req.UserAddressID == "" || len(req.RequestItems) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Missing required fields",
			nil,
		))
	}

	var requestItems []domain.RequestItem
	for _, item := range req.RequestItems {
		requestItems = append(requestItems, domain.RequestItem{
			TrashCategoryID: item.TrashCategory,
			EstimatedAmount: item.EstimatedAmount,
		})
	}

	requestPickup := &domain.RequestPickup{
		UserID:        userID,
		Request:       requestItems,
		RequestTime:   req.RequestTime,
		UserAddressID: req.UserAddressID,
		StatusRequest: "waiting driver",
	}

	service := services.NewRequestPickupService(repositories.NewRequestPickupRepository())
	if err := service.CreateRequestPickup(requestPickup); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create request pickup",
			nil,
		))
	}

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Request pickup created successfully",
		nil,
	))
}
