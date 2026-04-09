package providers

import (
	"fmt"
	"notification-service/internal/core/models"
)

// SMSNotifier implements the Notifier interface for text messages
type SMSNotifier struct {
	accountSID string
	authToken  string
	fromNumber string
}

// NewSMSNotifier initializes the SMS client (e.g., Twilio credentials)
func NewSMSNotifier(accountSID, authToken, fromNumber string) *SMSNotifier {
	return &SMSNotifier{
		accountSID: accountSID,
		authToken:  authToken,
		fromNumber: fromNumber,
	}
}

func (s *SMSNotifier) Send(userID string, message string) error {
	// PRO-TIP: In production, you would use a Twilio SDK or HTTP client here.
	// Example:
	/*
		// 1. Fetch user's phone number from DB using userID
		userPhone := getUserPhoneNumber(userID)

		// 2. Make API call to SMS Gateway (e.g., Twilio)
		client := twilio.NewRestClient(s.accountSID, s.authToken)
		params := &openapi.CreateMessageParams{}
		params.SetTo(userPhone)
		params.SetFrom(s.fromNumber)
		params.SetBody(message)

		_, err := client.Api.V2010.CreateMessage(params)
		if err != nil {
			return fmt.Errorf("failed to send SMS via Twilio: %v", err)
		}
	*/

	// Mocking the successful send for the assignment
	fmt.Printf("[SMS] Simulated sending SMS from %s to user '%s'. Message: %s\n", s.fromNumber, userID, message)
	return nil
}

func (s *SMSNotifier) Supports() models.ChannelType {
	return models.SMS
}
