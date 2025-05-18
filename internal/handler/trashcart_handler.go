package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	Service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{Service: service}
}

func (h *CartHandler) AddOrUpdateCartItem(c *fiber.Ctx) error {
	var body dto.BulkRequestCartItems
	if err := c.BodyParser(&body); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"Invalid JSON body"},
		})
	}

	if errors, ok := body.Validate(); !ok {
		return utils.ValidationErrorResponse(c, errors)
	}

	userID := c.Locals("userID").(string)
	for _, item := range body.Items {
		if err := services.AddOrUpdateCartItem(userID, item); err != nil {
			return utils.InternalServerErrorResponse(c, "Failed to update one or more items")
		}
	}

	return utils.SuccessResponse(c, nil, "Cart updated successfully")
}

func (h *CartHandler) DeleteCartItem(c *fiber.Ctx) error {
	trashID := c.Params("trashid")
	userID := c.Locals("userID").(string)

	err := services.DeleteCartItem(userID, trashID)
	if err != nil {
		if err.Error() == "no cart found" || err.Error() == "trashid not found" {
			return utils.GenericResponse(c, fiber.StatusNotFound, "Trash item not found in cart")
		}
		return utils.InternalServerErrorResponse(c, "Failed to delete item")
	}

	return utils.SuccessResponse(c, nil, "Item deleted")
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	if err := services.ClearCart(userID); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to clear cart")
	}
	return utils.SuccessResponse(c, nil, "Cart cleared")
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	cart, err := h.Service.GetCart(userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch cart")
	}

	return utils.SuccessResponse(c, cart, "User cart data successfully fetched")
}

func (h *CartHandler) CommitCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	err := h.Service.CommitCartToDatabase(userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to commit cart to database")
	}

	return utils.SuccessResponse(c, nil, "Cart committed to database")
}

// PUT /cart/refresh → refresh TTL Redis
func (h *CartHandler) RefreshCartTTL(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	err := services.RefreshCartTTL(userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to refresh cart TTL")
	}

	return utils.SuccessResponse(c, nil, "Cart TTL refreshed")
}
