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
			name: "positive test 1 elem",
			want: want{
				value: map[string]string{
					"PollCounter": "1",
				},
			},
		},
		{
			name: "positive test 2 elems",
			want: want{
				value: map[string]string{
					"PollCounter": "1",
					"Alloc":       "2000",
				},
			},
		},
		{
			name: "positive test empty",
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
			name: "positive test counter saved",
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
			name: "positive test gauge saved",
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
