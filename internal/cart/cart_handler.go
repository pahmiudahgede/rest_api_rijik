package cart

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartService CartService
}

func NewCartHandler(cartService CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddOrUpdateItem(c *fiber.Ctx) error {

	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	var req RequestCartItemDTO

	if err := c.BodyParser(&req); err != nil {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Payload tidak valid", map[string][]string{
			"request": {"Payload tidak valid"},
		})
	}

	hasErrors := req.Amount <= 0 || req.TrashID == ""
	if hasErrors {
		errs := make(map[string][]string)
		if req.Amount <= 0 {
			errs["amount"] = append(errs["amount"], "Amount harus lebih dari 0")
		}
		if req.TrashID == "" {
			errs["trash_id"] = append(errs["trash_id"], "Trash ID tidak boleh kosong")
		}
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validasi gagal", errs)
	}

	if err := h.cartService.AddOrUpdateItem(c.Context(), claims.UserID, req); err != nil {
		return utils.InternalServerError(c, "Gagal menambahkan item ke keranjang")
	}

	return utils.Success(c, "Item berhasil ditambahkan ke keranjang")
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	cart, err := h.cartService.GetCart(c.Context(), claims.UserID)
	if err != nil {
		return utils.InternalServerError(c, "Gagal mengambil data keranjang")
	}

	return utils.SuccessWithData(c, "Berhasil mengambil data keranjang", cart)
}

func (h *CartHandler) DeleteItem(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}
	trashID := c.Params("trash_id")

	if trashID == "" {
		return utils.BadRequest(c, "Trash ID tidak boleh kosong")
	}

	if err := h.cartService.DeleteItem(c.Context(), claims.UserID, trashID); err != nil {
		return utils.InternalServerError(c, "Gagal menghapus item dari keranjang")
	}

	return utils.Success(c, "Item berhasil dihapus dari keranjang")
}

func (h *CartHandler) Checkout(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	if err := h.cartService.Checkout(c.Context(), claims.UserID); err != nil {
		return utils.InternalServerError(c, "Gagal melakukan checkout keranjang")
	}

	return utils.Success(c, "Checkout berhasil. Permintaan pickup telah dibuat.")
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	err = h.cartService.ClearCart(c.Context(), claims.UserID)
	if err != nil {
		return utils.InternalServerError(c, "Gagal menghapus keranjang")
	}

	return utils.Success(c, "Keranjang berhasil dikosongkan")
}
