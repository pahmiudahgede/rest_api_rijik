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
			utils.FormatDateToIndonesianFormat(requestPickup.CreatedAt),
			utils.FormatDateToIndonesianFormat(requestPickup.UpdatedAt),
		))
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Request pickup by user has been fetched",
		requestPickupResponses,
	))
}

func CreateRequestPickup(c *fiber.Ctx) error {
	var req dto.RequestPickupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.FormatResponse(
			fiber.StatusBadRequest,
			"Invalid request body",
			nil,
		))
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"User not authenticated",
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
		StatusRequest: "Pending",
	}

	service := services.NewRequestPickupService(repositories.NewRequestPickupRepository())

	if err := service.CreateRequestPickup(requestPickup); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to create request pickup",
			nil,
		))
	}

	detail, err := service.GetRequestPickupByID(requestPickup.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch created request pickup",
			nil,
		))
	}

	var requestItemsDTO []dto.RequestItemDTO
	for _, item := range detail.Request {
		requestItemsDTO = append(requestItemsDTO, dto.RequestItemDTO{
			TrashCategory:   item.TrashCategory.Name,
			EstimatedAmount: item.EstimatedAmount,
		})
	}

	userAddressDTO := dto.UserAddressDTO{
		Province:    detail.UserAddress.Province,
		District:    detail.UserAddress.District,
		Subdistrict: detail.UserAddress.Subdistrict,
		PostalCode:  detail.UserAddress.PostalCode,
		Village:     detail.UserAddress.Village,
		Detail:      detail.UserAddress.Detail,
		Geography:   detail.UserAddress.Geography,
	}

	response := dto.NewRequestPickupResponse(
		detail.ID,
		detail.UserID,
		detail.RequestTime,
		detail.StatusRequest,
		requestItemsDTO,
		userAddressDTO,
		utils.FormatDateToIndonesianFormat(detail.CreatedAt),
		utils.FormatDateToIndonesianFormat(detail.UpdatedAt),
	)

	return c.Status(fiber.StatusCreated).JSON(utils.FormatResponse(
		fiber.StatusCreated,
		"Request pickup created successfully",
		response,
	))
}
