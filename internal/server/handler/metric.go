package handler

import (
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
	mType := chi.URLParam(r, "type")
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
	mType := chi.URLParam(r, "type")
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

