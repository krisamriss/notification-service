package providers

import (
	"context"
	"database/sql"
	"fmt"
	"notification-service/internal/core/models"
	"time"
)

type InAppNotifier struct {
	db *sql.DB
}

func NewInAppNotifier(db *sql.DB) *InAppNotifier {
	return &InAppNotifier{db: db}
}

func (i *InAppNotifier) Send(ctx context.Context, userID string, message string) error {
	if i.db == nil {
		fmt.Printf("[IN-APP] (mock) notification for user '%s' at %s — %s\n",
			userID, time.Now().Format(time.RFC3339), message)
		return nil
	}

	query := `
		INSERT INTO inapp_notifications (user_id, message, is_read, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := i.db.ExecContext(ctx, query, userID, message, false, time.Now())
	if err != nil {
		return fmt.Errorf("[IN-APP] failed to save notification for user %s: %w", userID, err)
	}

	fmt.Printf("[IN-APP] Saved notification for user '%s': %s\n", userID, message)
	return nil
}

func (i *InAppNotifier) Supports() models.ChannelType {
	return models.InApp
}
