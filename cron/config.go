package cron

import (
	"github.com/go-co-op/gocron"
	"github.com/wittano/filebot/setting"
	"time"
)

func NewScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Hours().Do(moveToTrashTask)
	if err != nil {
		setting.Logger().Fatal("Failed to register 'moveToTrash' task", err)
	}

	return s
}
