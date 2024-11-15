package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/worker"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_main(t *testing.T) {
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
			if err := logger.Init("debug"); err != nil {
				fmt.Printf("logger init failed: %v", err)
			}

			storage := repository.NewMemory()
			metricHandlers := handler.MetricHandlers{
				Storage: &storage,
			}

			const interval = 300
			const fileStoragePath = "./storage.json"
			const restore = true
			go worker.
				Init(interval, fileStoragePath, restore, &storage).
				Run()

			r := server.InitRouter(&metricHandlers)
			srv := httptest.NewServer(r)
			defer srv.Close()

			// делаем запрос с помощью библиотеки resty к адресу запущенного сервера,
			// который хранится в поле URL соответствующей структуры
			req := resty.New().R()
			req.Method = http.MethodPost
			req.URL = srv.URL + test.want.urlPath

			res, err := req.Send()
			fmt.Printf("%v", err)

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
