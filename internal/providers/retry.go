package providers

import (
	"context"
	"fmt"
	"notification-service/internal/core/models"
	"notification-service/internal/core/ports"
	"time"
)

type RetryNotifier struct {
	inner      ports.Notifier
	maxRetries int
	baseDelay  time.Duration
}

func NewRetryNotifier(inner ports.Notifier, maxRetries int, baseDelay time.Duration) *RetryNotifier {
	return &RetryNotifier{inner: inner, maxRetries: maxRetries, baseDelay: baseDelay}
}

func (r *RetryNotifier) Send(ctx context.Context, userID, message string) error {
	var err error
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if err = r.inner.Send(ctx, userID, message); err == nil {
			return nil
		}
		if attempt == r.maxRetries {
			break
		}
		delay := r.baseDelay * (1 << attempt)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("after %d attempts: %w", r.maxRetries+1, err)
}

func (r *RetryNotifier) Supports() models.ChannelType {
	return r.inner.Supports()
}
