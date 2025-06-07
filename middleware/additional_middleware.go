package middleware
/* 
import (
	"fmt"
	"time"

	"rijig/utils"

	"github.com/gofiber/fiber/v2"
)

func RateLimitByUser(maxRequests int, duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		key := fmt.Sprintf("rate_limit:%s:%s", claims.UserID, c.Route().Path)

		count, err := utils.IncrementCounter(key, duration)
		if err != nil {

			return c.Next()
		}

		if count > int64(maxRequests) {

			ttl, _ := utils.GetTTL(key)
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Rate limit exceeded",
				"message":     "Terlalu banyak permintaan, silakan coba lagi nanti",
				"retry_after": int64(ttl.Seconds()),
				"limit":       maxRequests,
				"remaining":   0,
			})
		}

		c.Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", maxRequests-int(count)))
		c.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(duration).Unix()))

		return c.Next()
	}
}

func SessionValidation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		if claims.SessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid session",
				"message": "Session tidak valid",
			})
		}

		sessionKey := fmt.Sprintf("session:%s", claims.SessionID)
		var sessionData map[string]interface{}
		err = utils.GetCache(sessionKey, &sessionData)
		if err != nil {
			if err.Error() == "ErrCacheMiss" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Session not found",
					"message": "Session tidak ditemukan, silakan login kembali",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Session error",
				"message": "Terjadi kesalahan saat validasi session",
			})
		}

		if sessionData["user_id"] != claims.UserID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Session mismatch",
				"message": "Session tidak sesuai dengan user",
			})
		}

		if expiryInterface, exists := sessionData["expires_at"]; exists {
			if expiry, ok := expiryInterface.(float64); ok {
				if time.Now().Unix() > int64(expiry) {

					utils.DeleteCache(sessionKey)
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error":   "Session expired",
						"message": "Session telah berakhir, silakan login kembali",
					})
				}
			}
		}

		return c.Next()
	}
}

func RequireApprovedRegistration() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		if claims.RegistrationStatus == utils.RegStatusRejected {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Registration rejected",
				"message": "Registrasi Anda ditolak, silakan hubungi admin",
			})
		}

		if claims.RegistrationStatus == utils.RegStatusPending {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Registration pending",
				"message": "Registrasi Anda masih menunggu persetujuan admin",
			})
		}

		if claims.RegistrationStatus != utils.RegStatusComplete {
			progress := utils.GetUserRegistrationProgress(claims.UserID)
			nextStep := utils.GetNextRegistrationStep(claims.Role, progress)

			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":               "Registration incomplete",
				"message":             "Silakan lengkapi registrasi terlebih dahulu",
				"registration_status": claims.RegistrationStatus,
				"next_step":           nextStep,
			})
		}

		return c.Next()
	}
}

func ConditionalAuth(condition func(*utils.JWTClaims) bool, errorMessage string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		if !condition(claims) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Condition not met",
				"message": errorMessage,
			})
		}

		return c.Next()
	}
}

func RequireSpecificRole(role string) fiber.Handler {
	return ConditionalAuth(
		func(claims *utils.JWTClaims) bool {
			return claims.Role == role
		},
		fmt.Sprintf("Akses ini hanya untuk role %s", role),
	)
}

func RequireCompleteRegistrationAndSpecificRole(role string) fiber.Handler {
	return ConditionalAuth(
		func(claims *utils.JWTClaims) bool {
			return claims.Role == role && utils.IsRegistrationComplete(claims.RegistrationStatus)
		},
		fmt.Sprintf("Akses ini hanya untuk role %s dengan registrasi lengkap", role),
	)
}

func DeviceValidation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		deviceID := c.Get("X-Device-ID")
		if deviceID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Device ID required",
				"message": "Device ID diperlukan",
			})
		}

		if claims.DeviceID != deviceID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Device mismatch",
				"message": "Token tidak valid untuk device ini",
			})
		}

		return c.Next()
	}
}
 */