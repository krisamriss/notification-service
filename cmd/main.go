package main

import (
	"context"
	"fmt"
	"notification-service/internal/core/models"
	"notification-service/internal/core/ports"
	"notification-service/internal/dispatcher"
	"notification-service/internal/providers"
	schedpkg "notification-service/internal/scheduler"
	"notification-service/internal/services"
	"notification-service/internal/templates"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	notifiers := []ports.Notifier{
		providers.NewRetryNotifier(
			providers.NewEmailNotifier(
				os.Getenv("SMTP_HOST"),
				os.Getenv("SMTP_PORT"),
				os.Getenv("SMTP_USERNAME"),
				os.Getenv("SMTP_PASSWORD"),
				os.Getenv("SMTP_FROM"),
			),
			3, time.Second,
		),
		providers.NewRetryNotifier(
			providers.NewSlackNotifier(os.Getenv("SLACK_WEBHOOK_URL")),
			3, time.Second,
		),
		providers.NewInAppNotifier(nil),
	}

	templateEngine := templates.NewGoTemplateEngine()

	svc := services.NewNotificationService(notifiers, templateEngine, nil)

	d := dispatcher.New(runtime.NumCPU()*4, 500, svc.Process)

	sch := schedpkg.New(func(ctx context.Context, req models.NotificationRequest) error {
		req.ScheduledAt = nil
		return d.Submit(ctx, req)
	})
	svc.SetScheduler(sch)

	go sch.Run(ctx)

	if err := d.Submit(ctx, models.NotificationRequest{
		UserID:     "user_999",
		Channel:    models.InApp,
		CustomBody: "Your profile has been updated successfully.",
	}); err != nil {
		fmt.Printf("submit error: %v\n", err)
	}

	if err := d.Submit(ctx, models.NotificationRequest{
		UserID:       "U123456",
		Channel:      models.Slack,
		TemplateName: "alert",
		Data:         map[string]interface{}{"message": "Server CPU > 90%", "level": "CRITICAL"},
	}); err != nil {
		fmt.Printf("submit error: %v\n", err)
	}

	scheduledAt := time.Now().Add(2 * time.Second)
	if err := d.Submit(ctx, models.NotificationRequest{
		UserID:      "user_123",
		Channel:     models.InApp,
		CustomBody:  "This is a scheduled notification.",
		ScheduledAt: &scheduledAt,
	}); err != nil {
		fmt.Printf("submit error: %v\n", err)
	}

	go func() {
		results := svc.Broadcast(ctx, []models.ChannelType{models.InApp, models.Slack},
			models.NotificationRequest{
				UserID:     "user_456",
				CustomBody: "System maintenance in 10 minutes.",
			})
		for ch, err := range results {
			if err != nil {
				fmt.Printf("[BROADCAST] %s failed: %v\n", ch, err)
			}
		}
	}()

	time.Sleep(4 * time.Second)

	<-ctx.Done()
	fmt.Println("\nShutting down — draining worker pool...")
	d.Shutdown()
	fmt.Println("Shutdown complete.")
}
