package ports

import "notification-service/internal/core/models"

// Notifier is the Strategy interface for different channels
type Notifier interface {
	Send(userID string, message string) error
	Supports() models.ChannelType
}

// TemplateEngine parses data into templates
type TemplateEngine interface {
	Render(templateName string, customBody string, data map[string]interface{}) (string, error)
}

// Scheduler handles async and delayed tasks
type Scheduler interface {
	Schedule(req models.NotificationRequest) error
}
