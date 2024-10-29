package handler

import (
	"errors"
	// "fmt"
	"net/http"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
	"github.com/go-chi/chi/v5"
)

type MetricHandlers struct {
	Storage repository.Storage
}

func (h *MetricHandlers) Update(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := chi.URLParam(r, "name")
	mValue, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.Storage.Save(mType, mName, mValue)
	w.Write([]byte("Metrics are updated!"))
}

func (h *MetricHandlers) Find(w http.ResponseWriter, r *http.Request) {
	mType, err := h.getType(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := chi.URLParam(r, "name")

	check := func(str string, err error) {
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(str))
	}

	switch mType {
	case "counter":
		value, err := h.Storage.FindCounter(mName)
		check(value.String(), err)
	default:
		value, err := h.Storage.FindGauge(mName)
		check(value.String(), err)
	}
}

func (h *MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	html := service.MetricsHTML(h.Storage.Get())
	w.Write([]byte(html))
}

func (h *MetricHandlers) getType(r *http.Request) (string, error) {
	t := chi.URLParam(r, "type")
	if t != "counter" && t != "gauge" {
		return "", errors.New("metric's type is invalid")
	}

	return t, nil
}
