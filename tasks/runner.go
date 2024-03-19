package tasks

import (
	"context"
	"time"
)

type taskRunner func(ctx context.Context) error

func RunTaskWithInterval(ctx context.Context, interval time.Duration, task taskRunner) {
	var err error

	newCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(err)

	go func() {
		timer := time.NewTicker(interval)
		defer timer.Stop()

		for {
			select {
			case <-newCtx.Done():
				break
			case <-timer.C:
				if err = task(newCtx); err != nil {
					break
				}
			}
		}
	}()
}
