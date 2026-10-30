package crypto

import (
	"hw3/internal/domain"
	"slices"
	"strings"
	"time"
)

type Service struct {
	repo CryptoRepository
	cg   PriceProvider
}

func NewService(repo CryptoRepository, cg PriceProvider) *Service {
	return &Service{
		repo: repo,
		cg:   cg,
	}
}

func (s *Service) GetListCrypto() ([]domain.Crypto, error) {
	list, err := s.repo.List()
	if err != nil {
		return []domain.Crypto{}, err
	}
	return list, nil
}

func (s *Service) AddCrypto(symbol string) (domain.Crypto, error) {
	symbol = strings.ToUpper(symbol)
	name, price, err := s.cg.GetPrice(symbol)
	if err != nil {
		return domain.Crypto{}, err
	}

	newCrypto := domain.Crypto{
		Symbol:       symbol,
		Name:         name,
		CurrentPrice: price,
		LastUpdated:  time.Now().UTC(),
	}
	if err := s.repo.Create(newCrypto); err != nil {
		return domain.Crypto{}, err
	}
	return newCrypto, nil
}

func (s *Service) GetInfoCrypto(symbol string) (domain.Crypto, error) {
	symbol = strings.ToUpper(symbol)
	getCrypto, err := s.repo.Get(symbol)
	if err != nil {
		return domain.Crypto{}, err
	}
	return getCrypto, nil
}

func (s *Service) UpdateCrypto(symbol string) (domain.Crypto, error) {
	symbol = strings.ToUpper(symbol)
	_, price, err := s.cg.GetPrice(symbol)
	if err != nil {
		return domain.Crypto{}, err
	}

	newPrice := domain.PricePoint{
		Price:     price,
		Timestamp: time.Now().UTC(),
	}
	if err = s.repo.UpdatePrice(symbol, newPrice); err != nil {
		return domain.Crypto{}, err
	}
	updatedCrypto, err := s.GetInfoCrypto(symbol)
	if err != nil {
		return domain.Crypto{}, err
	}
	return updatedCrypto, nil
}

func (s *Service) GetHistory(symbol string) ([]domain.PricePoint, error) {
	symbol = strings.ToUpper(symbol)

	historyList, err := s.repo.History(symbol)

	if err != nil {
		return []domain.PricePoint{}, err
	}

	return historyList, nil
}

func (s *Service) GetStats(symbol string) (currentPrice float64, stats domain.Stats, err error) {
	symbol = strings.ToUpper(symbol)
	crypto, err := s.GetInfoCrypto(symbol)

	if err != nil {
		return 0, domain.Stats{}, err
	}
	currentPrice = crypto.CurrentPrice
	list, err := s.GetHistory(symbol)

	if err != nil {
		return 0, domain.Stats{}, err
	}
	newStats := domain.Stats{}
	if len(list) == 0 {
		return currentPrice, domain.Stats{}, nil
	}

	prices := make([]float64, 0, len(list))
	var sum float64
	for _, point := range list {
		prices = append(prices, point.Price)
		sum += point.Price
	}
	newStats.MinPrice = slices.Min(prices)
	newStats.MaxPrice = slices.Max(prices)
	newStats.AvgPrice = sum / float64(len(list))
	newStats.RecordsCount = len(list)

	if list[0].Price == 0 {
		return currentPrice, newStats, nil
	}
	newStats.PriceChange = list[len(list)-1].Price - list[0].Price
	newStats.PriceChangePercent = newStats.PriceChange / list[0].Price * 100
	return currentPrice, newStats, nil
}

func (s *Service) DeleteCrypto(symbol string) error {
	symbol = strings.ToUpper(symbol)

	err := s.repo.Delete(symbol)

	if err != nil {
		return err
	}

	return nil
}
