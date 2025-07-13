package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
}

var whatsappService *WhatsAppService

func InitWhatsApp() {
	var err error
	var connectionString string

	// Check if running on Railway (DATABASE_URL provided)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Println("Using Railway DATABASE_URL for WhatsApp")
		connectionString, err = convertURLToLibPQFormat(databaseURL)
		if err != nil {
			log.Fatalf("Failed to convert DATABASE_URL for WhatsApp: %v", err)
		}
	} else {
		// Fallback to individual environment variables (for local development)
		log.Println("Using individual database environment variables for WhatsApp")
		connectionString = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
		)
	}

	log.Printf("WhatsApp connecting to database with DSN: %s", sanitizeConnectionString(connectionString))

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("postgres", connectionString, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	whatsappService = &WhatsAppService{
		Container: container,
	}

	log.Println("WhatsApp database connected successfully!")
}

// convertURLToLibPQFormat converts DATABASE_URL to lib/pq connection string format
func convertURLToLibPQFormat(databaseURL string) (string, error) {
	// Parse the URL
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse DATABASE_URL: %v", err)
	}

	// Extract components
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "5432" // default PostgreSQL port
	}

	dbname := strings.TrimPrefix(u.Path, "/")
	user := u.User.Username()
	password, _ := u.User.Password()

	// Determine SSL mode
	sslmode := "require" // Railway requires SSL
	if u.Query().Get("sslmode") != "" {
		sslmode = u.Query().Get("sslmode")
	}

	// Build lib/pq connection string
	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		user, password, dbname, host, port, sslmode,
	)

	return connectionString, nil
}

// sanitizeConnectionString removes password from connection string for logging
func sanitizeConnectionString(connectionString string) string {
	// Hide password in logs for security
	if strings.Contains(connectionString, "password=") {
		parts := strings.Split(connectionString, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, "password=") {
				parts[i] = "password=****"
			}
		}
		return strings.Join(parts, " ")
	}
	return connectionString
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
	case *events.Disconnected:
		fmt.Println("WhatsApp client disconnected!")
	}
}

func (wa *WhatsAppService) GenerateQR() (string, error) {
	if wa.Container == nil {
		return "", fmt.Errorf("container is not initialized")
	}

	deviceStore, err := wa.Container.GetFirstDevice()
	if err != nil {
		return "", fmt.Errorf("failed to get first device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	wa.Client = whatsmeow.NewClient(deviceStore, clientLog)
	wa.Client.AddEventHandler(eventHandler)

	if wa.Client.Store.ID == nil {
		fmt.Println("Client is not logged in, generating QR code...")
		qrChan, _ := wa.Client.GetQRChannel(context.Background())
		err = wa.Client.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR code generated:", evt.Code)
				png, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					return "", fmt.Errorf("failed to create QR code: %v", err)
				}
				dataURI := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
				return dataURI, nil
			} else {
				fmt.Println("Login event:", evt.Event)
				if evt.Event == "success" {
					return "success", nil
				}
			}
		}
	} else {
		fmt.Println("Client already logged in, connecting...")
		err = wa.Client.Connect()
		if err != nil {
			return "", fmt.Errorf("failed to connect: %v", err)
		}
		return "already_connected", nil
	}

	return "", fmt.Errorf("failed to generate QR code")
}

func (wa *WhatsAppService) SendMessage(phoneNumber, message string) error {
	if wa.Client == nil {
		return fmt.Errorf("client not initialized")
	}

	targetJID, err := types.ParseJID(phoneNumber + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid phone number: %v", err)
	}

	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	_, err = wa.Client.SendMessage(context.Background(), targetJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (wa *WhatsAppService) Logout() error {
	if wa.Client == nil {
		return fmt.Errorf("no active client session")
	}

	err := wa.Client.Logout()
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
