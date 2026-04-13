package ports

import (
	"context"
	"notification-service/internal/core/models"
)

type Notifier interface {
	Send(ctx context.Context, userID string, message string) error
	Supports() models.ChannelType
}

type TemplateEngine interface {
	Render(templateName string, customBody string, data map[string]interface{}) (string, error)
}

type Scheduler interface {
	Schedule(ctx context.Context, req models.NotificationRequest) error
}
