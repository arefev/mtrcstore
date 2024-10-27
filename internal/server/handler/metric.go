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

	fmt.Printf("Type %s, Name %s, Value %f\n", mType, mName, mValue)
	fmt.Printf("Storage has %+v\n", h.Storage)
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

	resp := fmt.Sprintf("Type %s, Name %s, Value %f\n", mType, mName, value)
	w.Write([]byte(resp))
}

func (h *MetricHandlers) Get(w http.ResponseWriter, r *http.Request) {
	html := service.MetricsHtml(h.Storage.Get())
	fmt.Printf("Metrics %+v\n", h.Storage.Get())
	w.Write([]byte(html))
}

