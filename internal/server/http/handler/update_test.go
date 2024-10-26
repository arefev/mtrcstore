package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateHandler_update(t *testing.T) {
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
				urlPath:      "/update/counter/test/123",
				code:         http.StatusOK,
				response:     `Metrics are updated!`,
				contentType:  "text/plain; charset=utf-8",
				testStorage:  false,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 123,
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
				storageValue: 123,
			},
		},
		{
			name: "not found test",
			want: want{
				urlPath:      "/update/counter/test",
				code:         http.StatusNotFound,
				response:     ``,
				contentType:  "",
				testStorage:  false,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 123,
			},
		},
		{
			name: "positive counter storage test #1",
			want: want{
				urlPath:      "/update/counter/test/123",
				code:         http.StatusOK,
				response:     `Metrics are updated!`,
				contentType:  "text/plain; charset=utf-8",
				testStorage:  true,
				storageType:  "counter",
				storageName:  "test",
				storageValue: 123,
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
			request := httptest.NewRequest(http.MethodGet, test.want.urlPath, nil)

			storage := repository.NewMemory()
			handler := UpdateHandler{
				Storage: &storage,
			}

			// создаём новый Recorder
			w := httptest.NewRecorder()
			handler.update(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

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
