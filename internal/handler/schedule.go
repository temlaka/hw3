package handler

import (
	"encoding/json"
	"errors"
	"hw3/internal/domain"
	"hw3/internal/usecase/schedule"
	"log"
	"net/http"
	"time"
)

type scheduleDTO struct {
	Enabled         bool      `json:"enabled"`
	IntervalSeconds int       `json:"interval_seconds"`
	LastUpdate      time.Time `json:"last_update"`
	NextUpdate      time.Time `json:"next_update"`
}

type ScheduleHandler struct {
	svc *schedule.Service
}

func NewScheduleHandler(svc *schedule.Service) *ScheduleHandler {
	return &ScheduleHandler{
		svc: svc,
	}
}

func (h *ScheduleHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	x := h.svc.Get()
	respSchedule := scheduleDTO{
		Enabled:         x.Enabled,
		IntervalSeconds: x.IntervalSeconds,
		LastUpdate:      x.LastUpdate,
		NextUpdate:      x.NextUpdate,
	}
	writeJSON(w, http.StatusOK, respSchedule)
}

func (h *ScheduleHandler) SetSettings(w http.ResponseWriter, r *http.Request) {

	var bodyDecode struct {
		Enabled         bool `json:"enabled"`
		IntervalSeconds int  `json:"interval_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&bodyDecode); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.svc.EditSettings(bodyDecode.Enabled, bodyDecode.IntervalSeconds)
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, "interval_seconds must be between 10 and 3600", http.StatusBadRequest)
		return
	case err != nil:
		log.Printf("edit schedule: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, bodyDecode)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	count, timestamp, err := h.svc.Trigger()
	if err != nil {
		log.Printf("trigger update: %v", err)
		writeError(w, "internal error", http.StatusInternalServerError)
		return
	}
	var resp struct {
		UpdatedCount int       `json:"updated_count"`
		TimeStamp    time.Time `json:"timestamp"`
	}
	resp.UpdatedCount = count
	resp.TimeStamp = timestamp
	writeJSON(w, http.StatusOK, resp)
}
