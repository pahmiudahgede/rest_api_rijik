package config

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
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

	connectionString := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("postgres", connectionString, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	whatsappService = &WhatsAppService{
		Container: container,
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
