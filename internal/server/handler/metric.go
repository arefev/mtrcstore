package handler

import (
	"errors"
	"fmt"
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

	value, err := h.Storage.Find(mType, mName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	strVal := strconv.FormatFloat(value, 'f', -1, 64)
	resp := fmt.Sprintf("%s\n", strVal)
	w.Write([]byte(resp))
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
