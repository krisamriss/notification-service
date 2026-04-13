package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"notification-service/internal/core/models"
	"sync"
)

type ProcessFunc func(ctx context.Context, req models.NotificationRequest) error

type job struct {
	ctx context.Context
	req models.NotificationRequest
}

type Dispatcher struct {
	jobs    chan job
	wg      sync.WaitGroup
	process ProcessFunc
}

func New(workers, queueSize int, process ProcessFunc) *Dispatcher {
	d := &Dispatcher{
		jobs:    make(chan job, queueSize),
		process: process,
	}
	for i := 0; i < workers; i++ {
		d.wg.Add(1)
		go func() {
			defer d.wg.Done()
			for j := range d.jobs {
				if err := d.process(j.ctx, j.req); err != nil {
					fmt.Printf("[DISPATCHER] failed for user %s on channel %s: %v\n",
						j.req.UserID, j.req.Channel, err)
				}
			}
		}()
	}
	return d
}

func (d *Dispatcher) Submit(ctx context.Context, req models.NotificationRequest) error {
	select {
	case d.jobs <- job{ctx: ctx, req: req}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.New("dispatcher: queue full, notification dropped")
	}
}

func (d *Dispatcher) Shutdown() {
	close(d.jobs)
	d.wg.Wait()
}
