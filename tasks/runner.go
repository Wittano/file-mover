package tasks

import (
	"context"
	"github.com/wittano/filebot/setting"
	"time"
)

type taskRunner func(ctx context.Context) error

func Run(ctx context.Context, interval time.Duration, task taskRunner) {
	var err error

	newCtx, cancel := context.WithCancelCause(ctx)
	defer cancel(err)

	go func() {
		timer := time.NewTicker(interval)
		defer timer.Stop()

		for {
			select {
			case <-newCtx.Done():
				return
			case <-timer.C:
				if err = task(newCtx); err != nil {
					setting.Logger().Error("Error during execute scheduled task", err)
					return
				}
			}
		}
	}()
}
