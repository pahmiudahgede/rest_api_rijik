package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailVerificationData struct {
	Token     string `json:"token"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
	Used      bool   `json:"used"`
	CreatedAt int64  `json:"created_at"`
}

const (
	EMAIL_VERIFICATION_TOKEN_EXPIRY = 24 * time.Hour
	EMAIL_VERIFICATION_TOKEN_LENGTH = 32
)

func GenerateEmailVerificationToken() (string, error) {
	bytes := make([]byte, EMAIL_VERIFICATION_TOKEN_LENGTH)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate email verification token: %v", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func StoreEmailVerificationToken(email, userID, token string) error {
	key := fmt.Sprintf("email_verification:%s", email)

	DeleteCache(key)

	data := EmailVerificationData{
		Token:     token,
		Email:     email,
		UserID:    userID,
		ExpiresAt: time.Now().Add(EMAIL_VERIFICATION_TOKEN_EXPIRY).Unix(),
		Used:      false,
		CreatedAt: time.Now().Unix(),
	}

	return SetCache(key, data, EMAIL_VERIFICATION_TOKEN_EXPIRY)
}

func ValidateEmailVerificationToken(email, inputToken string) (*EmailVerificationData, error) {
	key := fmt.Sprintf("email_verification:%s", email)

	var data EmailVerificationData
	err := GetCache(key, &data)
	if err != nil {
		return nil, fmt.Errorf("token verifikasi tidak ditemukan atau sudah kadaluarsa")
	}

	if time.Now().Unix() > data.ExpiresAt {
		DeleteCache(key)
		return nil, fmt.Errorf("token verifikasi sudah kadaluarsa")
	}

	if data.Used {
		return nil, fmt.Errorf("token verifikasi sudah digunakan")
	}

	// Validate token
	if !ConstantTimeCompare(data.Token, inputToken) {
		return nil, fmt.Errorf("token verifikasi tidak valid")
	}

	return &data, nil
}

// Mark email verification token as used
func MarkEmailVerificationTokenAsUsed(email string) error {
	key := fmt.Sprintf("email_verification:%s", email)

	var data EmailVerificationData
	err := GetCache(key, &data)
	if err != nil {
		return err
	}

	data.Used = true
	remaining := time.Until(time.Unix(data.ExpiresAt, 0))

	return SetCache(key, data, remaining)
}

// Check if email verification token exists and still valid
func IsEmailVerificationTokenValid(email string) bool {
	key := fmt.Sprintf("email_verification:%s", email)

	var data EmailVerificationData
	err := GetCache(key, &data)
	if err != nil {
		return false
	}

	return time.Now().Unix() <= data.ExpiresAt && !data.Used
}

// Get remaining email verification token time
func GetEmailVerificationTokenRemainingTime(email string) (time.Duration, error) {
	key := fmt.Sprintf("email_verification:%s", email)

	var data EmailVerificationData
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

// Send email verification email
func (e *EmailService) SendEmailVerificationEmail(email, name, token string) error {
	// Create verification URL - in real app this would be frontend URL
	verificationURL := fmt.Sprintf("http://localhost:3000/verify-email?token=%s&email=%s", token, email)

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.from, e.fromName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikasi Email Administrator - Rijig")

	// Email template
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        .container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
        .header { background-color: #2E7D32; color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; background-color: #f9f9f9; }
        .verify-button { display: inline-block; background-color: #2E7D32; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; font-weight: bold; margin: 20px 0; }
        .verify-button:hover { background-color: #1B5E20; }
        .token-box { font-size: 14px; color: #666; background-color: white; padding: 15px; border-left: 4px solid #2E7D32; margin: 20px 0; word-break: break-all; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .success-icon { font-size: 48px; text-align: center; margin: 20px 0; }
        .info-box { background-color: #E8F5E8; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Verifikasi Email</h1>
        </div>
        <div class="content">
            <div class="success-icon">🎉</div>
            
            <p>Selamat <strong>%s</strong>!</p>
            <p>Akun Administrator Anda telah berhasil dibuat. Untuk mengaktifkan akun dan mulai menggunakan sistem Rijig, silakan verifikasi email Anda dengan mengklik tombol di bawah ini:</p>
            
            <div style="text-align: center;">
                <a href="%s" class="verify-button">Verifikasi Email Saya</a>
            </div>
            
            <p>Atau copy paste link berikut ke browser Anda:</p>
            <div class="token-box">%s</div>
            
            <div class="info-box">
                <p><strong>Informasi Penting:</strong></p>
                <ul>
                    <li>Link verifikasi berlaku selama <strong>24 jam</strong></li>
                    <li>Setelah verifikasi, Anda dapat login ke sistem</li>
                    <li>Link hanya dapat digunakan sekali</li>
                    <li>Jangan bagikan link ini kepada siapapun</li>
                </ul>
            </div>
            
            <p><strong>Langkah selanjutnya setelah verifikasi:</strong></p>
            <ol>
                <li>Login menggunakan email dan password</li>
                <li>Masukkan kode OTP yang dikirim ke email</li>
                <li>Mulai menggunakan sistem Rijig</li>
            </ol>
            
            <p style="color: #666; font-style: italic;">Jika Anda tidak membuat akun ini, abaikan email ini.</p>
        </div>
        <div class="footer">
            <p>Email ini dikirim otomatis oleh sistem Rijig Waste Management<br>
            Jangan balas email ini.</p>
        </div>
    </div>
</body>
</html>
	`, name, verificationURL, verificationURL)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.host, e.port, e.username, e.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email verification email: %v", err)
	}

	return nil
}
