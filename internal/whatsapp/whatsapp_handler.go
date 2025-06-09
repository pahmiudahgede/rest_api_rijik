package whatsapp

import (
	"html/template"
	"path/filepath"
	"rijig/config"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Meta map[string]interface{} `json:"meta"`
	Data interface{}            `json:"data,omitempty"`
}

func WhatsAppQRPageHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": "WhatsApp service not initialized",
			},
		})
	}

	// Jika sudah login, tampilkan halaman success
	if wa.IsLoggedIn() {
		templatePath := filepath.Join("internal", "whatsapp", "success_scan.html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
				Meta: map[string]interface{}{
					"status":  "error",
					"message": "Unable to load success template: " + err.Error(),
				},
			})
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
	}

	qrDataURI, err := wa.GenerateQR()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": "Failed to generate QR code: " + err.Error(),
			},
		})
	}

	if qrDataURI == "success" {
		// Login berhasil, tampilkan halaman success
		templatePath := filepath.Join("internal", "whatsapp", "success_scan.html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
				Meta: map[string]interface{}{
					"status":  "error",
					"message": "Unable to load success template: " + err.Error(),
				},
			})
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
	}

	if qrDataURI == "already_connected" {
		// Sudah terhubung, tampilkan halaman success
		templatePath := filepath.Join("internal", "whatsapp", "success_scan.html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
				Meta: map[string]interface{}{
					"status":  "error",
					"message": "Unable to load success template: " + err.Error(),
				},
			})
		}

		c.Set("Content-Type", "text/html")
		return tmpl.Execute(c.Response().BodyWriter(), nil)
	}

	// Tampilkan QR code scanner
	templatePath := filepath.Join("internal", "whatsapp", "scanner.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": "Unable to load scanner template: " + err.Error(),
			},
		})
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c.Response().BodyWriter(), template.URL(qrDataURI))
}

func WhatsAppLogoutHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": "WhatsApp service not initialized",
			},
		})
	}

	err := wa.Logout()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Meta: map[string]interface{}{
			"status":  "success",
			"message": "Successfully logged out and session deleted",
		},
	})
}

func WhatsAppStatusHandler(c *fiber.Ctx) error {
	wa := config.GetWhatsAppService()
	if wa == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Meta: map[string]interface{}{
				"status":  "error",
				"message": "WhatsApp service not initialized",
			},
		})
	}

	status := map[string]interface{}{
		"is_connected": wa.IsConnected(),
		"is_logged_in": wa.IsLoggedIn(),
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Meta: map[string]interface{}{
			"status":  "success",
			"message": "WhatsApp status retrieved successfully",
		},
		Data: status,
	})
}
