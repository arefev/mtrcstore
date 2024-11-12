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
	"go.uber.org/zap"
)

type MetricHandlers struct {
	Storage repository.Storage
}

func (h *MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	var model model.Metric
	data := json.NewDecoder(r.Body)
	
	if err := data.Decode(&model); err != nil {
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

	resp := json.NewEncoder(w)
	if err := resp.Encode(model); err != nil {
		logger.Log.Error("handler Update metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	data := json.NewDecoder(r.Body)
	
	if err := data.Decode(&metric); err != nil {
		logger.Log.Error("handler Find metrics: json decode failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(metric.MType); err != nil {
		logger.Log.Error("handler Find metrics fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	check := func(metric *model.Metric, err error) {
		if err != nil {
			log.Printf("handler Find metric fail: %s", err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp := json.NewEncoder(w)
		if err := resp.Encode(metric); err != nil {
			logger.Log.Error("handler Find metrics: response writer failed", zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	switch metric.MType {
	case "counter":
		value, err := h.Storage.FindCounter(metric.ID)
		check(&value, err)
		return
	default:
		value, err := h.Storage.FindGauge(metric.ID)
		check(&value, err)
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

func (h *MetricHandlers) checkType(t string) error {
	if t != "counter" && t != "gauge" {
		return errors.New("metric's type is invalid")
	}

	return nil
}
