package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func GetRequestPickupsByUser(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.FormatResponse(
			fiber.StatusUnauthorized,
			"User not authenticated",
			nil,
		))
	}

	requestPickups, err := services.GetRequestPickupsByUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.FormatResponse(
			fiber.StatusInternalServerError,
			"Failed to fetch request pickups",
			nil,
		))
	}

	var requestPickupResponses []dto.RequestPickupResponse
	for _, requestPickup := range requestPickups {
		var requestItems []dto.RequestItemDTO
		for _, item := range requestPickup.Request {
			requestItems = append(requestItems, dto.RequestItemDTO{
				TrashCategory:   item.TrashCategory.Name,
				EstimatedAmount: item.EstimatedAmount,
			})
		}

		userAddress := dto.UserAddressDTO{
			Province:    requestPickup.UserAddress.Province,
			District:    requestPickup.UserAddress.District,
			Subdistrict: requestPickup.UserAddress.Subdistrict,
			PostalCode:  requestPickup.UserAddress.PostalCode,
			Village:     requestPickup.UserAddress.Village,
			Detail:      requestPickup.UserAddress.Detail,
			Geography:   requestPickup.UserAddress.Geography,
		}

		requestPickupResponse := dto.NewRequestPickupResponse(
			requestPickup.ID,
			requestPickup.UserID,
			requestPickup.RequestTime,
			requestPickup.StatusRequest,
			requestItems,
			userAddress,
		)

		requestPickupResponses = append(requestPickupResponses, requestPickupResponse)
	}

	return c.Status(fiber.StatusOK).JSON(utils.FormatResponse(
		fiber.StatusOK,
		"Request pickup by user has been fetched",
		requestPickupResponses,
	))
}
