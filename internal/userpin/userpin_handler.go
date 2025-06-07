package userpin

import (
	"rijig/middleware"
	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

type UserPinHandler struct {
	service UserPinService
}

func NewUserPinHandler(service UserPinService) *UserPinHandler {
	return &UserPinHandler{service}
}

// userID, ok := c.Locals("user_id").(string)
//
//	if !ok || userID == "" {
//		return utils.Unauthorized(c, "user_id is missing or invalid")
//	}
func (h *UserPinHandler) CreateUserPinHandler(c *fiber.Ctx) error {
	// Ambil klaim pengguna yang sudah diautentikasi
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	// Parsing body request untuk PIN
	var req RequestPinDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	// Validasi request PIN
	if errs, ok := req.ValidateRequestPinDTO(); !ok {
		return utils.ResponseErrorData(c, fiber.StatusBadRequest, "Validation error", errs)
	}

	// Panggil service untuk membuat PIN
	err = h.service.CreateUserPin(c.Context(), claims.UserID, &req)
	if err != nil {
		if err.Error() == "PIN already created" {
			return utils.BadRequest(c, err.Error()) // Jika PIN sudah ada, kembalikan error 400
		}
		return utils.InternalServerError(c, err.Error()) // Jika terjadi error lain, internal server error
	}

	// Mengembalikan response sukses jika berhasil
	return utils.Success(c, "PIN created successfully")
}

func (h *UserPinHandler) VerifyPinHandler(c *fiber.Ctx) error {
	// userID, ok := c.Locals("user_id").(string)
	// if !ok || userID == "" {
	// 	return utils.Unauthorized(c, "user_id is missing or invalid")
	// }
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return err
	}

	var req RequestPinDTO
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}

	token, err := h.service.VerifyUserPin(c.Context(), claims.UserID, &req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}

	return utils.SuccessWithData(c, "PIN verified successfully", fiber.Map{
		"token": token,
	})
}
