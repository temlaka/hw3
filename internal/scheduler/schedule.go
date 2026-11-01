package scheduler

import (
	"context"
	"hw3/internal/usecase/schedule"
	"log"
	"time"
)

func UpdateAllPrice(ctx context.Context, svc *schedule.Service) {
	for {

		settings := svc.Get()
		select {
		case <-ctx.Done():
			log.Println("schedule stopped")
			return
		case <-time.After(time.Duration(settings.IntervalSeconds) * time.Second):

			settings := svc.Get()
			if !settings.Enabled {
				continue
			}
			_, _, err := svc.Trigger()
			if err != nil {
				log.Printf("trigger: %v", err)
			}
		}
	}
}
