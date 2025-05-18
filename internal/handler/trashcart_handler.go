package handler

import (
	"rijig/dto"
	"rijig/internal/services"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	CartService services.CartService
}

func NewCartHandler(service services.CartService) *CartHandler {
	return &CartHandler{
		CartService: service,
	}
}

// GET /cart - Get cart by user ID
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "unauthorized or invalid user")
	}

	cart, err := h.CartService.GetCartByUserID(userID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "failed to retrieve cart")
	}

	if cart == nil {
		return utils.SuccessResponse(c, nil, "Cart is empty")
	}

	return utils.SuccessResponse(c, cart, "User cart data successfully fetched")
}

// POST /cart - Create new cart
func (h *CartHandler) CreateCart(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "unauthorized or invalid user")
	}

	var reqItems []dto.RequestCartItems
	if err := c.BodyParser(&reqItems); err != nil {
		return utils.ValidationErrorResponse(c, map[string][]string{
			"body": {"invalid JSON format"},
		})
	}

	// Logic dipindahkan ke service
	if err := h.CartService.CreateCartFromDTO(userID, reqItems); err != nil {
		if ve, ok := err.(dto.ValidationErrors); ok {
			return utils.ValidationErrorResponse(c, ve.Errors)
		}
		return utils.InternalServerErrorResponse(c, "failed to create cart")
	}

	return utils.CreateResponse(c, nil, "Cart created successfully")
}


// DELETE /cart/:id - Delete cart by cartID
func (h *CartHandler) DeleteCart(c *fiber.Ctx) error {
	cartID := c.Params("id")
	if cartID == "" {
		return utils.ErrorResponse(c, "Cart ID is required")
	}

	if err := h.CartService.DeleteCart(cartID); err != nil {
		return utils.InternalServerErrorResponse(c, "failed to delete cart")
	}

	return utils.SuccessResponse(c, nil, "Cart deleted successfully")
}

// POST /cart/commit - Simpan cart dari Redis ke DB
func (h *CartHandler) CommitCart(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return utils.ErrorResponse(c, "unauthorized or invalid user")
	}

	err := h.CartService.CommitCartFromRedis(userID)
	if err != nil {
		if err.Error() == "cart not found in redis" {
			return utils.ErrorResponse(c, "Cart tidak ditemukan atau sudah expired")
		}
		return utils.InternalServerErrorResponse(c, "Gagal menyimpan cart ke database")
	}

	return utils.SuccessResponse(c, nil, "Cart berhasil disimpan ke database")
}
