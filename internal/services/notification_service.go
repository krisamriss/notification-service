package services

import (
	"errors"
	"fmt"
	"notification-service/internal/core/models"
	"notification-service/internal/core/ports"
	"time"
)

type NotificationService struct {
	notifiers []ports.Notifier // Holds Email, Slack, InApp providers
	templates ports.TemplateEngine
	scheduler ports.Scheduler
}

func NewNotificationService(n []ports.Notifier, t ports.TemplateEngine, s ports.Scheduler) *NotificationService {
	return &NotificationService{
		notifiers: n,
		templates: t,
		scheduler: s,
	}
}

func (s *NotificationService) Process(req models.NotificationRequest) error {
	// 1. Check Scheduling
	if req.ScheduledAt != nil && req.ScheduledAt.After(time.Now()) {
		fmt.Println("Delegating to Scheduler...")
		return s.scheduler.Schedule(req)
	}

	// 2. Render Template
	message, err := s.templates.Render(req.TemplateName, req.CustomBody, req.Data)
	if err != nil {
		return fmt.Errorf("template rendering failed: %v", err)
	}

	// 3. Channel Mapping (Strategy Pattern)
	for _, notifier := range s.notifiers {
		if notifier.Supports() == req.Channel {
			// Send instantly
			return notifier.Send(req.UserID, message)
		}
	}

	return errors.New("unsupported notification channel")
}
