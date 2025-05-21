package handler

import (
	"context"
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler interface {
	GetCart(c *fiber.Ctx) error
	AddOrUpdateCartItem(c *fiber.Ctx) error
	AddMultipleCartItems(c *fiber.Ctx) error
	DeleteCartItem(c *fiber.Ctx) error
	ClearCart(c *fiber.Ctx) error
}

type cartHandler struct {
	service services.CartService
}

func NewCartHandler(service services.CartService) CartHandler {
	return &cartHandler{service: service}
}

// GET /cart
func (h *cartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	cart, err := h.service.GetCart(context.Background(), userID)
	if err != nil {
		return utils.ErrorResponse(c, "Cart belum dibuat atau sudah kadaluarsa")
	}

	return utils.SuccessResponse(c, cart, "Data cart berhasil diambil")
}

// POST /cart/item
func (h *cartHandler) AddOrUpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var item dto.RequestCartItemDTO
	if err := c.BodyParser(&item); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{"body": {"format tidak valid"}})
	}

	if item.TrashID == "" || item.Amount <= 0 {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"trash_id": {"harus diisi"},
			"amount":   {"harus lebih dari 0"},
		})
	}

	if err := h.service.AddOrUpdateItem(context.Background(), userID, item); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Item berhasil ditambahkan/diupdate di cart")
}

// POST /cart/items
func (h *cartHandler) AddMultipleCartItems(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var payload dto.RequestCartDTO
	if err := c.BodyParser(&payload); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"format tidak valid"},
		})
	}

	if errs, ok := payload.ValidateRequestCartDTO(); !ok {
		return utils.ValidationErrorResponse(c, errs)
	}

	for _, item := range payload.CartItems {
		if err := h.service.AddOrUpdateItem(context.Background(), userID, item); err != nil {
			return utils.InternalServerErrorResponse(c, err.Error())
		}
	}

	return utils.SuccessResponse(c, nil, "Semua item berhasil ditambahkan/diupdate ke cart")
}


// DELETE /cart/item/:trashID
func (h *cartHandler) DeleteCartItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	trashID := c.Params("trashID")

	if trashID == "" {
		return utils.ValidationErrorResponse(c, map[string][]string{"trash_id": {"tidak boleh kosong"}})
	}

	err := h.service.DeleteItem(context.Background(), userID, trashID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Item berhasil dihapus dari cart")
}

// DELETE /cart
func (h *cartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.service.ClearCart(context.Background(), userID); err != nil {
		return utils.InternalServerErrorResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, nil, "Seluruh cart berhasil dihapus")
}
