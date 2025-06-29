package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type ResetPasswordData struct {
	Token     string `json:"token"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
	Used      bool   `json:"used"`
	CreatedAt int64  `json:"created_at"`
}

const (
	RESET_TOKEN_EXPIRY = 30 * time.Minute
	RESET_TOKEN_LENGTH = 32
)

// Generate secure reset token
func GenerateResetToken() (string, error) {
	bytes := make([]byte, RESET_TOKEN_LENGTH)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Store reset password token di Redis
func StoreResetToken(email, userID, token string) error {
	key := fmt.Sprintf("reset_password:%s", email)

	// Delete any existing reset token for this email
	DeleteCache(key)

	data := ResetPasswordData{
		Token:     token,
		Email:     email,
		UserID:    userID,
		ExpiresAt: time.Now().Add(RESET_TOKEN_EXPIRY).Unix(),
		Used:      false,
		CreatedAt: time.Now().Unix(),
	}

	return SetCache(key, data, RESET_TOKEN_EXPIRY)
}

// Validate reset password token
func ValidateResetToken(email, inputToken string) (*ResetPasswordData, error) {
	key := fmt.Sprintf("reset_password:%s", email)

	var data ResetPasswordData
	err := GetCache(key, &data)
	if err != nil {
		return nil, fmt.Errorf("token reset tidak ditemukan atau sudah kadaluarsa")
	}

	// Check if token is expired
	if time.Now().Unix() > data.ExpiresAt {
		DeleteCache(key)
		return nil, fmt.Errorf("token reset sudah kadaluarsa")
	}

	// Check if token is already used
	if data.Used {
		return nil, fmt.Errorf("token reset sudah digunakan")
	}

	// Validate token
	if !ConstantTimeCompare(data.Token, inputToken) {
		return nil, fmt.Errorf("token reset tidak valid")
	}

	return &data, nil
}

// Mark reset token as used
func MarkResetTokenAsUsed(email string) error {
	key := fmt.Sprintf("reset_password:%s", email)

	var data ResetPasswordData
	err := GetCache(key, &data)
	if err != nil {
		return err
	}

	data.Used = true
	remaining := time.Until(time.Unix(data.ExpiresAt, 0))

	return SetCache(key, data, remaining)
}

// Check if reset token exists and still valid
func IsResetTokenValid(email string) bool {
	key := fmt.Sprintf("reset_password:%s", email)

	var data ResetPasswordData
	err := GetCache(key, &data)
	if err != nil {
		return false
	}

	return time.Now().Unix() <= data.ExpiresAt && !data.Used
}

// Get remaining reset token time
func GetResetTokenRemainingTime(email string) (time.Duration, error) {
	key := fmt.Sprintf("reset_password:%s", email)

	var data ResetPasswordData
	err := GetCache(key, &data)
	if err != nil {
		return 0, err
	}

	remaining := time.Until(time.Unix(data.ExpiresAt, 0))
	if remaining < 0 {
		return 0, fmt.Errorf("token expired")
	}

	return remaining, nil
}

// Send reset password email
func (e *EmailService) SendResetPasswordEmail(email, name, token string) error {
	// Create reset URL - in real app this would be frontend URL
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s&email=%s", token, email)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.from, e.fromName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Reset Password Administrator - Rijig")

	// Email template
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        .container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
        .header { background-color: #d32f2f; color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; background-color: #f9f9f9; }
        .reset-button { display: inline-block; background-color: #d32f2f; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-weight: bold; margin: 20px 0; }
        .reset-button:hover { background-color: #b71c1c; }
        .token-box { font-size: 14px; color: #666; background-color: white; padding: 15px; border-left: 4px solid #d32f2f; margin: 20px 0; word-break: break-all; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .warning { color: #d32f2f; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Reset Password</h1>
        </div>
        <div class="content">
            <p>Halo <strong>%s</strong>,</p>
            <p>Kami menerima permintaan untuk reset password akun Administrator Anda.</p>
            
            <p>Klik tombol di bawah ini untuk reset password:</p>
            <div style="text-align: center;">
                <a href="%s" class="reset-button">Reset Password</a>
            </div>
            
            <p>Atau copy paste link berikut ke browser Anda:</p>
            <div class="token-box">%s</div>
            
            <p><strong>Penting:</strong></p>
            <ul>
                <li>Link ini berlaku selama <strong>30 menit</strong></li>
                <li>Link hanya dapat digunakan sekali</li>
                <li>Jangan bagikan link ini kepada siapapun</li>
            </ul>
            
            <p class="warning">⚠️ Jika Anda tidak melakukan permintaan reset password, abaikan email ini dan password Anda tidak akan berubah.</p>
        </div>
        <div class="footer">
            <p>Email ini dikirim otomatis oleh sistem Rijig Waste Management<br>
            Jangan balas email ini.</p>
        </div>
    </div>
</body>
</html>
	`, name, resetURL, resetURL)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.host, e.port, e.username, e.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send reset password email: %v", err)
	}

	return nil
}
