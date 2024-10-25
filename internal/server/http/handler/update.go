package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arefev/mtrcstore/internal/server/http/middleware"
)

type MetricParam struct {
	Type  string
	Name  string
	Value float64
}

func UpdateHandler(mux *http.ServeMux) {
	mux.Handle("/update/", middleware.Post(http.HandlerFunc(update)))
}

func update(w http.ResponseWriter, r *http.Request) {
	parsed, err := parseUrl(r.URL.Path)

	if err != nil {
		fmt.Printf("Url parse error: %s\n", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params, err := getMetric(parsed)
	if err != nil {
		fmt.Printf("Value parse error: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := checkType(&params); err != nil {
		fmt.Printf("Type error: %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Printf("Params is %+v\n", params)
	w.Write([]byte("Metrics are updated!"))
}

func parseUrl(path string) ([]string, error) {
	const pathSize = 4
	var parsed = make([]string, pathSize)
	parsed = strings.Split(path, "/")
	if len(parsed) > 1 {
		parsed = parsed[1:]
	}

	if len(parsed) != pathSize {
		return parsed, errors.New("path url must be view as /update/{type}/{name}/{value}")
	}

	return parsed, nil
}

func getMetric(split []string) (MetricParam, error) {
	value, err := strconv.ParseFloat(split[3], 64)
	if err != nil {
		return MetricParam{}, err
	}

	param := MetricParam{
		Type:  split[1],
		Name:  split[2],
		Value: value,
	}

	return param, nil
}

func checkType(p *MetricParam) error {
	if p.Type != "counter" && p.Type != "gauge" {
		return errors.New("type is invalid")
	}

	return nil
}
