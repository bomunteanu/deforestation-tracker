package jobs

import (
	"log"

	"github.com/robfig/cron/v3"
)

var (
	jobCron *cron.Cron
)

func init() {
	jobCron = cron.New()
	jobCron.Start()
}

func StartWeeklyJob(areaID uint, getImage func(uint) error) {
	job := func() {
		getImage(areaID)
	}

	schedule := "0 0 * * 0" // Every Sunday at midnight
	_, err := jobCron.AddFunc(schedule, job)
	if err != nil {
		log.Fatalf("Error scheduling job: %v", err)
	}
}

func StopAllJobs() {
	jobCron.Stop()
}
