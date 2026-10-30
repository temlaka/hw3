package crypto

import "hw3/internal/domain"

type CryptoRepository interface {
	Create(c domain.Crypto) error
	List() ([]domain.Crypto, error)
	Get(symbol string) (domain.Crypto, error)
	UpdatePrice(symbol string, p domain.PricePoint) error
	History(symbol string) ([]domain.PricePoint, error)
	Delete(symbol string) error
}
