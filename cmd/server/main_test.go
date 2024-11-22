package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/mocks"
	"github.com/arefev/mtrcstore/internal/server/model"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_Get(t *testing.T) {
	type want struct {
		value map[string]string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "1 elem",
			want: want{
				value: map[string]string{
					"PollCounter": "1",
				},
			},
		},
		{
			name: "2 elems",
			want: want{
				value: map[string]string{
					"PollCounter": "1",
					"Alloc":       "2000",
				},
			},
		},
		{
			name: "empty",
			want: want{
				value: map[string]string{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)
			storage.EXPECT().Get().MaxTimes(1).Return(test.want.value)

			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = srv.URL

			res, err := req.Send()
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode())

			for k, v := range test.want.value {
				require.Contains(
					t,
					string(res.Body()),
					fmt.Sprintf("<li><strong>%s</strong>: %s</li>", k, v),
				)
			}
		})
	}
}

func Test_UpdateShortUrl(t *testing.T) {
	const urlPath string = "/update"
	var delta int64 = 1
	var value float64 = 1.55

	type want struct {
		metric     model.Metric
		err        error
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "save counter",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				err:        nil,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "save gauge",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "bad request, when returned error",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
				},
				err:        fmt.Errorf("saved failed"),
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)
			storage.EXPECT().Save(test.want.metric).MaxTimes(1).Return(test.want.err)

			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.SetHeader("Content-type", "application/json")
			req.URL = srv.URL + urlPath

			jsonValue, err := json.Marshal(test.want.metric)
			require.NoError(t, err)

			req.Body = jsonValue

			res, err := req.Send()
			require.NoError(t, err)
			require.Equal(t, test.want.statusCode, res.StatusCode())

			if test.want.err == nil {
				require.Contains(t, string(res.Body()), string(jsonValue))
			}
		})
	}
}

func Test_UpdateFullUrl(t *testing.T) {
	var delta int64 = 1
	var value float64 = 1

	type want struct {
		metric     model.Metric
		err        error
		urlPath    string
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "save counter",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				err:        nil,
				urlPath:    "/update/counter/PollCounter/1",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "save gauge",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/update/gauge/Alloc/1",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "bad request with invalid value",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/update/gauge/Alloc/test",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "bad request with invalid type",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/update/test/Alloc/1",
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "404 with invalid template path №1",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/update/counter",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "404 with invalid template path №2",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/update/gauge/Alloc/1/test",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "bad request when saving failed",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        fmt.Errorf("saved failed"),
				urlPath:    "/update/gauge/Alloc/1",
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)
			storage.EXPECT().Save(test.want.metric).MaxTimes(1).Return(test.want.err)

			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL + test.want.urlPath

			require.NoError(t, err)

			res, err := req.Send()
			require.NoError(t, err)
			require.Equal(t, test.want.statusCode, res.StatusCode())
		})
	}
}

func Test_FindShortUrl(t *testing.T) {
	const urlPath string = "/value"
	var delta int64 = 1
	var value float64 = 1.55

	type want struct {
		metric     model.Metric
		data       model.Metric
		err        error
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "counter found",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				data: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
				},
				err:        nil,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "gauge found",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				data: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
				},
				err:        nil,
				statusCode: http.StatusOK,
			},
		},
		{
			name: "404 when counter id not found",
			want: want{
				metric: model.Metric{},
				data: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
				},
				err:        fmt.Errorf("not found"),
				statusCode: http.StatusNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)
			storage.
				EXPECT().
				Find(test.want.data.ID, test.want.data.MType).
				MaxTimes(1).
				Return(test.want.metric, test.want.err)

			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodPost
			req.SetHeader("Content-type", "application/json")
			req.URL = srv.URL + urlPath

			jsonData, err := json.Marshal(test.want.data)
			require.NoError(t, err)

			req.Body = jsonData
			res, err := req.Send()
			require.NoError(t, err)

			jsonValue, err := json.Marshal(test.want.metric)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, test.want.statusCode, res.StatusCode())

			if test.want.err == nil {
				require.Contains(t, string(res.Body()), string(jsonValue))
			}
		})
	}
}
func Test_FindFullUrl(t *testing.T) {
	var delta int64 = 1
	var value float64 = 1.55

	type want struct {
		metric     model.Metric
		err        error
		urlPath    string
		statusCode int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "counter found",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				err:        nil,
				urlPath:    "/value/counter/PollCounter",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "gauge found",
			want: want{
				metric: model.Metric{
					ID:    "Alloc",
					MType: "gauge",
					Value: &value,
				},
				err:        nil,
				urlPath:    "/value/gauge/Alloc",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "404 when counter id not found",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				err:        fmt.Errorf("not found"),
				urlPath:    "/value/counter/PollCounter",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "bad request when invalid type",
			want: want{
				metric: model.Metric{
					ID:    "PollCounter",
					MType: "counter",
					Delta: &delta,
				},
				err:        fmt.Errorf("bad request"),
				urlPath:    "/value/test/PollCounter",
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storage := mocks.NewMockStorage(ctrl)
			storage.
				EXPECT().
				Find(test.want.metric.ID, test.want.metric.MType).
				MaxTimes(1).
				Return(test.want.metric, test.want.err)

			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = srv.URL + test.want.urlPath

			require.NoError(t, err)

			res, err := req.Send()
			require.NoError(t, err)
			require.Equal(t, test.want.statusCode, res.StatusCode())
			if test.want.err == nil {
				var value string
				if test.want.metric.MType == "counter" {
					value = test.want.metric.DeltaString()
				} else {
					value = test.want.metric.ValueString()
				}
				require.Contains(t, string(res.Body()), value)
			}
		})
	}
}
