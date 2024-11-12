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
		ID: mName,
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
}

func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		log.Printf("handler Find metric fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := r.PathValue("name")

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
		m, err := h.Storage.FindCounter(mName)
		check(m.DeltaString(), err)
		return
	default:
		m, err := h.Storage.FindGauge(mName)
		check(m.ValueString(), err)
		return
	}
}

func (h *MetricHandlers) UpdateJson(w http.ResponseWriter, r *http.Request) {
	var model model.Metric
	data := json.NewDecoder(r.Body)
	
	if err := data.Decode(&model); err != nil {
		logger.Log.Error("handler UpdateJson metrics: json decode failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(model.MType); err != nil {
		logger.Log.Error("handler UpdateJson metrics fail", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Save(model); err != nil {
		logger.Log.Error("handler UpdateJson metrics: data save failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := json.NewEncoder(w)
	if err := resp.Encode(model); err != nil {
		logger.Log.Error("handler UpdateJson metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *MetricHandlers) FindJson(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	data := json.NewDecoder(r.Body)
	
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


	check := func(metric *model.Metric, err error) {
		if err != nil {
			log.Printf("handler FindJson metric fail: %s", err.Error())
			w.WriteHeader(http.StatusNotFound)
			return
		}

		resp := json.NewEncoder(w)
		if err := resp.Encode(metric); err != nil {
			logger.Log.Error("handler FindJson metrics: response writer failed", zap.Error(err))
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

func (h *MetricHandlers) getType(r *http.Request) (string, error) {
	t := r.PathValue("type")
	return t, h.checkType(t)
}

func (h *MetricHandlers) checkType(t string) error {
	if t != "counter" && t != "gauge" {
		return errors.New("metric's type is invalid")
	}

	return nil
}
