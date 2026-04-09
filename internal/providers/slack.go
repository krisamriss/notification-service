package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"notification-service/internal/core/models"
)

type SlackNotifier struct {
	webhookURL string // Slack webhook
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
	}
}

func (s *SlackNotifier) Send(userID string, message string) error {

	payload := map[string]string{
		"text": fmt.Sprintf("<@%s> %s", userID, message),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %v", err)
	}
	
	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send slack message: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status: %d", resp.StatusCode)
	}
	

	fmt.Printf("[SLACK] sending payload to webhook: %s\n", string(jsonData))
	return nil
}

func (s *SlackNotifier) Supports() models.ChannelType {
	return models.Slack
}