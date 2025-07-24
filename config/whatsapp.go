package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsAppService struct {
	Client       *whatsmeow.Client
	Container    *sqlstore.Container
	QRChan       chan whatsmeow.QRChannelItem
	Connected    chan bool
	mu           sync.RWMutex
	isConnecting bool
	loginSuccess chan bool
}

var whatsappService *WhatsAppService

func InitWhatsApp() {
	var err error
	var connectionString string

	// Check if DATABASE_URL is available (Railway provides this)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		connectionString = processDatabaseURL(databaseURL)
		log.Println("Using DATABASE_URL for connection")
	} else {
		// Fallback to individual environment variables
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")

		// Validate required environment variables
		if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" {
			log.Fatal("Missing required database environment variables")
		}

		connectionString = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName,
		)
		log.Println("Using individual DB environment variables")
	}

	log.Printf("Attempting to connect to database...")

	// Use more appropriate logging level
	logLevel := "INFO"
	if os.Getenv("DEBUG") == "true" {
		logLevel = "DEBUG"
	}
	dbLog := waLog.Stdout("Database", logLevel, true)

	// Add context with longer timeout for production
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	container, err := sqlstore.New(ctx, "postgres", connectionString, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	log.Println("Successfully connected to WhatsApp database")

	whatsappService = &WhatsAppService{
		Container:    container,
		Connected:    make(chan bool, 1),
		loginSuccess: make(chan bool, 1),
	}
}

func GetWhatsAppService() *WhatsAppService {
	return whatsappService
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		log.Printf("📨 Received message: %s", v.Message.GetConversation())
	case *events.Connected:
		log.Println("✅ WhatsApp client connected!")
		if whatsappService != nil && whatsappService.Connected != nil {
			whatsappService.mu.Lock()
			whatsappService.isConnecting = false
			whatsappService.mu.Unlock()

			select {
			case whatsappService.Connected <- true:
			default:
			}
		}
	case *events.Disconnected:
		log.Println("❌ WhatsApp client disconnected!")
		if whatsappService != nil {
			whatsappService.mu.Lock()
			whatsappService.isConnecting = false
			whatsappService.mu.Unlock()
		}
	case *events.LoggedOut:
		log.Println("🚪 WhatsApp client logged out!")
	case *events.PairSuccess:
		log.Println("🎉 WhatsApp pairing successful!")
		if whatsappService != nil && whatsappService.loginSuccess != nil {
			select {
			case whatsappService.loginSuccess <- true:
			default:
			}
		}
	case *events.ConnectFailure:
		log.Printf("❌ Connection failure: %v", v.Reason)
	case *events.StreamError:
		log.Printf("🌊 Stream error: %v", v.Code)
	case *events.QR:
		// events.QR memiliki field Codes ([]string), bukan Event
		if len(v.Codes) > 0 {
			log.Printf("📱 QR Code received with %d codes", len(v.Codes))
		} else {
			log.Printf("📱 QR Code event (no codes available)")
		}
	case *events.QRScannedWithoutMultidevice:
		log.Println("📱 QR scanned but multidevice not enabled")
	default:
		log.Printf("📱 WhatsApp event: %T", v)
	}
}

