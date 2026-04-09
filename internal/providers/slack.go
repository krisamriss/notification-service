package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notification-service/internal/core/models"
)

// SlackNotifier implements the Notifier interface for Slack
type SlackNotifier struct {
	webhookURL string // Slack webhook URL loaded from config/env
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Send(userID string, message string) error {
	// Slack expects a specific JSON payload format
	payload := map[string]string{
		"text": fmt.Sprintf("<@%s> %s", userID, message), // Mentions the user
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %v", err)
	}

	// PRO-TIP: In production, you would actually make this HTTP call.
	// Uncomment the below code in a real scenario:
	
	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send slack message: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status: %d", resp.StatusCode)
	}
	

	// Mocking the successful send for this assignment
	fmt.Printf("[SLACK] Simulated sending payload to webhook: %s\n", string(jsonData))
	return nil
}

func (s *SlackNotifier) Supports() models.ChannelType {
	return models.Slack
}