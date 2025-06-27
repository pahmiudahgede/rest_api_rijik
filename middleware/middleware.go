package middleware

import (
	"crypto/subtle"
	"fmt"
	"rijig/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type AuthConfig struct {
	RequiredTokenType  utils.TokenType
	RequiredRoles      []string
	RequiredStatuses   []string
	RequiredStep       int
	RequireComplete    bool
	SkipAuth           bool
	AllowPartialToken  bool
	CustomErrorHandler ErrorHandler
}

type ErrorHandler func(c *fiber.Ctx, err error) error

type AuthContext struct {
	Claims    *utils.JWTClaims
	StepInfo  *utils.RegistrationStepInfo
	IsAdmin   bool
	CanAccess bool
}

type AuthError struct {
	Code    string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AuthError) Error() string {
	return e.Message
}

var (
	ErrMissingToken = &AuthError{
		Code:    "MISSING_TOKEN",
		Message: "Token akses diperlukan",
	}

	ErrInvalidTokenFormat = &AuthError{
		Code:    "INVALID_TOKEN_FORMAT",
		Message: "Format token tidak valid",
	}

	ErrInvalidToken = &AuthError{
		Code:    "INVALID_TOKEN",
		Message: "Token tidak valid atau telah kadaluarsa",
	}

	ErrUserContextNotFound = &AuthError{
		Code:    "USER_CONTEXT_NOT_FOUND",
		Message: "Silakan login terlebih dahulu",
	}

	ErrInsufficientPermissions = &AuthError{
		Code:    "INSUFFICIENT_PERMISSIONS",
		Message: "Akses ditolak untuk role ini",
	}

	ErrRegistrationIncomplete = &AuthError{
		Code:    "REGISTRATION_INCOMPLETE",
		Message: "Registrasi belum lengkap",
	}

	ErrRegistrationNotApproved = &AuthError{
		Code:    "REGISTRATION_NOT_APPROVED",
		Message: "Registrasi belum disetujui",
	}

	ErrInvalidTokenType = &AuthError{
		Code:    "INVALID_TOKEN_TYPE",
		Message: "Tipe token tidak sesuai",
	}

	ErrStepNotAccessible = &AuthError{
		Code:    "STEP_NOT_ACCESSIBLE",
		Message: "Step registrasi belum dapat diakses",
	}

	ErrAwaitingApproval = &AuthError{
		Code:    "AWAITING_ADMIN_APPROVAL",
		Message: "Menunggu persetujuan admin",
	}

	ErrInvalidRegistrationStatus = &AuthError{
		Code:    "INVALID_REGISTRATION_STATUS",
		Message: "Status registrasi tidak sesuai",
	}
)

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	if authErr, ok := err.(*AuthError); ok {
		statusCode := getStatusCodeForError(authErr.Code)
		return c.Status(statusCode).JSON(authErr)
	}

	return c.Status(fiber.StatusInternalServerError).JSON(&AuthError{
		Code:    "INTERNAL_ERROR",
		Message: "Terjadi kesalahan internal",
	})
}

func getStatusCodeForError(errorCode string) int {
	switch errorCode {
	case "MISSING_TOKEN", "INVALID_TOKEN_FORMAT", "INVALID_TOKEN", "USER_CONTEXT_NOT_FOUND":
		return fiber.StatusUnauthorized
	case "INSUFFICIENT_PERMISSIONS", "REGISTRATION_INCOMPLETE", "REGISTRATION_NOT_APPROVED",
		"INVALID_TOKEN_TYPE", "STEP_NOT_ACCESSIBLE", "AWAITING_ADMIN_APPROVAL",
		"INVALID_REGISTRATION_STATUS":
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}

func AuthMiddleware(config ...AuthConfig) fiber.Handler {
	cfg := AuthConfig{}
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.CustomErrorHandler == nil {
		cfg.CustomErrorHandler = defaultErrorHandler
	}

	return func(c *fiber.Ctx) error {

		if cfg.SkipAuth {
			return c.Next()
		}

		claims, err := extractAndValidateToken(c)
		if err != nil {
			return cfg.CustomErrorHandler(c, err)
		}

		authCtx := createAuthContext(claims)

		if err := validateAuthConfig(authCtx, cfg); err != nil {
			return cfg.CustomErrorHandler(c, err)
		}

		c.Locals("user", claims)
		c.Locals("auth_context", authCtx)

		return c.Next()
	}
}

func extractAndValidateToken(c *fiber.Ctx) (*utils.JWTClaims, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, ErrMissingToken
	}

	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return nil, ErrInvalidTokenFormat
	}

	claims, err := utils.ValidateAccessToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func createAuthContext(claims *utils.JWTClaims) *AuthContext {
	stepInfo := utils.GetRegistrationStepInfo(
		claims.Role,
		claims.RegistrationProgress,
		claims.RegistrationStatus,
	)

	return &AuthContext{
		Claims:    claims,
		StepInfo:  stepInfo,
		IsAdmin:   claims.Role == utils.RoleAdministrator,
		CanAccess: stepInfo.IsAccessible,
	}
}