func (wa *WhatsAppService) GenerateQR() (string, error) {
	wa.mu.Lock()
	defer wa.mu.Unlock()

	if wa.Container == nil {
		return "", fmt.Errorf("container is not initialized")
	}

	// Prevent multiple concurrent QR generations
	if wa.isConnecting {
		return "", fmt.Errorf("QR generation already in progress")
	}

	// Cleanup existing client if any
	if wa.Client != nil {
		log.Println("🧹 Cleaning up existing client...")
		wa.Client.Disconnect()
		time.Sleep(2 * time.Second) // Give time for cleanup
		wa.Client = nil
	}

	wa.isConnecting = true
	log.Println("🔄 Starting QR generation...")

	// Add context with timeout for GetFirstDevice - sesuai dokumentasi terbaru
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	deviceStore, err := wa.Container.GetFirstDevice(ctx)
	if err != nil {
		wa.isConnecting = false
		return "", fmt.Errorf("failed to get first device: %v", err)
	}

	// Use appropriate logging level
	logLevel := "INFO"
	if os.Getenv("DEBUG") == "true" {
		logLevel = "DEBUG"
	}
	clientLog := waLog.Stdout("Client", logLevel, true)

	wa.Client = whatsmeow.NewClient(deviceStore, clientLog)
	wa.Client.AddEventHandler(eventHandler)

	// Set properties untuk production stability - sesuai dokumentasi terbaru
	wa.Client.EnableAutoReconnect = true
	wa.Client.AutoTrustIdentity = true

	if wa.Client.Store.ID == nil {
		log.Println("📱 Client is not logged in, generating QR code...")

		// PENTING: Gunakan context yang tidak akan di-cancel untuk QR monitoring
		// Context ini harus hidup selama proses scan QR berlangsung
		qrCtx := context.Background() // Tidak menggunakan timeout untuk monitoring

		qrChan, err := wa.Client.GetQRChannel(qrCtx)
		if err != nil {
			wa.isConnecting = false
			return "", fmt.Errorf("failed to get QR channel: %v", err)
		}

		// Connect dengan error handling yang lebih baik
		err = wa.Client.Connect()
		if err != nil {
			wa.isConnecting = false
			return "", fmt.Errorf("failed to connect: %v", err)
		}

		// Wait for FIRST QR code dengan timeout yang wajar
		qrTimeout := time.NewTimer(60 * time.Second)
		defer qrTimeout.Stop()

		select {
		case evt := <-qrChan:
			switch evt.Event {
			case "code":
				log.Println("📱 QR code generated successfully")
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					wa.isConnecting = false
					return "", fmt.Errorf("failed to create QR code image: %v", err)
				}

				// KUNCI: Start monitoring di background TANPA cancel context
				// Ini akan terus berjalan sampai login berhasil atau gagal
				go wa.handleQREventsEnhanced(qrChan)

				dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
				log.Println("✅ QR code ready for scanning")
				log.Println("⏳ Waiting for QR scan... (QR monitoring is active)")
				return dataURI, nil

			case "success":
				log.Println("🎉 Login successful immediately!")
				wa.isConnecting = false
				return "success", nil

			case "err-client-outdated":
				wa.isConnecting = false
				return "", fmt.Errorf("WhatsApp client is outdated. Please update whatsmeow library")

			case "err-client-banned":
				wa.isConnecting = false
				return "", fmt.Errorf("WhatsApp client is banned")

			case "timeout":
				wa.isConnecting = false
				return "", fmt.Errorf("QR code generation timeout")

			default:
				log.Printf("🔄 QR Login event: %s", evt.Event)
				if evt.Error != nil {
					wa.isConnecting = false
					return "", fmt.Errorf("QR login error: %v", evt.Error)
				}
			}

		case <-qrTimeout.C:
			wa.isConnecting = false
			return "", fmt.Errorf("timeout waiting for first QR code")
		}
	} else {
		log.Println("📱 Client already logged in, attempting to connect...")

		err = wa.Client.Connect()
		if err != nil {
			wa.isConnecting = false
			return "", fmt.Errorf("failed to connect existing session: %v", err)
		}

		// Wait longer for connection establishment in production
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			if wa.Client.IsConnected() {
				wa.isConnecting = false
				log.Println("✅ Successfully connected with existing session")
				return "already_connected", nil
			}
		}

		wa.isConnecting = false
		return "", fmt.Errorf("failed to establish connection with existing session")
	}

	wa.isConnecting = false
	return "", fmt.Errorf("unexpected end of QR generation process")
}

// Enhanced QR event handler TANPA context timeout untuk memungkinkan scan yang lama
func (wa *WhatsAppService) handleQREventsEnhanced(qrChan <-chan whatsmeow.QRChannelItem) {
	log.Println("🔍 Starting QR event monitoring (no timeout)...")

	// Monitor QR events sampai login berhasil atau channel ditutup
	for evt := range qrChan {
		switch evt.Event {
		case "success":
			log.Println("🎉 Login successful after QR scan!")
			wa.mu.Lock()
			wa.isConnecting = false
			wa.mu.Unlock()

			select {
			case wa.loginSuccess <- true:
			default:
			}

			// Start keep-alive setelah login berhasil
			go wa.StartKeepAlive()
			return

		case "timeout":
			log.Println("⏰ QR code scan timeout - generating new QR code")
			// Jangan return, biarkan QR baru di-generate

		case "err-client-outdated":
			log.Println("📱 Client outdated error")
			wa.mu.Lock()
			wa.isConnecting = false
			wa.mu.Unlock()
			return

		case "err-client-banned":
			log.Println("🚫 Client banned error")
			wa.mu.Lock()
			wa.isConnecting = false
			wa.mu.Unlock()
			return

		case "err-scanned-without-multidevice":
			log.Println("📱 QR scanned but multidevice not enabled - user needs to enable multidevice")
			// Jangan return, biarkan user enable multidevice dan scan lagi

		case "code":
			log.Println("🔄 New QR code generated (old one expired)")

		default:
			log.Printf("📱 QR event: %s", evt.Event)
			if evt.Error != nil {
				log.Printf("❌ QR Error: %v", evt.Error)
			}
		}
	}

	log.Println("📱 QR channel closed")
	wa.mu.Lock()
	wa.isConnecting = false
	wa.mu.Unlock()
}

