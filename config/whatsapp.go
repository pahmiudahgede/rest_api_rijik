package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	Client    *whatsmeow.Client
	Container *sqlstore.Container
	QRChan    chan whatsmeow.QRChannelItem
	Connected chan bool
}

var whatsappService *WhatsAppService

func InitWhatsApp() {
	var err error
	var connectionString string

	// Check if DATABASE_URL is available (Railway provides this)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		// connectionString = databaseURL
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

	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Add context with timeout for database connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	container, err := sqlstore.New(ctx, "postgres", connectionString, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	log.Println("Successfully connected to WhatsApp database")

	whatsappService = &WhatsAppService{
		Container: container,
		Connected: make(chan bool, 1),
	}
}

func GetWhatsAppService() *WhatsAppService {
	return whatsappService
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	case *events.Connected:
		fmt.Println("WhatsApp client connected!")
		if whatsappService != nil && whatsappService.Connected != nil {
			select {
			case whatsappService.Connected <- true:
			default:
			}
		}
	case *events.Disconnected:
		fmt.Println("WhatsApp client disconnected!")
	case *events.LoggedOut:
		fmt.Println("WhatsApp client logged out!")
	case *events.PairSuccess:
		fmt.Println("WhatsApp pairing successful!")
	}
}

func (wa *WhatsAppService) GenerateQR() (string, error) {
	if wa.Container == nil {
		return "", fmt.Errorf("container is not initialized")
	}

	// Cleanup existing client if any
	if wa.Client != nil {
		wa.Client.Disconnect()
		wa.Client = nil
	}

	// Add context with timeout for GetFirstDevice
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deviceStore, err := wa.Container.GetFirstDevice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get first device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	wa.Client = whatsmeow.NewClient(deviceStore, clientLog)
	wa.Client.AddEventHandler(eventHandler)

	if wa.Client.Store.ID == nil {
		fmt.Println("Client is not logged in, generating QR code...")

		// Use context with reasonable timeout for QR generation
		qrCtx, qrCancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer qrCancel()

		qrChan, err := wa.Client.GetQRChannel(qrCtx)
		if err != nil {
			return "", fmt.Errorf("failed to get QR channel: %v", err)
		}

		err = wa.Client.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect: %v", err)
		}

		// Wait for the first QR code
		select {
		case evt := <-qrChan:
			switch evt.Event {
			case "code":
				fmt.Println("QR code generated:", evt.Code)
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					return "", fmt.Errorf("failed to create QR code: %v", err)
				}

				// Start goroutine to handle subsequent events (like success)
				go wa.handleQREvents(qrChan)

				dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
				return dataURI, nil

			case "success":
				fmt.Println("Login successful!")
				return "success", nil

			case "err-client-outdated":
				return "", fmt.Errorf("WhatsApp client is outdated. Please update the whatsmeow library")

			case "err-client-banned":
				return "", fmt.Errorf("WhatsApp client is banned")

			default:
				fmt.Printf("Login event: %s\n", evt.Event)
				if evt.Error != nil {
					return "", fmt.Errorf("login error: %v", evt.Error)
				}
			}

		case <-qrCtx.Done():
			return "", fmt.Errorf("timeout waiting for QR code: %v", qrCtx.Err())
		}
	} else {
		fmt.Println("Client already logged in, connecting...")
		err = wa.Client.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect: %v", err)
		}

		// Wait a bit to ensure connection is established
		time.Sleep(3 * time.Second)

		if wa.Client.IsConnected() {
			return "already_connected", nil
		} else {
			return "", fmt.Errorf("failed to establish connection")
		}
	}

	return "", fmt.Errorf("failed to generate QR code")
}

// Handle QR events in background
func (wa *WhatsAppService) handleQREvents(qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		switch evt.Event {
		case "success":
			fmt.Println("Login successful after QR scan!")
		case "timeout":
			fmt.Println("QR code scan timeout")
		case "err-client-outdated":
			fmt.Println("Client outdated error")
		case "err-client-banned":
			fmt.Println("Client banned error")
		default:
			fmt.Printf("Background login event: %s\n", evt.Event)
		}
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

	return nil
}

func (wa *WhatsAppService) Logout() error {
	if wa.Client == nil {
		return fmt.Errorf("no active client session")
	}

	// Add context with timeout for Logout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := wa.Client.Logout(ctx)
	if err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}

	wa.Client.Disconnect()
	wa.Client = nil
	return nil
}

func (wa *WhatsAppService) IsConnected() bool {
	return wa.Client != nil && wa.Client.IsConnected()
}

func (wa *WhatsAppService) IsLoggedIn() bool {
	return wa.Client != nil && wa.Client.Store.ID != nil
}

func (wa *WhatsAppService) GetLoginStatus() map[string]interface{} {
	status := map[string]interface{}{
		"is_connected": wa.IsConnected(),
		"is_logged_in": wa.IsLoggedIn(),
		"timestamp":    time.Now().Unix(),
	}

	if wa.Client != nil && wa.Client.Store.ID != nil {
		status["device_id"] = wa.Client.Store.ID.User
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
	if wa.Client != nil {
		wa.Client.Disconnect()
		wa.Client = nil
	}
}

func (wa *WhatsAppService) GracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("Shutting down WhatsApp client...")
		if wa.Client != nil {
			wa.Client.Disconnect()
		}
		os.Exit(0)
	}()
}
