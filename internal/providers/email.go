package providers

import (
	"fmt"
	"notification-service/internal/core/models"
)

type EmailNotifier struct {
	// Add SMTP credentials/clients here
}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (e *EmailNotifier) Send(userID string, message string) error {
	// TODO: Add actual SendGrid / SMTP logic here
	fmt.Printf("[EMAIL] Sending to %s: %s\n", userID, message)
	return nil
}

func (e *EmailNotifier) Supports() models.ChannelType {
	return models.Email
}
