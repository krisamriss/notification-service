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
	// check if  Scheduling
	if req.ScheduledAt != nil && req.ScheduledAt.After(time.Now()) {
		
		return s.scheduler.Schedule(req)
	}

	// here i am loading the  template
	message, err := s.templates.Render(req.TemplateName, req.CustomBody, req.Data)
	if err != nil {
		return fmt.Errorf("template rendering failed: %v", err)
	}

	// channel Mapping
	for _, notifier := range s.notifiers {
		if notifier.Supports() == req.Channel {
			// sending
			return notifier.Send(req.UserID, message)
		}
	}

	return errors.New("unsupported notification channel")
}
