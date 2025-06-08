package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsAppManager struct {
	Client      *whatsmeow.Client
	container   *sqlstore.Container
	isConnected bool
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	shutdownCh  chan struct{}
}

var (
	waManager *WhatsAppManager
	once      sync.Once
)

func GetWhatsAppManager() *WhatsAppManager {
	once.Do(func() {
		ctx, cancel := context.WithCancel(context.Background())
		waManager = &WhatsAppManager{
			ctx:        ctx,
			cancel:     cancel,
			shutdownCh: make(chan struct{}),
		}
	})
	return waManager
}

func InitWhatsApp() {
	manager := GetWhatsAppManager()

	log.Println("Initializing WhatsApp client...")

	if err := manager.setupDatabase(); err != nil {
		log.Fatalf("Failed to setup WhatsApp database: %v", err)
	}

	if err := manager.setupClient(); err != nil {
		log.Fatalf("Failed to setup WhatsApp client: %v", err)
	}

	if err := manager.handleAuthentication(); err != nil {
		log.Fatalf("Failed to authenticate WhatsApp: %v", err)
	}

	manager.setupEventHandlers()

	go manager.handleShutdown()

	log.Println("WhatsApp client initialized successfully and ready to send messages!")
}

func (w *WhatsAppManager) setupDatabase() error {
	dbLog := waLog.Stdout("WhatsApp-DB", "ERROR", true)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	w.container, err = sqlstore.New("postgres", dsn, dbLog)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("WhatsApp database connection established")
	return nil
}

func (w *WhatsAppManager) setupClient() error {
	deviceStore, err := w.container.GetFirstDevice()
	if err != nil {
		return fmt.Errorf("failed to get device store: %v", err)
	}

	clientLog := waLog.Stdout("WhatsApp-Client", "ERROR", true)
	w.Client = whatsmeow.NewClient(deviceStore, clientLog)

	return nil
}

func (w *WhatsAppManager) handleAuthentication() error {
	if w.Client.Store.ID == nil {
		log.Println("WhatsApp client not logged in, generating QR code...")
		return w.authenticateWithQR()
	}

	log.Println("WhatsApp client already logged in, connecting...")
	return w.connect()
}

func (w *WhatsAppManager) authenticateWithQR() error {
	qrChan, err := w.Client.GetQRChannel(w.ctx)
	if err != nil {
		return fmt.Errorf("failed to get QR channel: %v", err)
	}

	if err := w.Client.Connect(); err != nil {
		return fmt.Errorf("failed to connect client: %v", err)
	}

	qrTimeout := time.NewTimer(3 * time.Minute)
	defer qrTimeout.Stop()

	for {
		select {
		case evt := <-qrChan:
			switch evt.Event {
			case "code":
				fmt.Println("\n=== QR CODE UNTUK LOGIN WHATSAPP ===")
				generateQRCode(evt.Code)
				fmt.Println("Scan QR code di atas dengan WhatsApp Anda")
				fmt.Println("QR code akan expired dalam 3 menit")
			case "success":
				log.Println("✅ WhatsApp login successful!")
				w.setConnected(true)
				return nil
			case "timeout":
				return fmt.Errorf("QR code expired, please restart")
			default:
				log.Printf("Login status: %s", evt.Event)
			}
		case <-qrTimeout.C:
			return fmt.Errorf("QR code authentication timeout after 3 minutes")
		case <-w.ctx.Done():
			return fmt.Errorf("authentication cancelled")
		}
	}
}

func (w *WhatsAppManager) connect() error {
	if err := w.Client.Connect(); err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	time.Sleep(2 * time.Second)
	w.setConnected(true)
	return nil
}

func (w *WhatsAppManager) setupEventHandlers() {
	w.Client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Connected:
			log.Println("✅ WhatsApp client connected")
			w.setConnected(true)
		case *events.Disconnected:
			log.Println("❌ WhatsApp client disconnected")
			w.setConnected(false)
		case *events.LoggedOut:
			log.Println("🚪 WhatsApp client logged out")
			w.setConnected(false)
		case *events.Message:
			log.Printf("📨 Message received from %s", v.Info.Sender)
		}
	})
}

func (w *WhatsAppManager) setConnected(status bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isConnected = status
}

