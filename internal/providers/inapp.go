package providers

import (
	"fmt"
	"notification-service/internal/core/models"
	"time"
)

// InAppNotifier implements the Notifier interface for In-App (Web/Mobile) notifications
type InAppNotifier struct {
	// In production, inject your Database connection or WebSocket client here
	// db *sql.DB 
}

func NewInAppNotifier() *InAppNotifier {
	return &InAppNotifier{}
}

func (i *InAppNotifier) Send(userID string, message string) error {
	// PRO-TIP: In production, this would be an SQL/NoSQL insert query:
	/*
		query := `INSERT INTO inapp_notifications (user_id, message, is_read, created_at) VALUES ($1, $2, $3, $4)`
		_, err := i.db.Exec(query, userID, message, false, time.Now())
		if err != nil {
			return fmt.Errorf("failed to save in-app notification to DB: %v", err)
		}
		
		// (Optional) Push to WebSocket channel if user is currently online
		// websocketServer.BroadcastToUser(userID, message)
	*/

	// Mocking the DB insert for this assignment
	fmt.Printf("[IN-APP] Simulated saving notification to Database for user '%s' at %v. Message: %s\n", 
		userID, time.Now().Format(time.RFC3339), message)
	return nil
}

func (i *InAppNotifier) Supports() models.ChannelType {
	return models.InApp
}