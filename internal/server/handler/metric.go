package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"github.com/arefev/mtrcstore/internal/server/worker"
	"go.uber.org/zap"
)

type MetricHandlers struct {
	Storage repository.Storage
}

func (h *MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		log.Printf("handler Update metrics fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := r.PathValue("name")
	mValue, err := strconv.ParseFloat(r.PathValue("value"), 64)

	if err != nil {
		log.Printf("handler Update metrics fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ival := int64(mValue)
	metric := model.Metric{
		ID:    mName,
		MType: mType,
		Value: &mValue,
		Delta: &ival,
	}

	if err := h.Storage.Save(metric); err != nil {
		log.Printf("handler Update metrics: data save failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte("Metrics are updated!")); err != nil {
		log.Printf("handler Update metrics: response writer failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	worker.Worker.SaveEvent()
}

func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		log.Printf("handler Find metric fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.Storage.Find(r.PathValue("name"), mType)
	if err != nil {
		log.Printf("handler Find metric failed: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var value string
	switch mType {
	case repository.CounterName:
		value = metric.DeltaString()
	default:
		value = metric.ValueString()
	}

	if _, err := w.Write([]byte(value)); err != nil {
		log.Printf("handler Find metrics: response writer failed: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	data := json.NewDecoder(r.Body)

	w.Header().Add("Content-type", "application/json")

	if err := data.Decode(&metric); err != nil {
		logger.Log.Error("handler UpdateJson metrics: json decode failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(metric.MType); err != nil {
		logger.Log.Error("handler UpdateJson metrics fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Save(metric); err != nil {
		logger.Log.Error("handler UpdateJson metrics: data save failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := json.NewEncoder(w)
	if err := resp.Encode(metric); err != nil {
		logger.Log.Error("handler UpdateJson metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	worker.Worker.SaveEvent()
}

func (h *MetricHandlers) FindJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	data := json.NewDecoder(r.Body)

	w.Header().Add("Content-type", "application/json")

	if err := data.Decode(&metric); err != nil {
		logger.Log.Error("handler FindJson metrics: json decode failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(metric.MType); err != nil {
		logger.Log.Error("handler FindJson metrics fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.Storage.Find(metric.ID, metric.MType)
	if err != nil {
		log.Printf("handler FindJson metric failed: %s", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := json.NewEncoder(w)
	if err := resp.Encode(value); err != nil {
		logger.Log.Error("handler FindJson metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := service.ListHTML(w, h.Storage.Get()); err != nil {
		log.Printf("handler Get fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) getType(r *http.Request) (string, error) {
	t := r.PathValue("type")
	return t, h.checkType(t)
}

func (h *MetricHandlers) checkType(t string) error {
	if t != repository.CounterName && t != repository.GaugeName {
		return errors.New("metric's type is invalid")
	}

	return nil
}