func validateAuthConfig(authCtx *AuthContext, cfg AuthConfig) error {
	claims := authCtx.Claims

	if cfg.RequiredTokenType != "" {
		if claims.TokenType != cfg.RequiredTokenType {
			return &AuthError{
				Code:    "INVALID_TOKEN_TYPE",
				Message: fmt.Sprintf("Endpoint memerlukan token type: %s", cfg.RequiredTokenType),
				Details: fiber.Map{
					"current_token_type":  claims.TokenType,
					"required_token_type": cfg.RequiredTokenType,
				},
			}
		}
	}

	if len(cfg.RequiredRoles) > 0 {
		if !contains(cfg.RequiredRoles, claims.Role) {
			return &AuthError{
				Code:    "INSUFFICIENT_PERMISSIONS",
				Message: "Akses ditolak untuk role ini",
				Details: fiber.Map{
					"user_role":     claims.Role,
					"allowed_roles": cfg.RequiredRoles,
				},
			}
		}
	}

	if len(cfg.RequiredStatuses) > 0 {
		if !contains(cfg.RequiredStatuses, claims.RegistrationStatus) {
			return &AuthError{
				Code:    "INVALID_REGISTRATION_STATUS",
				Message: "Status registrasi tidak sesuai",
				Details: fiber.Map{
					"current_status":   claims.RegistrationStatus,
					"allowed_statuses": cfg.RequiredStatuses,
					"next_step":        authCtx.StepInfo.Description,
				},
			}
		}
	}

	if cfg.RequiredStep > 0 {
		if claims.RegistrationProgress < cfg.RequiredStep {
			return &AuthError{
				Code:    "STEP_NOT_ACCESSIBLE",
				Message: "Step registrasi belum dapat diakses",
				Details: fiber.Map{
					"current_step":      claims.RegistrationProgress,
					"required_step":     cfg.RequiredStep,
					"current_step_info": authCtx.StepInfo.Description,
				},
			}
		}

		if authCtx.StepInfo.RequiresAdminApproval && !authCtx.CanAccess {
			return &AuthError{
				Code:    "AWAITING_ADMIN_APPROVAL",
				Message: "Menunggu persetujuan admin",
				Details: fiber.Map{
					"status": claims.RegistrationStatus,
				},
			}
		}
	}

	if cfg.RequireComplete {
		if claims.TokenType != utils.TokenTypeFull {
			return &AuthError{
				Code:    "REGISTRATION_INCOMPLETE",
				Message: "Registrasi belum lengkap",
				Details: fiber.Map{
					"registration_status":     claims.RegistrationStatus,
					"registration_progress":   claims.RegistrationProgress,
					"next_step":               authCtx.StepInfo.Description,
					"requires_admin_approval": authCtx.StepInfo.RequiresAdminApproval,
					"can_proceed":             authCtx.CanAccess,
				},
			}
		}

		if !utils.IsRegistrationComplete(claims.RegistrationStatus) {
			return ErrRegistrationNotApproved
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RequireAuth() fiber.Handler {
	return AuthMiddleware()
}

func RequireFullToken() fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredTokenType: utils.TokenTypeFull,
		RequireComplete:   true,
	})
}

func RequirePartialToken() fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredTokenType: utils.TokenTypePartial,
	})
}

func RequireRoles(roles ...string) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredRoles: roles,
	})
}

func RequireAdminRole() fiber.Handler {
	return RequireRoles(utils.RoleAdministrator)
}

func RequireRegistrationStep(step int) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredStep: step,
	})
}

func RequireRegistrationStatus(statuses ...string) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredStatuses: statuses,
	})
}

func RequireTokenType(tokenType utils.TokenType) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredTokenType: tokenType,
	})
}

func RequireCompleteRegistrationForRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		if contains(roles, claims.Role) {
			return RequireFullToken()(c)
		}

		return c.Next()
	}
}

func RequireRoleAndComplete(roles ...string) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredRoles:   roles,
		RequireComplete: true,
	})
}

func RequireRoleAndStep(step int, roles ...string) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredRoles: roles,
		RequiredStep:  step,
	})
}

