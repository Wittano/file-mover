package tasks

import (
	"context"
	"time"
)

func RunTaskWithInterval(ctx context.Context, interval time.Duration, task func(cancel context.CancelFunc)) {
	_, cancel := context.WithCancel(ctx)

	go func() {
		timer := time.NewTicker(interval)
		defer timer.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				task(cancel)
			}
		}
	}()
}
