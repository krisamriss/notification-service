package services

import (
	"context"
	"errors"
	"fmt"
	"notification-service/internal/core/models"
	"notification-service/internal/core/ports"
	"sync"
	"time"
)

type NotificationService struct {
	notifiers []ports.Notifier
	templates ports.TemplateEngine
	scheduler ports.Scheduler
	mu        sync.RWMutex
}

func NewNotificationService(n []ports.Notifier, t ports.TemplateEngine, s ports.Scheduler) *NotificationService {
	return &NotificationService{
		notifiers: n,
		templates: t,
		scheduler: s,
	}
}

func (s *NotificationService) SetScheduler(sch ports.Scheduler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scheduler = sch
}

func (s *NotificationService) Process(ctx context.Context, req models.NotificationRequest) error {
	if req.ScheduledAt != nil && req.ScheduledAt.After(time.Now()) {
		s.mu.RLock()
		sch := s.scheduler
		s.mu.RUnlock()
		if sch != nil {
			return sch.Schedule(ctx, req)
		}
	}

	message, err := s.templates.Render(req.TemplateName, req.CustomBody, req.Data)
	if err != nil {
		return fmt.Errorf("template rendering failed: %v", err)
	}

	for _, notifier := range s.notifiers {
		if notifier.Supports() == req.Channel {
			return notifier.Send(ctx, req.UserID, message)
		}
	}

	return errors.New("unsupported notification channel")
}

func (s *NotificationService) Broadcast(ctx context.Context, channels []models.ChannelType, req models.NotificationRequest) map[models.ChannelType]error {
	results := make(map[models.ChannelType]error, len(channels))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(ch models.ChannelType) {
			defer wg.Done()
			r := req
			r.Channel = ch
			err := s.Process(ctx, r)
			mu.Lock()
			results[ch] = err
			mu.Unlock()
		}(ch)
	}
	wg.Wait()
	return results
}
