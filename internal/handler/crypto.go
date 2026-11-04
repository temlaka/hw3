package handler

import (
	"encoding/json"
	"errors"
	"hw3/internal/domain"
	"hw3/internal/usecase/crypto"
	"log"
	"net/http"
	"strings"
	"time"
)

type CryptoHandler struct {
	svc *crypto.Service
}

type cryptoResponse struct {
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

type pricePointResponse struct {
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}

type history struct {
	Symbol  string               `json:"symbol"`
	History []pricePointResponse `json:"history"`
}

type stats struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
	Stats        struct {
		MinPrice           float64 `json:"min_price"`
		MaxPrice           float64 `json:"max_price"`
		AvgPrice           float64 `json:"avg_price"`
		PriceChange        float64 `json:"price_change"`
		PriceChangePercent float64 `json:"price_change_percent"`
		RecordsCount       int     `json:"records_count"`
	} `json:"stats"`
}

func NewCryptoHandler(svc *crypto.Service) *CryptoHandler {
	return &CryptoHandler{
		svc: svc,
	}
}

func toCryptoResponse(from domain.Crypto) cryptoResponse {
	to := cryptoResponse{
		Symbol:       from.Symbol,
		Name:         from.Name,
		CurrentPrice: from.CurrentPrice,
		LastUpdated:  from.LastUpdated,
	}
	return to
}

func (h *CryptoHandler) GetCryptos(w http.ResponseWriter, r *http.Request) {

	list, err := h.svc.GetListCrypto()
	if err != nil {
		log.Printf("get crypto: %v", err)
		writeError(w, "internal", http.StatusInternalServerError)
		return
	}

	var cryptos struct {
		Cryptos []cryptoResponse `json:"cryptos"`
	}

	cryptos.Cryptos = make([]cryptoResponse, 0, len(list))

	for _, item := range list {
		cryptos.Cryptos = append(cryptos.Cryptos, toCryptoResponse(item))
	}
	writeJSON(w, http.StatusOK, cryptos)
}

func (h *CryptoHandler) AddCrypto(w http.ResponseWriter, r *http.Request) {
	bodyDecode := make(map[string]string, 0)

	if err := json.NewDecoder(r.Body).Decode(&bodyDecode); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	newCrypto, err := h.svc.AddCrypto(bodyDecode["symbol"])
	switch {
	case errors.Is(err, domain.ErrUnknownSymbol):
		writeError(w, "unknown symbol", http.StatusBadRequest)
		return
	case errors.Is(err, domain.ErrCryptoExists):
		writeError(w, "crypto already tracked", http.StatusConflict)
		return
	case err != nil:
		log.Printf("add crypto: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	var jsonResp struct {
		Crypto cryptoResponse `json:"crypto"`
	}
	jsonResp.Crypto = toCryptoResponse(newCrypto)
	writeJSON(w, http.StatusCreated, jsonResp)
}

func (h *CryptoHandler) GetCryptoInfo(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	getCrypto, err := h.svc.GetInfoCrypto(symbol)
	switch {
	case errors.Is(err, domain.ErrCryptoNotFound):
		writeError(w, "crypto not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get crypto: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, toCryptoResponse(getCrypto))

}

func (h *CryptoHandler) UpdateCryptoInfo(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	updatedCrypto, err := h.svc.UpdateCrypto(symbol)
	switch {
	case errors.Is(err, domain.ErrCryptoNotFound), errors.Is(err, domain.ErrUnknownSymbol):
		writeError(w, "crypto not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("refresh crypto: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	var jsonResp struct {
		Crypto cryptoResponse `json:"crypto"`
	}
	jsonResp.Crypto = toCryptoResponse(updatedCrypto)
	writeJSON(w, http.StatusOK, jsonResp)

}

func (h *CryptoHandler) GetCryptoHistory(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	historyList, err := h.svc.GetHistory(symbol)
	switch {
	case errors.Is(err, domain.ErrCryptoNotFound):
		writeError(w, "crypto not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get history: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}

	var respHistory history
	respHistory.Symbol = strings.ToUpper(symbol)

	respHistory.History = make([]pricePointResponse, 0)
	for _, item := range historyList {

		x := pricePointResponse{}

		x.Price = item.Price
		x.Timestamp = item.Timestamp

		respHistory.History = append(respHistory.History, x)
	}

	writeJSON(w, http.StatusOK, respHistory)

}

func (h *CryptoHandler) GetStatsCrypto(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	currentPrice, getStats, err := h.svc.GetStats(symbol)
	switch {
	case errors.Is(err, domain.ErrCryptoNotFound):
		writeError(w, "crypto not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("get stats: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	var respStats stats
	respStats.CurrentPrice = currentPrice
	respStats.Symbol = strings.ToUpper(symbol)
	respStats.Stats.MinPrice = getStats.MinPrice
	respStats.Stats.MaxPrice = getStats.MaxPrice
	respStats.Stats.AvgPrice = getStats.AvgPrice
	respStats.Stats.PriceChange = getStats.PriceChange
	respStats.Stats.PriceChangePercent = getStats.PriceChangePercent
	respStats.Stats.RecordsCount = getStats.RecordsCount
	writeJSON(w, http.StatusOK, respStats)

}

func (h *CryptoHandler) DeleteCrypto(w http.ResponseWriter, r *http.Request) {
	symbol := r.PathValue("symbol")
	err := h.svc.DeleteCrypto(symbol)
	switch {
	case errors.Is(err, domain.ErrCryptoNotFound):
		writeError(w, "crypto not found", http.StatusNotFound)
		return
	case err != nil:
		log.Printf("delete crypto: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	var x struct{}
	writeJSON(w, http.StatusOK, x)

}
