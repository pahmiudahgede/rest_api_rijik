package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

type OTPData struct {
	Code      string `json:"code"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expires_at"`
	Attempts  int    `json:"attempts"`
	CreatedAt int64  `json:"created_at"`
}

const (
	OTP_LENGTH       = 6
	OTP_EXPIRY       = 5 * time.Minute
	MAX_OTP_ATTEMPTS = 3
)

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi("587")

	return &EmailService{
		host:     "smtp.gmail.com",
		port:     port,
		username: os.Getenv("SMTP_FROM_EMAIL"),
		password: os.Getenv("GMAIL_APP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM_EMAIL"),
		fromName: os.Getenv("SMTP_FROM_NAME"),
	}
}

func StoreOTP(email, otp string) error {
	key := fmt.Sprintf("otp:admin:%s", email)

	data := OTPData{
		Code:      otp,
		Email:     email,
		ExpiresAt: time.Now().Add(OTP_EXPIRY).Unix(),
		Attempts:  0,
		CreatedAt: time.Now().Unix(),
	}

	return SetCache(key, data, OTP_EXPIRY)
}

func ValidateOTP(email, inputOTP string) error {
	key := fmt.Sprintf("otp:admin:%s", email)

	var data OTPData
	err := GetCache(key, &data)
	if err != nil {
		return fmt.Errorf("OTP tidak ditemukan atau sudah kadaluarsa")
	}

	if time.Now().Unix() > data.ExpiresAt {
		DeleteCache(key)
		return fmt.Errorf("OTP sudah kadaluarsa")
	}

	if data.Attempts >= MAX_OTP_ATTEMPTS {
		DeleteCache(key)
		return fmt.Errorf("OTP diblokir karena terlalu banyak percobaan salah")
	}

	if data.Code != inputOTP {

		data.Attempts++
		SetCache(key, data, time.Until(time.Unix(data.ExpiresAt, 0)))
		return fmt.Errorf("OTP tidak valid. Sisa percobaan: %d", MAX_OTP_ATTEMPTS-data.Attempts)
	}

	DeleteCache(key)
	return nil
}

func (e *EmailService) SendOTPEmail(email, name, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.from, e.fromName))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Kode Verifikasi Login Administrator - Rijig")

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        .container { max-width: 600px; margin: 0 auto; font-family: Arial, sans-serif; }
        .header { background-color: #2E7D32; color: white; padding: 20px; text-align: center; }
        .content { padding: 30px; background-color: #f9f9f9; }
        .otp-code { font-size: 32px; font-weight: bold; color: #2E7D32; text-align: center; letter-spacing: 5px; margin: 20px 0; padding: 15px; background-color: white; border: 2px dashed #2E7D32; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .warning { color: #d32f2f; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Kode Verifikasi Login</h1>
        </div>
        <div class="content">
            <p>Halo <strong>%s</strong>,</p>
            <p>Anda telah meminta untuk login sebagai Administrator. Gunakan kode verifikasi berikut:</p>
            
            <div class="otp-code">%s</div>
            
            <p><strong>Penting:</strong></p>
            <ul>
                <li>Kode ini berlaku selama <strong>5 menit</strong></li>
                <li>Jangan berikan kode ini kepada siapapun</li>
                <li>Maksimal 3 kali percobaan</li>
            </ul>
            
            <p class="warning">⚠️ Jika Anda tidak melakukan permintaan login ini, abaikan email ini.</p>
        </div>
        <div class="footer">
            <p>Email ini dikirim otomatis oleh sistem Rijig Waste Management<br>
            Jangan balas email ini.</p>
        </div>
    </div>
</body>
</html>
	`, name, otp)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.host, e.port, e.username, e.password)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func IsOTPValid(email string) bool {
	key := fmt.Sprintf("otp:admin:%s", email)

	var data OTPData
	err := GetCache(key, &data)
	if err != nil {
		return false
	}

	return time.Now().Unix() <= data.ExpiresAt && data.Attempts < MAX_OTP_ATTEMPTS
}

func GetOTPRemainingTime(email string) (time.Duration, error) {
	key := fmt.Sprintf("otp:admin:%s", email)

	var data OTPData
	err := GetCache(key, &data)
	if err != nil {
		return 0, err
	}

	remaining := time.Until(time.Unix(data.ExpiresAt, 0))
	if remaining < 0 {
		return 0, fmt.Errorf("OTP expired")
	}

	return remaining, nil
}
