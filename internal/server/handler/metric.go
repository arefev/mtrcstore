package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type MetricHandlers struct {
	Storage repository.Storage
	db      *sql.DB
	log     *zap.Logger
}

func NewMetricHandlers(s repository.Storage, db *sql.DB, log *zap.Logger) *MetricHandlers {
	m := MetricHandlers{
		Storage: s,
		db:      db,
		log:     log,
	}
	return &m
}

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

	delta := int64(mValue)
	metric := model.Metric{
		ID:    mName,
		MType: mType,
		Value: &mValue,
		Delta: &delta,
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

func (h *MetricHandlers) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	if err := h.db.PingContext(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte("DB connected!")); err != nil {
		h.log.Error("handler Ping metrics: response writer failed", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
