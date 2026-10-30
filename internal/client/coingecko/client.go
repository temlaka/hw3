package coingecko

import (
	"encoding/json"
	"fmt"
	"hw3/internal/domain"
	"net/http"
	"strings"
	"sync"
	"time"
)

type coin struct {
	id   string
	name string
}

type Client struct {
	mu      sync.RWMutex
	client  *http.Client
	cache   map[string]coin
	baseURL string
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		client:  &http.Client{Timeout: 10 * time.Second},
		cache:   make(map[string]coin),
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

//https://api.coingecko.com/api/v3 base url

func (c *Client) search(symbol string) (coin, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search?query=%s", c.baseURL, symbol), nil)
	if err != nil {
		return coin{}, fmt.Errorf("search %s: build request: %w", symbol, err)
	}

	req.Header.Set("x-cg-demo-api-key", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return coin{}, fmt.Errorf("search %s: %w", symbol, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return coin{}, fmt.Errorf("search %s: unexpected status %d", symbol, resp.StatusCode)

	}
	var body struct {
		Coins []struct {
			Id     string `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
		} `json:"coins"`
	}
	// example{"id":"solana","name":"Solana","api_symbol":"solana","symbol":"SOL","market_cap_rank":7,
	// "thumb":"https://coin-images.coingecko.com/coins/images/4128/thumb/solana.png",
	// "large":"https://coin-images.coingecko.com/coins/images/4128/large/solana.png"}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {

		return coin{}, fmt.Errorf("search %s: decode response: %w", symbol, err)
	}

	for _, item := range body.Coins {
		if strings.EqualFold(item.Symbol, symbol) {
			return coin{
				id:   item.Id,
				name: item.Name,
			}, nil
		}
	}

	return coin{}, fmt.Errorf("search %s: %w", symbol, domain.ErrUnknownSymbol)

}

func (c *Client) price(id string) (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/simple/price?vs_currencies=usd&ids=%s", c.baseURL, id), nil)
	if err != nil {
		return 0, fmt.Errorf("price %s: build request: %w", id, err)
	}

	req.Header.Set("x-cg-demo-api-key", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("price %s: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("price %s: unexpected status %d", id, resp.StatusCode)

	}
	body := make(map[string]map[string]float64, 0)

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {

		return 0, fmt.Errorf("price %s: decode response: %w", id, err)
	}
	prices, ok := body[id]["usd"]
	if !ok {
		return 0, fmt.Errorf("price %s: no usd price in response", id)
	}
	return prices, nil

}

func (c *Client) GetPrice(symbol string) (name string, price float64, err error) {
	c.mu.RLock()
	symbol = strings.ToUpper(symbol)
	getCoin, ok := c.cache[symbol]
	c.mu.RUnlock()
	if !ok {
		getCoin, err = c.search(symbol)
		if err != nil {
			return "", 0, err
		}
		c.mu.Lock()
		c.cache[symbol] = getCoin
		c.mu.Unlock()
	}
	price, err = c.price(getCoin.id)
	if err != nil {
		return "", 0, err
	}

	return getCoin.name, price, nil

}
