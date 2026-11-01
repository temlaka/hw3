package schedule

import (
	"hw3/internal/domain"
	"time"
)

type PriceUpdater interface {
	UpdateCrypto(symbol string) (domain.Crypto, error)
	GetListCrypto() ([]domain.Crypto, error)
}
type ScheduleRepository interface {
	Get() domain.Schedule
	Edit(enabled bool, interval int, lastUpdated time.Time, nextUpdate time.Time) error
	Update() error
}
