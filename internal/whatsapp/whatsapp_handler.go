package whatsapp

import (
	"regexp"
	"rijig/config"
	"rijig/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type QRResponse struct {
	QRCode    string `json:"qr_code,omitempty"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type StatusResponse struct {
	IsConnected bool   `json:"is_connected"`
	IsLoggedIn  bool   `json:"is_logged_in"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}

type SendMessageRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Message     string `json:"message" validate:"required"`
}

type SendMessageResponse struct {
	PhoneNumber string `json:"phone_number"`
	Timestamp   int64  `json:"timestamp"`
}

func GenerateQRHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	if wa.IsLoggedIn() {
		data := QRResponse{
			Status:    "logged_in",
			Message:   "WhatsApp is already connected and logged in",
			Timestamp: time.Now().Unix(),
		}
		return utils.SuccessWithData(c, "Already logged in", data)
	}

	qrDataURI, err := wa.GenerateQR()
	if err != nil {
		return utils.InternalServerError(c, "Failed to generate QR code: "+err.Error())
	}

	switch qrDataURI {
	case "success":
		data := QRResponse{
			Status:    "login_success",
			Message:   "WhatsApp login successful",
			Timestamp: time.Now().Unix(),
		}
		return utils.SuccessWithData(c, "Successfully logged in", data)

	case "already_connected":
		data := QRResponse{
			Status:    "already_connected",
			Message:   "WhatsApp is already connected",
			Timestamp: time.Now().Unix(),
		}
		return utils.SuccessWithData(c, "Already connected", data)

	default:

		data := QRResponse{
			QRCode:    qrDataURI,
			Status:    "qr_generated",
			Message:   "Scan QR code with WhatsApp to login",
			Timestamp: time.Now().Unix(),
		}
		return utils.SuccessWithData(c, "QR code generated successfully", data)
	}
}

func CheckLoginStatusHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	// if !wa.IsLoggedIn() {
	// 	return utils.Unauthorized(c, "WhatsApp not logged in")
	// }

	isConnected := wa.IsConnected()
	isLoggedIn := wa.IsLoggedIn()

	var status string
	var message string

	if isLoggedIn && isConnected {
		status = "connected_and_logged_in"
		message = "WhatsApp is connected and logged in"
	} else if isLoggedIn {
		status = "logged_in_but_disconnected"
		message = "WhatsApp is logged in but disconnected"
	} else if isConnected {
		status = "connected_but_not_logged_in"
		message = "WhatsApp is connected but not logged in"
	} else {
		status = "disconnected"
		message = "WhatsApp is disconnected"
	}

	data := StatusResponse{
		IsConnected: isConnected,
		IsLoggedIn:  isLoggedIn,
		Status:      status,
		Message:     message,
		Timestamp:   time.Now().Unix(),
	}

	return utils.SuccessWithData(c, "Status retrieved successfully", data)
}

func WhatsAppLogoutHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	if !wa.IsLoggedIn() {
		return utils.BadRequest(c, "No active session to logout")
	}

	err := wa.Logout()
	if err != nil {
		return utils.InternalServerError(c, "Failed to logout: "+err.Error())
	}

	data := map[string]interface{}{
		"timestamp": time.Now().Unix(),
	}

	return utils.SuccessWithData(c, "Successfully logged out and session deleted", data)
}

func SendMessageHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	if !wa.IsLoggedIn() {
		return utils.Unauthorized(c, "WhatsApp not logged in")
	}

	req := GetValidatedSendMessageRequest(c)
	if req == nil {
		return utils.BadRequest(c, "Invalid request data")
	}

	err := wa.SendMessage(req.PhoneNumber, req.Message)
	if err != nil {
		return utils.InternalServerError(c, "Failed to send message: "+err.Error())
	}

	data := SendMessageResponse{
		PhoneNumber: req.PhoneNumber,
		Timestamp:   time.Now().Unix(),
	}

	return utils.SuccessWithData(c, "Message sent successfully", data)
}

func GetDeviceInfoHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	if !wa.IsLoggedIn() {
		return utils.Unauthorized(c, "WhatsApp not logged in")
	}

	var deviceInfo map[string]interface{}
	if wa.Client != nil && wa.Client.Store.ID != nil {
		deviceInfo = map[string]interface{}{
			"device_id":    wa.Client.Store.ID.User,
			"device_name":  wa.Client.Store.ID.Device,
			"is_logged_in": wa.IsLoggedIn(),
			"is_connected": wa.IsConnected(),
			"timestamp":    time.Now().Unix(),
		}
	} else {
		deviceInfo = map[string]interface{}{
			"device_id":    nil,
			"device_name":  nil,
			"is_logged_in": false,
			"is_connected": false,
			"timestamp":    time.Now().Unix(),
		}
	}

	return utils.SuccessWithData(c, "Device info retrieved successfully", deviceInfo)
}

func HealthCheckHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return utils.InternalServerError(c, "WhatsApp service not initialized")
	}

	healthData := map[string]interface{}{
		"service_status":   "running",
		"container_status": wa.Container != nil,
		"client_status":    wa.Client != nil,
		"is_connected":     wa.IsConnected(),
		"is_logged_in":     wa.IsLoggedIn(),
		"timestamp":        time.Now().Unix(),
	}

	message := "WhatsApp service is healthy"
	if !wa.IsConnected() || !wa.IsLoggedIn() {
		message = "WhatsApp service is running but not fully operational"
	}

	return utils.SuccessWithData(c, message, healthData)
}

func validatePhoneNumber(phoneNumber string) error {

	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "+", "")

	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return fiber.NewError(fiber.StatusBadRequest, "Phone number must contain only digits")
	}

	if len(cleaned) < 10 {
		return fiber.NewError(fiber.StatusBadRequest, "Phone number too short. Include country code (e.g., 628123456789)")
	}

	if len(cleaned) > 15 {
		return fiber.NewError(fiber.StatusBadRequest, "Phone number too long")
	}

	return nil
}

func validateMessage(message string) error {

	if strings.TrimSpace(message) == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Message cannot be empty")
	}

	if len(message) > 4096 {
		return fiber.NewError(fiber.StatusBadRequest, "Message too long. Maximum 4096 characters allowed")
	}

	return nil
}

func ValidateSendMessageRequest(c *fiber.Ctx) error {
	var req SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid JSON format: "+err.Error())
	}

	if err := validatePhoneNumber(req.PhoneNumber); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	if err := validateMessage(req.Message); err != nil {
		return utils.BadRequest(c, err.Error())
	}

	req.PhoneNumber = strings.ReplaceAll(req.PhoneNumber, " ", "")
	req.PhoneNumber = strings.ReplaceAll(req.PhoneNumber, "-", "")
	req.PhoneNumber = strings.ReplaceAll(req.PhoneNumber, "+", "")

	c.Locals("validatedRequest", req)

	return c.Next()
}

func GetValidatedSendMessageRequest(c *fiber.Ctx) *SendMessageRequest {
	if req, ok := c.Locals("validatedRequest").(SendMessageRequest); ok {
		return &req
	}
	return nil
}

func ValidateContentType() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Method() == "GET" {
			return c.Next()
		}

		contentType := c.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return utils.BadRequest(c, "Content-Type must be application/json")
		}

		return c.Next()
	}
}