func (wa *WhatsAppService) WaitForLoginSuccess(timeout time.Duration) error {
	log.Printf("⏳ Waiting for login success (timeout: %v)...", timeout)

	select {
	case <-wa.loginSuccess:
		log.Println("✅ Login success confirmed!")
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for login success")
	}
}

func (wa *WhatsAppService) SendMessage(phoneNumber, message string) error {
	if wa.Client == nil {
		return fmt.Errorf("client not initialized")
	}

	if !wa.Client.IsConnected() {
		return fmt.Errorf("client not connected")
	}

	targetJID, err := types.ParseJID(phoneNumber + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid phone number: %v", err)
	}

	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = wa.Client.SendMessage(ctx, targetJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("✅ Message sent to %s", phoneNumber)
	return nil
}

func (wa *WhatsAppService) Logout() error {
	wa.mu.Lock()
	defer wa.mu.Unlock()

	if wa.Client == nil {
		return fmt.Errorf("no active client session")
	}

	log.Println("🚪 Logging out...")

	// Add context with timeout for Logout - sesuai dokumentasi terbaru
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := wa.Client.Logout(ctx)
	if err != nil {
		log.Printf("❌ Logout error: %v", err)
		// Continue with disconnection even if logout fails
	}

	wa.Client.Disconnect()
	wa.Client = nil
	wa.isConnecting = false

	log.Println("✅ Logged out successfully")
	return err
}

func (wa *WhatsAppService) IsConnected() bool {
	wa.mu.RLock()
	defer wa.mu.RUnlock()
	return wa.Client != nil && wa.Client.IsConnected()
}

func (wa *WhatsAppService) IsLoggedIn() bool {
	wa.mu.RLock()
	defer wa.mu.RUnlock()
	return wa.Client != nil && wa.Client.Store.ID != nil
}

func (wa *WhatsAppService) IsConnecting() bool {
	wa.mu.RLock()
	defer wa.mu.RUnlock()
	return wa.isConnecting
}

func (wa *WhatsAppService) GetLoginStatus() map[string]interface{} {
	wa.mu.RLock()
	defer wa.mu.RUnlock()

	status := map[string]interface{}{
		"is_connected":  wa.IsConnected(),
		"is_logged_in":  wa.IsLoggedIn(),
		"is_connecting": wa.isConnecting,
		"timestamp":     time.Now().Unix(),
	}

	if wa.Client != nil && wa.Client.Store.ID != nil {
		status["device_id"] = wa.Client.Store.ID.User
	}

	// Add QR monitoring status
	if wa.isConnecting {
		status["qr_monitoring_active"] = true
		status["message"] = "QR code is ready for scanning. Please scan with WhatsApp mobile app."
	} else if wa.IsLoggedIn() && wa.IsConnected() {
		status["message"] = "Successfully logged in and connected."
	} else if wa.IsLoggedIn() && !wa.IsConnected() {
		status["message"] = "Logged in but not connected. Attempting to reconnect..."
	} else {
		status["message"] = "Not logged in. Please generate QR code."
	}

	return status
}

func (wa *WhatsAppService) WaitForLogin(timeout time.Duration) error {
	if wa.Client == nil {
		return fmt.Errorf("client not initialized")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeoutChan := time.After(timeout)

	for {
		select {
		case <-ticker.C:
			if wa.IsLoggedIn() && wa.IsConnected() {
				return nil
			}
		case <-timeoutChan:
			return fmt.Errorf("timeout waiting for login")
		}
	}
}

func (wa *WhatsAppService) Disconnect() {
	wa.mu.Lock()
	defer wa.mu.Unlock()

	if wa.Client != nil {
		log.Println("🔌 Disconnecting client...")
		wa.Client.Disconnect()
		wa.Client = nil
	}
	wa.isConnecting = false
}

func (wa *WhatsAppService) GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("🛑 Shutting down WhatsApp client...")
		if wa.Client != nil {
			wa.Disconnect()
		}
		os.Exit(0)
	}()
}

// Keep-alive mechanism untuk menjaga koneksi tetap hidup di production
func (wa *WhatsAppService) StartKeepAlive() {
	if wa.Client == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if wa.IsConnected() {
				// Send presence untuk keep connection alive
				wa.Client.SendPresence(types.PresenceAvailable)
				log.Println("🔄 Keep-alive signal sent")
			}
		}
	}()
}
