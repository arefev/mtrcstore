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
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/go-resty/resty/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UpdateFullUrl(t *testing.T) {
	const interval = 300
	const fileStoragePath = "./storage.json"
	const restore = true
			
	type want struct {
		urlPath      string
		code         int
		response     string
		contentType  string
		testStorage  bool
		storageType  string
		storageName  string
		storageValue float64
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				urlPath:      "/update/counter/test/1",
				code:         http.StatusOK,
				response:     `Metrics are updated!`,
				contentType:  "text/plain; charset=utf-8",
				testStorage:  false,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 1,
			},
		},
		{
			name: "bad request test",
			want: want{
				urlPath:      "/update/counter/test/asd",
				code:         http.StatusBadRequest,
				response:     ``,
				contentType:  "",
				testStorage:  false,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 1,
			},
		},
		{
			name: "not found test",
			want: want{
				urlPath:      "/update/counter/test",
				code:         http.StatusNotFound,
				response:     "404 page not found\n",
				contentType:  "text/plain; charset=utf-8",
				testStorage:  false,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 1,
			},
		},
		{
			name: "positive counter storage test #1",
			want: want{
				urlPath:      "/update/counter/test/1",
				code:         http.StatusOK,
				response:     `Metrics are updated!`,
				contentType:  "text/plain; charset=utf-8",
				testStorage:  true,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 1,
			},
		},
		{
			name: "positive gauge storage test #1",
			want: want{
				urlPath:      "/update/gauge/test/45.56",
				code:         http.StatusOK,
				response:     `Metrics are updated!`,
				contentType:  "text/plain; charset=utf-8",
				testStorage:  true,
				storageType:  "gauge",
				storageName:  "test",
				storageValue: 45.56,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cLog, err := logger.Build("debug")
			require.NoError(t, err)

			storage := repository.NewFile(interval, fileStoragePath, restore, cLog)
			metricHandlers := handler.NewMetricHandlers(storage, cLog)

			r := server.InitRouter(metricHandlers, cLog)
			srv := httptest.NewServer(r)
			defer srv.Close()

			// делаем запрос с помощью библиотеки resty к адресу запущенного сервера,
			// который хранится в поле URL соответствующей структуры
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL + test.want.urlPath

			res, err := req.Send()

			require.NoError(t, err)
			assert.Equal(t, test.want.code, res.StatusCode())
			assert.Equal(t, test.want.response, string(res.Body()))
			assert.Equal(t, test.want.contentType, res.Header().Get("Content-type"))

			// Проверка сохранения данных
			if test.want.testStorage {
				var sValue float64
				tValue := test.want.storageValue

				switch test.want.storageType {
				case "counter":
					sValue = float64(storage.Counter[test.want.storageName])
				default:
					sValue = float64(storage.Gauge[test.want.storageName])
				}

				assert.Equal(t, tValue, sValue)
			}
		})
	}
}

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
