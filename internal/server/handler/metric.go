package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type MetricHandlers struct {
	Storage repository.Storage
}

func (h *MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var model model.Metric
	json := json.NewDecoder(r.Body)
	if err := json.Decode(&model); err != nil {
		logger.Log.Error("handler Update metrics: json decode failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(model.MType); err != nil {
		logger.Log.Error("handler Update metrics fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Save(model); err != nil {
		logger.Log.Error("handler Update metrics: data save failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		log.Printf("handler Find metric fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := chi.URLParam(r, "name")

	check := func(str string, err error) {
		if err != nil {
			log.Printf("handler Find metric fail: %s", err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := w.Write([]byte(str)); err != nil {
			log.Printf("handler Find metrics: response writer failed: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	switch mType {
	case "counter":
		value, err := h.Storage.FindCounter(mName)
		check(value.String(), err)
		return
	default:
		value, err := h.Storage.FindGauge(mName)
		check(value.String(), err)
		return
	}
}

func (h *MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	if err := service.ListHTML(w, h.Storage.Get()); err != nil {
		log.Printf("handler Get fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) getType(r *http.Request) (string, error) {
	t := chi.URLParam(r, "type")
	return t, h.checkType(t)
}

func (h *MetricHandlers) checkType(t string) error {
	if t != "counter" && t != "gauge" {
		return errors.New("metric's type is invalid")
	}

	return nil
}
