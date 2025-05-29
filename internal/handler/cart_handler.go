package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartService services.CartService
}

func NewCartHandler(cartService services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddOrUpdateItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	var req dto.RequestCartItemDTO

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"request": {"Payload tidak valid"},
		})
	}

	hasErrors, _ := req.Amount > 0 && req.TrashID != "", true
	if !hasErrors {
		errs := make(map[string][]string)
		if req.Amount <= 0 {
			errs["amount"] = append(errs["amount"], "Amount harus lebih dari 0")
		}
		if req.TrashID == "" {
			errs["trash_id"] = append(errs["trash_id"], "Trash ID tidak boleh kosong")
		}
		return utils.ValidationErrorResponse(c, errs)
	}

	if err := h.cartService.AddOrUpdateItem(c.Context(), userID, req); err != nil {
		return utils.InternalServerErrorResponse(c, "Gagal menambahkan item ke keranjang")
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Item berhasil ditambahkan ke keranjang")
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	cart, err := h.cartService.GetCart(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, "Gagal mengambil data keranjang")
	}

	return utils.SuccessResponse(c, cart, "Berhasil mengambil data keranjang")
}

func (h *CartHandler) DeleteItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	trashID := c.Params("trash_id")

	if trashID == "" {
		return utils.GenericResponse(c, fiber.StatusBadRequest, "Trash ID tidak boleh kosong")
	}

	if err := h.cartService.DeleteItem(c.Context(), userID, trashID); err != nil {
		return utils.InternalServerErrorResponse(c, "Gagal menghapus item dari keranjang")
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Item berhasil dihapus dari keranjang")
}

func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.cartService.Checkout(c.Context(), userID); err != nil {
		return utils.InternalServerErrorResponse(c, "Gagal melakukan checkout keranjang")
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Checkout berhasil. Permintaan pickup telah dibuat.")
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	err := h.cartService.ClearCart(c.Context(), userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Gagal menghapus keranjang")
	}

	return utils.GenericResponse(c, fiber.StatusOK, "Keranjang berhasil dikosongkan")
}