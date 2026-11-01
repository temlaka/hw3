package schedule

import (
	"fmt"
	"hw3/internal/domain"
	"log"
	"time"
)

const (
	minIntervalSeconds = 10
	maxIntervalSeconds = 3600
)

type Service struct {
	schedulerepo ScheduleRepository
	priceUpd     PriceUpdater
}

func NewService(srepo ScheduleRepository, crypo PriceUpdater) *Service {
	return &Service{
		schedulerepo: srepo,
		priceUpd:     crypo,
	}
}

func (s *Service) Get() domain.Schedule {

	return s.schedulerepo.Get()
}

func (s *Service) EditSettings(enabled bool, interval int) error {
	if interval < minIntervalSeconds || interval > maxIntervalSeconds {
		return fmt.Errorf("interval %d out of range %d..%d: %w",
			interval, minIntervalSeconds, maxIntervalSeconds, domain.ErrInvalidInput)
	}

	x := s.Get()
	if err := s.schedulerepo.Edit(enabled, interval, x.LastUpdate, x.NextUpdate); err != nil {
		return fmt.Errorf("edit schedule: %w", err)
	}
	return nil
}

func (s *Service) Trigger() (count int, timestamp time.Time, err error) {
	list, err := s.priceUpd.GetListCrypto()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("list cryptos: %w", err)
	}
	count = 0
	for _, item := range list {
		_, err = s.priceUpd.UpdateCrypto(item.Symbol)
		if err != nil {
			log.Printf("trigger: update %s: %v", item.Symbol, err)
			continue
		}
		count += 1
	}
	if err = s.schedulerepo.Update(); err != nil {
		return count, time.Time{}, fmt.Errorf("save schedule: %w", err)
	}

	return count, time.Now().UTC(), nil
}
