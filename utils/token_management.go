package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"rijig/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	RoleAdministrator = "administrator"
	RolePengelola     = "pengelola"
	RolePengepul      = "pengepul"
	RoleMasyarakat    = "masyarakat"
)

const (
	TokenTypePartial TokenType = "partial"
	TokenTypeFull    TokenType = "full"
	TokenTypeRefresh TokenType = "refresh"
)

const (
	RegStatusIncomplete = "uncomplete"
	RegStatusPending    = "awaiting_approval"
	RegStatusConfirmed  = "approved"
	RegStatusComplete   = "complete"
	RegStatusRejected   = "rejected"
)

const (
	ProgressOTPVerified   = 1
	ProgressDataSubmitted = 2
	ProgressComplete      = 3
)

type JWTClaims struct {
	UserID               string    `json:"user_id"`
	Role                 string    `json:"role"`
	DeviceID             string    `json:"device_id"`
	RegistrationStatus   string    `json:"registration_status"`
	RegistrationProgress int       `json:"registration_progress"`
	TokenType            TokenType `json:"token_type"`
	SessionID            string    `json:"session_id,omitempty"`
	jwt.RegisteredClaims
}

type RefreshTokenData struct {
	RefreshToken         string `json:"refresh_token"`
	ExpiresAt            int64  `json:"expires_at"`
	DeviceID             string `json:"device_id"`
	Role                 string `json:"role"`
	RegistrationStatus   string `json:"registration_status"`
	RegistrationProgress int    `json:"registration_progress"`
	SessionID            string `json:"session_id"`
	CreatedAt            int64  `json:"created_at"`
}

type TokenResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token,omitempty"`
	ExpiresIn             int64     `json:"expires_in"`
	TokenType             TokenType `json:"token_type"`
	RegistrationStatus    string    `json:"registration_status"`
	RegistrationProgress  int       `json:"registration_progress"`
	NextStep              string    `json:"next_step,omitempty"`
	SessionID             string    `json:"session_id"`
	RequiresAdminApproval bool      `json:"requires_admin_approval,omitempty"`
}

type RegistrationStepInfo struct {
	Step                  int    `json:"step"`
	Status                string `json:"status"`
	Description           string `json:"description"`
	RequiresAdminApproval bool   `json:"requires_admin_approval"`
	IsAccessible          bool   `json:"is_accessible"`
	IsCompleted           bool   `json:"is_completed"`
}

func GetTTLFromEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	ttl, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return ttl
}

