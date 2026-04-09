package main

import (
	"notification-service/internal/core/models"
	"notification-service/internal/core/ports"
	"notification-service/internal/providers"
	"notification-service/internal/services"
	"notification-service/internal/templates"
	"os"
	"time"
)

// For production we can use any pub-sub
type MockScheduler struct{}

func (m *MockScheduler) Schedule(req models.NotificationRequest) error {
	/*
		scheduler will push the message to the queue . 
		for now i am just doing time.sleep 
		to assume the message is pushed to queue and then the messsage is sent
		to the user
	*/
	go func() {
		delay := time.Until(*req.ScheduledAt)
		time.Sleep(delay)
		req.ScheduledAt = nil

	}()
	return nil
}

func main() {
	/* 
		Initialize Providers (Channels)
	*/
	emailNotifier := providers.NewEmailNotifier(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("SMTP_FROM"),
	)
	slackNotifier := providers.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL"))
	inAppNotifier := providers.NewInAppNotifier(nil)

	notifiers := []ports.Notifier{
		emailNotifier,
		slackNotifier,
		inAppNotifier,
		
	}

	templateEngine := templates.NewGoTemplateEngine()
	scheduler := &MockScheduler{}

	svc := services.NewNotificationService(notifiers, templateEngine, scheduler)

	/* 
		i could have the test.go for the test cases running
		but for the simplicity i am putting the test function running in the main.go
	*/
	reqInApp := models.NotificationRequest{
		UserID:     "user_999",
		Channel:    models.InApp,
		CustomBody: "Your profile has been updated successfully.",
	}
	svc.Process(reqInApp)

	reqSlack := models.NotificationRequest{
		UserID:       "U123456", // Slack User ID
		Channel:      models.Slack,
		TemplateName: "alert",
		Data:         map[string]interface{}{"message": "Server CPU > 90%", "level": "CRITICAL"},
	}
	svc.Process(reqSlack)
}
