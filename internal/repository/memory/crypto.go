package memory

import (
	"fmt"
	"hw3/internal/domain"
	"sync"
)

type CryptoRepo struct {
	history      map[string][]domain.PricePoint
	cryptos      map[string]domain.Crypto
	mu           sync.RWMutex
	historyLimit int
}

func NewCryptoRepo(limit int) *CryptoRepo {
	return &CryptoRepo{
		history:      make(map[string][]domain.PricePoint),
		cryptos:      make(map[string]domain.Crypto),
		historyLimit: limit,
	}
}

func (cr *CryptoRepo) Create(c domain.Crypto) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, ok := cr.cryptos[c.Symbol]
	if ok {
		return fmt.Errorf("create %s: %w", c.Symbol, domain.ErrCryptoExists)
	}
	cr.cryptos[c.Symbol] = c
	cr.history[c.Symbol] = append(cr.history[c.Symbol], domain.PricePoint{
		Price:     c.CurrentPrice,
		Timestamp: c.LastUpdated,
	})
	return nil
}

func (cr *CryptoRepo) List() ([]domain.Crypto, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	list := make([]domain.Crypto, 0, len(cr.cryptos))
	for _, value := range cr.cryptos {
		list = append(list, value)
	}
	return list, nil
}

func (cr *CryptoRepo) Get(symbol string) (domain.Crypto, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	getCrypto, ok := cr.cryptos[symbol]
	if !ok {
		return domain.Crypto{}, fmt.Errorf("get %s: %w", symbol, domain.ErrCryptoNotFound)
	}
	return getCrypto, nil
}

func (cr *CryptoRepo) UpdatePrice(symbol string, p domain.PricePoint) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	getCrypto, exist := cr.cryptos[symbol]
	if !exist {
		return fmt.Errorf("update price %s: %w", symbol, domain.ErrCryptoNotFound)
	}
	getCrypto.CurrentPrice = p.Price
	getCrypto.LastUpdated = p.Timestamp
	cr.cryptos[symbol] = getCrypto
	if cr.historyLimit <= len(cr.history[symbol]) {
		cr.history[symbol] = cr.history[symbol][1:]
	}
	cr.history[symbol] = append(cr.history[symbol], p)
	return nil
}
func (cr *CryptoRepo) History(symbol string) ([]domain.PricePoint, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()
	historySlice, exist := cr.history[symbol]
	if !exist {
		return nil, fmt.Errorf("history %s: %w", symbol, domain.ErrCryptoNotFound)
	}
	res := make([]domain.PricePoint, len(historySlice))
	copy(res, historySlice)
	return res, nil
}

func (cr *CryptoRepo) Delete(symbol string) error {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	_, ok := cr.cryptos[symbol]
	if !ok {
		return fmt.Errorf("delete %s: %w", symbol, domain.ErrCryptoNotFound)
	}
	delete(cr.cryptos, symbol)
	delete(cr.history, symbol)

	return nil

}