var (
	ACCESS_TOKEN_EXPIRY  = GetTTLFromEnv("ACCESS_TOKEN_EXPIRY", 23*time.Hour)
	REFRESH_TOKEN_EXPIRY = GetTTLFromEnv("REFRESH_TOKEN_EXPIRY", 28*24*time.Hour)
	PARTIAL_TOKEN_EXPIRY = GetTTLFromEnv("PARTIAL_TOKEN_EXPIRY", 2*time.Hour)
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateJTI() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {

		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func ConstantTimeCompare(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	result := 0
	for i := 0; i < len(a); i++ {
		result |= int(a[i]) ^ int(b[i])
	}
	return result == 0
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is empty")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}

func GenerateAccessToken(userID, role, deviceID, registrationStatus string, registrationProgress int, tokenType TokenType, sessionID string) (string, error) {
	secretKey := config.GetSecretKey()
	if secretKey == "" {
		return "", fmt.Errorf("secret key not found")
	}

	if userID == "" || role == "" || deviceID == "" {
		return "", fmt.Errorf("required fields cannot be empty")
	}

	now := time.Now()
	expiry := ACCESS_TOKEN_EXPIRY
	if tokenType == TokenTypePartial {
		expiry = PARTIAL_TOKEN_EXPIRY
	}

	claims := JWTClaims{
		UserID:               userID,
		Role:                 role,
		DeviceID:             deviceID,
		RegistrationStatus:   registrationStatus,
		RegistrationProgress: registrationProgress,
		TokenType:            tokenType,
		SessionID:            sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "rijig-api",
			Subject:   userID,
			ID:        GenerateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	secretKey := config.GetSecretKey()
	if secretKey == "" {
		return nil, fmt.Errorf("secret key not found")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		if IsTokenBlacklisted(claims.ID) {
			return nil, fmt.Errorf("token has been revoked")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func ValidateTokenWithChecks(tokenString string, requiredTokenType TokenType, requireCompleteReg bool) (*JWTClaims, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	if requiredTokenType != "" && claims.TokenType != requiredTokenType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", requiredTokenType, claims.TokenType)
	}

	if requireCompleteReg && !IsRegistrationComplete(claims.RegistrationStatus) {
		return nil, fmt.Errorf("registration not complete")
	}

	return claims, nil
}

func ValidateTokenForStep(tokenString string, role string, requiredStep int) (*JWTClaims, error) {
	claims, err := ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Role != role {
		return nil, fmt.Errorf("role mismatch")
	}

	if claims.RegistrationProgress < requiredStep {
		return nil, fmt.Errorf("step not accessible yet: current step %d, required step %d",
			claims.RegistrationProgress, requiredStep)
	}

	return claims, nil
}

func StoreRefreshToken(userID, deviceID, refreshToken, role, registrationStatus string, registrationProgress int, sessionID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)

	DeleteCache(key)

	data := RefreshTokenData{
		RefreshToken:         refreshToken,
		ExpiresAt:            time.Now().Add(REFRESH_TOKEN_EXPIRY).Unix(),
		DeviceID:             deviceID,
		Role:                 role,
		RegistrationStatus:   registrationStatus,
		RegistrationProgress: registrationProgress,
		SessionID:            sessionID,
		CreatedAt:            time.Now().Unix(),
	}

	err := SetCache(key, data, REFRESH_TOKEN_EXPIRY)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %v", err)
	}

	return nil
}

func ValidateRefreshToken(userID, deviceID, refreshToken string) (*RefreshTokenData, error) {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)

	var data RefreshTokenData
	err := GetCache(key, &data)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found or invalid")
	}

	if !ConstantTimeCompare(data.RefreshToken, refreshToken) {
		return nil, fmt.Errorf("refresh token mismatch")
	}

	if time.Now().Unix() > data.ExpiresAt {
		DeleteCache(key)
		return nil, fmt.Errorf("refresh token expired")
	}

	return &data, nil
}

func RefreshAccessToken(userID, deviceID, refreshToken string) (*TokenResponse, error) {
	data, err := ValidateRefreshToken(userID, deviceID, refreshToken)
	if err != nil {
		return nil, err
	}

	tokenType := DetermineTokenType(data.RegistrationStatus)

	accessToken, err := GenerateAccessToken(
		userID,
		data.Role,
		deviceID,
		data.RegistrationStatus,
		data.RegistrationProgress,
		tokenType,
		data.SessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %v", err)
	}

	expiry := ACCESS_TOKEN_EXPIRY
	if tokenType == TokenTypePartial {
		expiry = PARTIAL_TOKEN_EXPIRY
	}

	nextStep := GetNextRegistrationStep(data.Role, data.RegistrationProgress, data.RegistrationStatus)
	requiresAdminApproval := RequiresAdminApproval(data.Role, data.RegistrationProgress, data.RegistrationStatus)

	return &TokenResponse{
		AccessToken:           accessToken,
		ExpiresIn:             int64(expiry.Seconds()),
		TokenType:             tokenType,
		RegistrationStatus:    data.RegistrationStatus,
		RegistrationProgress:  data.RegistrationProgress,
		NextStep:              nextStep,
		SessionID:             data.SessionID,
		RequiresAdminApproval: requiresAdminApproval,
	}, nil
}

func GenerateTokenPair(userID, role, deviceID, registrationStatus string, registrationProgress int) (*TokenResponse, error) {
	if userID == "" || role == "" || deviceID == "" {
		return nil, fmt.Errorf("required parameters cannot be empty")
	}

	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %v", err)
	}

	tokenType := DetermineTokenType(registrationStatus)

	accessToken, err := GenerateAccessToken(
		userID,
		role,
		deviceID,
		registrationStatus,
		registrationProgress,
		tokenType,
		sessionID,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = StoreRefreshToken(userID, deviceID, refreshToken, role, registrationStatus, registrationProgress, sessionID)
	if err != nil {
		return nil, err
	}

	expiry := ACCESS_TOKEN_EXPIRY
	if tokenType == TokenTypePartial {
		expiry = PARTIAL_TOKEN_EXPIRY
	}

	nextStep := GetNextRegistrationStep(role, registrationProgress, registrationStatus)
	requiresAdminApproval := RequiresAdminApproval(role, registrationProgress, registrationStatus)

	return &TokenResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		ExpiresIn:             int64(expiry.Seconds()),
		TokenType:             tokenType,
		RegistrationStatus:    registrationStatus,
		RegistrationProgress:  registrationProgress,
		NextStep:              nextStep,
		SessionID:             sessionID,
		RequiresAdminApproval: requiresAdminApproval,
	}, nil
}

func GenerateTokenForRole(userID, role, deviceID string, progress int, status string) (*TokenResponse, error) {
	switch role {
	case RoleAdministrator:

		return GenerateTokenPair(userID, role, deviceID, RegStatusComplete, ProgressComplete)
	default:

		return GenerateTokenPair(userID, role, deviceID, status, progress)
	}
}

func DetermineTokenType(registrationStatus string) TokenType {
	if registrationStatus == RegStatusComplete {
		return TokenTypeFull
	}
	return TokenTypePartial
}

func IsRegistrationComplete(registrationStatus string) bool {
	return registrationStatus == RegStatusComplete
}

func RequiresAdminApproval(role string, progress int, status string) bool {
	switch role {
	case RolePengelola, RolePengepul:
		return progress == ProgressDataSubmitted && status == RegStatusPending
	default:
		return false
	}
}

func GetNextRegistrationStep(role string, progress int, status string) string {
	switch role {
	case RoleAdministrator:
		return "completed"

	case RoleMasyarakat:
		switch progress {
		case ProgressOTPVerified:
			return "complete_personal_data"
		case ProgressDataSubmitted:
			return "create_pin"
		case ProgressComplete:
			return "completed"
		}

	case RolePengepul:
		switch progress {
		case ProgressOTPVerified:
			return "upload_ktp"
		case ProgressDataSubmitted:
			if status == RegStatusPending {
				return "awaiting_admin_approval"
			} else if status == RegStatusConfirmed {
				return "create_pin"
			}
		case ProgressComplete:
			return "completed"
		}

	case RolePengelola:
		switch progress {
		case ProgressOTPVerified:
			return "complete_company_data"
		case ProgressDataSubmitted:
			if status == RegStatusPending {
				return "awaiting_admin_approval"
			} else if status == RegStatusConfirmed {
				return "create_pin"
			}
		case ProgressComplete:
			return "completed"
		}
	}
	return "unknown"
}

func GetRegistrationStepInfo(role string, currentProgress int, currentStatus string) *RegistrationStepInfo {
	switch role {
	case RoleAdministrator:
		return &RegistrationStepInfo{
			Step:         ProgressComplete,
			Status:       RegStatusComplete,
			Description:  "Administrator registration complete",
			IsAccessible: true,
			IsCompleted:  true,
		}

	case RoleMasyarakat:
		return getMasyarakatStepInfo(currentProgress, currentStatus)

	case RolePengepul:
		return getPengepulStepInfo(currentProgress, currentStatus)

	case RolePengelola:
		return getPengelolaStepInfo(currentProgress, currentStatus)
	}

	return &RegistrationStepInfo{
		Step:         0,
		Status:       "unknown",
		Description:  "Unknown role",
		IsAccessible: false,
		IsCompleted:  false,
	}
}

func getMasyarakatStepInfo(progress int, status string) *RegistrationStepInfo {
	switch progress {
	case ProgressOTPVerified:
		return &RegistrationStepInfo{
			Step:         ProgressOTPVerified,
			Status:       status,
			Description:  "Complete personal data",
			IsAccessible: true,
			IsCompleted:  false,
		}
	case ProgressDataSubmitted:
		return &RegistrationStepInfo{
			Step:         ProgressDataSubmitted,
			Status:       status,
			Description:  "Create PIN",
			IsAccessible: true,
			IsCompleted:  false,
		}
	case ProgressComplete:
		return &RegistrationStepInfo{
			Step:         ProgressComplete,
			Status:       RegStatusComplete,
			Description:  "Registration complete",
			IsAccessible: true,
			IsCompleted:  true,
		}
	}
	return nil
}

func getPengepulStepInfo(progress int, status string) *RegistrationStepInfo {
	switch progress {
	case ProgressOTPVerified:
		return &RegistrationStepInfo{
			Step:         ProgressOTPVerified,
			Status:       status,
			Description:  "Upload KTP",
			IsAccessible: true,
			IsCompleted:  false,
		}
	case ProgressDataSubmitted:
		if status == RegStatusPending {
			return &RegistrationStepInfo{
				Step:                  ProgressDataSubmitted,
				Status:                status,
				Description:           "Awaiting admin approval",
				RequiresAdminApproval: true,
				IsAccessible:          false,
				IsCompleted:           false,
			}
		} else if status == RegStatusConfirmed {
			return &RegistrationStepInfo{
				Step:         ProgressDataSubmitted,
				Status:       status,
				Description:  "Create PIN",
				IsAccessible: true,
				IsCompleted:  false,
			}
		}
	case ProgressComplete:
		return &RegistrationStepInfo{
			Step:         ProgressComplete,
			Status:       RegStatusComplete,
			Description:  "Registration complete",
			IsAccessible: true,
			IsCompleted:  true,
		}
	}
	return nil
}

func getPengelolaStepInfo(progress int, status string) *RegistrationStepInfo {
	switch progress {
	case ProgressOTPVerified:
		return &RegistrationStepInfo{
			Step:         ProgressOTPVerified,
			Status:       status,
			Description:  "Complete company data",
			IsAccessible: true,
			IsCompleted:  false,
		}
	case ProgressDataSubmitted:
		if status == RegStatusPending {
			return &RegistrationStepInfo{
				Step:                  ProgressDataSubmitted,
				Status:                status,
				Description:           "Awaiting admin approval",
				RequiresAdminApproval: true,
				IsAccessible:          false,
				IsCompleted:           false,
			}
		} else if status == RegStatusConfirmed {
			return &RegistrationStepInfo{
				Step:         ProgressDataSubmitted,
				Status:       status,
				Description:  "Create PIN",
				IsAccessible: true,
				IsCompleted:  false,
			}
		}
	case ProgressComplete:
		return &RegistrationStepInfo{
			Step:         ProgressComplete,
			Status:       RegStatusComplete,
			Description:  "Registration complete",
			IsAccessible: true,
			IsCompleted:  true,
		}
	}
	return nil
}

func RevokeRefreshToken(userID, deviceID string) error {
	key := fmt.Sprintf("refresh_token:%s:%s", userID, deviceID)
	err := DeleteCache(key)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %v", err)
	}
	return nil
}

func RevokeAllRefreshTokens(userID string) error {
	pattern := fmt.Sprintf("refresh_token:%s:*", userID)
	err := ScanAndDelete(pattern)
	if err != nil {
		return fmt.Errorf("failed to revoke all refresh tokens: %v", err)
	}
	return nil
}

func BlacklistToken(jti string, expiresAt time.Time) error {
	key := fmt.Sprintf("blacklist:%s", jti)
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	return SetCache(key, true, ttl)
}

func IsTokenBlacklisted(jti string) bool {
	key := fmt.Sprintf("blacklist:%s", jti)
	var exists bool
	err := GetCache(key, &exists)
	return err == nil && exists
}
