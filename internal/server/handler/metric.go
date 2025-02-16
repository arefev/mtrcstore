package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"go.uber.org/zap"
)

// @Title MetricsStore API
// @Description Metrics storage service
// @Version 1.0

// @BasePath /
// @Host localhost:8080

// @Tag.name Info
// @Tag.description "Group of requests to get metrics"

// @Tag.name Update
// @Tag.description "Group of requests to update metrics"

type MetricHandlers struct {
	Storage repository.Storage
	log     *zap.Logger
}

func NewMetricHandlers(s repository.Storage, log *zap.Logger) *MetricHandlers {
	m := MetricHandlers{
		Storage: s,
		log:     log,
	}
	return &m
}

// Update godoc
// @Tags Update
// @Summary Update metric by type and name
// @ID updateMetric
// @Accept  text/html
// @Produce text/html
// @Param type path string true "metric type [counter, gauge]"
// @Param name path string true "metric name"
// @Param value path number true "metric value"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /update/{type}/{name}/{value} [post]
func (h *MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := r.PathValue("name")
	mValue, err := strconv.ParseFloat(r.PathValue("value"), 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric := model.Metric{
		ID:    mName,
		MType: mType,
	}

	switch mType {
	case repository.CounterName:
		delta := int64(mValue)
		metric.Delta = &delta
	default:
		metric.Value = &mValue
	}

	if err := h.Storage.Save(metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte("Metrics are updated!")); err != nil {
		h.log.Error("handler Update metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Find godoc
// @Tags Info
// @Summary Find metric by type and name
// @ID findMetric
// @Accept  text/html
// @Produce text/html
// @Param type path string true "metric type [counter, gauge]"
// @Param name path string true "metric name"
// @Success 200 {string} number "metric's value, for example 200.4"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /value/{type}/{name} [get]
func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric, err := h.Storage.Find(r.PathValue("name"), mType)
	if err != nil {
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
		h.log.Error("handler Find metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// UpdateJson godoc
// @Tags Update
// @Summary Update metric with json format
// @ID updateJSONMetric
// @Accept  application/json
// @Produce application/json
// @Param metric body model.Metric true "Metric's data"
// @Success 200 {object} model.Metric "Metric's data"
// @Failure 400
// @Failure 500
// @Router /update [post]
func (h *MetricHandlers) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	d := json.NewDecoder(r.Body)

	w.Header().Add("Content-type", "application/json")

	if err := d.Decode(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(metric.MType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Save(metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := json.NewEncoder(w)
	if err := resp.Encode(metric); err != nil {
		h.log.Error("handler UpdateJson metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// FindJson godoc
// @Tags Info
// @Summary Get metric info with json format
// @ID findJSONMetric
// @Accept  application/json
// @Produce application/json
// @Param metric body model.Metric true "Metric's data"
// @Success 200 {object} model.Metric "Metric's data"
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /value/ [post]
func (h *MetricHandlers) FindJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	data := json.NewDecoder(r.Body)

	w.Header().Add("Content-type", "application/json")

	if err := data.Decode(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.checkType(metric.MType); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value, err := h.Storage.Find(metric.ID, metric.MType)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp := json.NewEncoder(w)
	if err := resp.Encode(value); err != nil {
		h.log.Error("handler FindJson metric: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Get godoc
// @Tags Info
// @Summary Get metrics list
// @ID getMetric
// @Accept  text/html
// @Produce text/html
// @Success 200
// @Failure 500
// @Router / [get]
func (h *MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := service.ListHTML(w, h.Storage.Get()); err != nil {
		h.log.Error("handler Get failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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

// Ping godoc
// @Tags Info
// @Summary Check storage status
// @ID pingMetric
// @Accept  text/html
// @Produce text/html
// @Success 200
// @Failure 500
// @Router /ping [get]
func (h *MetricHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.Storage.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte("DB connected!")); err != nil {
		h.log.Error("handler Ping metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Updates godoc
// @Tags Update
// @Summary Mass update metrics with json format
// @ID updatesMetric
// @Accept  application/json
// @Produce application/json
// @Param metric body []model.Metric true "Metric's data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /updates/ [post]
func (h *MetricHandlers) Updates(w http.ResponseWriter, r *http.Request) {
	var metrics []model.Metric
	d := json.NewDecoder(r.Body)

	if err := d.Decode(&metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.MassSave(metrics); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte("Mass save successful!")); err != nil {
		h.log.Error("handler Updates metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
