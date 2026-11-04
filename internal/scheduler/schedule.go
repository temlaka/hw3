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

		if !settings.Enabled {
			select {
			case <-ctx.Done():
				log.Println("schedule stopped")
				return

			case <-svc.Changes():
				continue
			}
		}
		select {
		case <-ctx.Done():
			log.Println("schedule stopped")
			return

		case <-svc.Changes():
			continue
		case <-time.After(time.Until(settings.NextUpdate)):
			_, _, err := svc.Trigger()
			if err != nil {
				log.Printf("trigger: %v", err)
			}
		}
	}
}
