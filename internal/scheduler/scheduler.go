package scheduler

import (
	"container/heap"
	"context"
	"notification-service/internal/core/models"
	"sync"
	"time"
)

type FireFunc func(ctx context.Context, req models.NotificationRequest) error

type notificationHeap []models.NotificationRequest

func (h notificationHeap) Len() int { return len(h) }
func (h notificationHeap) Less(i, j int) bool {
	return (*h[i].ScheduledAt).Before(*h[j].ScheduledAt)
}
func (h notificationHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *notificationHeap) Push(x interface{}) {
	*h = append(*h, x.(models.NotificationRequest))
}
func (h *notificationHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

type HeapScheduler struct {
	mu      sync.Mutex
	h       notificationHeap
	trigger chan struct{}
	fire    FireFunc
}

func New(fire FireFunc) *HeapScheduler {
	s := &HeapScheduler{
		h:       make(notificationHeap, 0),
		trigger: make(chan struct{}, 1),
		fire:    fire,
	}
	heap.Init(&s.h)
	return s
}

func (s *HeapScheduler) Schedule(ctx context.Context, req models.NotificationRequest) error {
	s.mu.Lock()
	heap.Push(&s.h, req)
	s.mu.Unlock()

	select {
	case s.trigger <- struct{}{}:
	default:
	}
	return nil
}

func (s *HeapScheduler) Run(ctx context.Context) {
	for {
		s.mu.Lock()
		var waitDur time.Duration
		if s.h.Len() > 0 {
			waitDur = time.Until(*s.h[0].ScheduledAt)
			if waitDur < 0 {
				waitDur = 0
			}
		} else {
			waitDur = time.Hour
		}
		s.mu.Unlock()

		timer := time.NewTimer(waitDur)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-s.trigger:
			timer.Stop()
		case <-timer.C:
			s.mu.Lock()
			for s.h.Len() > 0 && !time.Now().Before(*s.h[0].ScheduledAt) {
				req := heap.Pop(&s.h).(models.NotificationRequest)
				go s.fire(ctx, req)
			}
			s.mu.Unlock()
		}
	}
}