func (w *WhatsAppManager) IsConnected() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isConnected
}

func generateQRCode(qrString string) {
	qrterminal.GenerateHalfBlock(qrString, qrterminal.M, os.Stdout)
}

func (w *WhatsAppManager) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		log.Println("Received shutdown signal...")
	case <-w.ctx.Done():
		log.Println("Context cancelled...")
	}

	w.shutdown()
}

func (w *WhatsAppManager) shutdown() {
	log.Println("Shutting down WhatsApp client...")

	w.cancel()

	if w.Client != nil {
		w.Client.Disconnect()
	}

	if w.container != nil {
		w.container.Close()
	}

	close(w.shutdownCh)
	log.Println("WhatsApp client shutdown completed")
}

func SendWhatsAppMessage(phone, message string) error {
	manager := GetWhatsAppManager()

	if manager.Client == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	if !manager.IsConnected() {
		return fmt.Errorf("WhatsApp client is not connected")
	}

	if phone == "" || message == "" {
		return fmt.Errorf("phone number and message cannot be empty")
	}

	if phone[0] == '0' {
		phone = "62" + phone[1:] // Convert 08xx menjadi 628xx
	}
	if phone[:2] != "62" {
		phone = "62" + phone // Tambahkan 62 jika belum ada
	}

	// Parse JID
	targetJID, err := types.ParseJID(phone + "@s.whatsapp.net")
	if err != nil {
		return fmt.Errorf("invalid phone number format: %v", err)
	}

	// Buat pesan
	msg := &waE2E.Message{
		Conversation: proto.String(message),
	}

	// Kirim dengan timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := manager.Client.SendMessage(ctx, targetJID, msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("✅ Message sent to %s (ID: %s)", phone, resp.ID)
	return nil
}

// SendWhatsAppMessageBatch - Kirim pesan ke multiple nomor
func SendWhatsAppMessageBatch(phoneNumbers []string, message string) []error {
	var errors []error

	for _, phone := range phoneNumbers {
		if err := SendWhatsAppMessage(phone, message); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %v", phone, err))
			continue
		}

		// Delay untuk menghindari rate limit
		time.Sleep(1 * time.Second)
	}

	return errors
}

// GetWhatsAppStatus - Cek status koneksi
func GetWhatsAppStatus() map[string]interface{} {
	manager := GetWhatsAppManager()

	status := map[string]interface{}{
		"initialized": manager.Client != nil,
		"connected":   manager.IsConnected(),
		"logged_in":   false,
		"jid":         "",
	}

	if manager.Client != nil && manager.Client.Store.ID != nil {
		status["logged_in"] = true
		status["jid"] = manager.Client.Store.ID.String()
	}

	return status
}

// LogoutWhatsApp - Logout dan cleanup
func LogoutWhatsApp() error {
	manager := GetWhatsAppManager()

	if manager.Client == nil {
		return fmt.Errorf("WhatsApp client is not initialized")
	}

	log.Println("Logging out WhatsApp...")

	// Logout
	err := manager.Client.Logout()
	if err != nil {
		log.Printf("Warning: Failed to logout properly: %v", err)
	}

	// Disconnect
	manager.Client.Disconnect()
	manager.setConnected(false)

	// Hapus device dari store
	if err := manager.removeDeviceFromStore(); err != nil {
		log.Printf("Warning: Failed to remove device: %v", err)
	}

	// Close database
	if manager.container != nil {
		manager.container.Close()
	}

	log.Println("✅ WhatsApp logout completed")
	return nil
}

func (w *WhatsAppManager) removeDeviceFromStore() error {
	deviceStore, err := w.container.GetFirstDevice()
	if err != nil {
		return err
	}

	if deviceStore != nil && deviceStore.ID != nil {
		return deviceStore.Delete()
	}

	return nil
}

// IsValidPhoneNumber - Validasi format nomor telepon Indonesia
func IsValidPhoneNumber(phone string) bool {
	// Minimal validasi untuk nomor Indonesia
	if len(phone) < 10 || len(phone) > 15 {
		return false
	}

	// Cek awalan nomor Indonesia
	if phone[:2] == "62" || phone[0] == '0' {
		return true
	}

	return false
}
