package cron

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

func NewScheduler() *gocron.Scheduler {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Hours().Do(moveToTrashTask)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
