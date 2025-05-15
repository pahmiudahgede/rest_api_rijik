package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

var WhatsAppClient *whatsmeow.Client
var container *sqlstore.Container

func InitWhatsApp() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	container, err = sqlstore.New("postgres", dsn, dbLog)
	if err != nil {
		log.Fatalf("Failed to connect to WhatsApp database: %v", err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Failed to get WhatsApp device: %v", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	WhatsAppClient = whatsmeow.NewClient(deviceStore, clientLog)

	if WhatsAppClient.Store.ID == nil {
		fmt.Println("WhatsApp Client is not logged in, generating QR Code...")

		qrChan, _ := WhatsAppClient.GetQRChannel(context.Background())
		err = WhatsAppClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect WhatsApp client: %v", err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR Code untuk login:")
				generateQRCode(evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		fmt.Println("WhatsApp Client sudah login, langsung terhubung...")
		err = WhatsAppClient.Connect()
		if err != nil {
			log.Fatalf("Failed to connect WhatsApp client: %v", err)
		}
	}

	log.Println("WhatsApp client connected successfully!")
	go handleShutdown()
}

func generateQRCode(qrString string) {
	qrterminal.GenerateHalfBlock(qrString, qrterminal.M, os.Stdout)
}

func handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down WhatsApp client...")
	WhatsAppClient.Disconnect()
	os.Exit(0)
}

func SendWhatsAppMessage(phone, message string) error {
	if WhatsAppClient == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	targetJID, _ := types.ParseJID(phone + "@s.whatsapp.net")
	msg := waE2E.Message{
		Conversation: proto.String(message),
	}

	_, err := WhatsAppClient.SendMessage(context.Background(), targetJID, &msg)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %v", err)
	}

	log.Printf("WhatsApp message sent successfully to: %s", phone)
	return nil
}

func LogoutWhatsApp() error {
	if WhatsAppClient == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	WhatsAppClient.Disconnect()

	err := removeWhatsAppDeviceFromContainer()
	if err != nil {
		return fmt.Errorf("failed to remove device from container: %v", err)
	}

	err = container.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %v", err)
	}

	log.Println("WhatsApp client disconnected and session cleared successfully.")
	return nil
}

func removeWhatsAppDeviceFromContainer() error {
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return fmt.Errorf("failed to get WhatsApp device: %v", err)
	}

	if deviceStore != nil {
		err := deviceStore.Delete()
		if err != nil {
			return fmt.Errorf("failed to remove device from store: %v", err)
		}
	}

	return nil
}