func RequireRoleAndStatus(statuses []string, roles ...string) fiber.Handler {
	return AuthMiddleware(AuthConfig{
		RequiredRoles:    roles,
		RequiredStatuses: statuses,
	})
}

func GetUserFromContext(c *fiber.Ctx) (*utils.JWTClaims, error) {
	claims, ok := c.Locals("user").(*utils.JWTClaims)
	if !ok {
		return nil, ErrUserContextNotFound
	}
	return claims, nil
}

func GetAuthContextFromContext(c *fiber.Ctx) (*AuthContext, error) {
	authCtx, ok := c.Locals("auth_context").(*AuthContext)
	if !ok {

		claims, err := GetUserFromContext(c)
		if err != nil {
			return nil, err
		}
		return createAuthContext(claims), nil
	}
	return authCtx, nil
}

func MustGetUserFromContext(c *fiber.Ctx) *utils.JWTClaims {
	claims, err := GetUserFromContext(c)
	if err != nil {
		panic("user context not found")
	}
	return claims
}

func GetUserID(c *fiber.Ctx) (string, error) {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

func GetUserRole(c *fiber.Ctx) (string, error) {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}

func IsAdmin(c *fiber.Ctx) bool {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return false
	}
	return claims.Role == utils.RoleAdministrator
}

func IsRegistrationComplete(c *fiber.Ctx) bool {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return false
	}
	return utils.IsRegistrationComplete(claims.RegistrationStatus)
}

func HasRole(c *fiber.Ctx, role string) bool {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return false
	}
	return claims.Role == role
}

func HasAnyRole(c *fiber.Ctx, roles ...string) bool {
	claims, err := GetUserFromContext(c)
	if err != nil {
		return false
	}
	return contains(roles, claims.Role)
}

type RateLimitConfig struct {
	MaxRequests int
	Window      time.Duration
	KeyFunc     func(*fiber.Ctx) string
	SkipFunc    func(*fiber.Ctx) bool
}

func AuthRateLimit(config RateLimitConfig) fiber.Handler {
	if config.KeyFunc == nil {
		config.KeyFunc = func(c *fiber.Ctx) string {
			claims, err := GetUserFromContext(c)
			if err != nil {
				return c.IP()
			}
			return fmt.Sprintf("user:%s", claims.UserID)
		}
	}

	return func(c *fiber.Ctx) error {
		if config.SkipFunc != nil && config.SkipFunc(c) {
			return c.Next()
		}

		key := fmt.Sprintf("rate_limit:%s", config.KeyFunc(c))

		var count int
		err := utils.GetCache(key, &count)
		if err != nil {
			count = 0
		}

		if count >= config.MaxRequests {
			return c.Status(fiber.StatusTooManyRequests).JSON(&AuthError{
				Code:    "RATE_LIMIT_EXCEEDED",
				Message: "Terlalu banyak permintaan, coba lagi nanti",
				Details: fiber.Map{
					"max_requests": config.MaxRequests,
					"window":       config.Window.String(),
				},
			})
		}

		count++
		utils.SetCache(key, count, config.Window)

		return c.Next()
	}
}

func DeviceValidation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		deviceID := claims.DeviceID
		if deviceID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(&AuthError{
				Code:    "MISSING_DEVICE_ID",
				Message: "Device ID diperlukan",
			})
		}

		if subtle.ConstantTimeCompare([]byte(claims.DeviceID), []byte(deviceID)) != 1 {
			return c.Status(fiber.StatusForbidden).JSON(&AuthError{
				Code:    "DEVICE_MISMATCH",
				Message: "Device tidak cocok dengan token",
			})
		}

		return c.Next()
	}
}

func SessionValidation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, err := GetUserFromContext(c)
		if err != nil {
			return err
		}

		sessionKey := fmt.Sprintf("session:%s", claims.SessionID)
		var sessionData interface{}
		err = utils.GetCache(sessionKey, &sessionData)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(&AuthError{
				Code:    "SESSION_EXPIRED",
				Message: "Sesi telah berakhir, silakan login kembali",
			})
		}

		return c.Next()
	}
}

func AuthLogger() fiber.Handler {
	return logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - User: ${locals:user_id} Role: ${locals:user_role} IP: ${ip}\n",
		CustomTags: map[string]logger.LogFunc{
			"user_id": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				if claims, err := GetUserFromContext(c); err == nil {
					return output.WriteString(claims.UserID)
				}
				return output.WriteString("anonymous")
			},
			"user_role": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
				if claims, err := GetUserFromContext(c); err == nil {
					return output.WriteString(claims.Role)
				}
				return output.WriteString("none")
			},
		},
	})
}
