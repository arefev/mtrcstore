package handler

import (
	"errors"
	"log"
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
		log.Printf("handler Update metrics fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mName := chi.URLParam(r, "name")
	mValue, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)

	if err != nil {
		log.Printf("handler Update metrics fail: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.Storage.Save(mType, mName, mValue); err != nil {
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
	if t != "counter" && t != "gauge" {
		return "", errors.New("metric's type is invalid")
	}

	return t, nil
}
