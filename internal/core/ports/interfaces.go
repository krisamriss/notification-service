package ports

import "notification-service/internal/core/models"

type Notifier interface {
	Send(userID string, message string) error
	Supports() models.ChannelType
}

type TemplateEngine interface {
	Render(templateName string, customBody string, data map[string]interface{}) (string, error)
}

// Scheduler handles async and delayed tasks
type Scheduler interface {
	Schedule(req models.NotificationRequest) error
}
