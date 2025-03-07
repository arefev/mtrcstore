package handler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/arefev/mtrcstore/internal/server"
	"github.com/arefev/mtrcstore/internal/server/handler"
	"github.com/arefev/mtrcstore/internal/server/logger"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/go-resty/resty/v2"
)

func ExampleMetricHandlers_Update() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/update/gauge/Alloc/1.55"

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))

	// Output:
	// Response code: 200
}

func ExampleMetricHandlers_Find() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/update/gauge/Alloc/1.55"

	_, err = req.Send()
	if err != nil {
		panic(err)
	}

	req = resty.New().R()
	req.Method = http.MethodGet
	req.URL = srv.URL + "/value/gauge/Alloc"

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))
	fmt.Println("Response body: " + string(response.Body()))

	// Output:
	// Response code: 200
	// Response body: 1.55
}

func ExampleMetricHandlers_UpdateJSON() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	data := `
	{
		"id": "PollCount",
		"type": "counter",
		"delta": 4
	}
	`
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/update/"
	req.Body = data

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))

	// Output:
	// Response code: 200
}

func ExampleMetricHandlers_FindJSON() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/update/gauge/Alloc/1.55"

	_, err = req.Send()
	if err != nil {
		panic(err)
	}

	data := `
	{
		"id": "Alloc",
		"type": "gauge"
	}
	`
	req = resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/value/"
	req.Body = data

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))
	fmt.Println("Response body: " + string(response.Body()))

	// Output:
	// Response code: 200
	// Response body: {"value":1.55,"id":"Alloc","type":"gauge"}
}

func ExampleMetricHandlers_Get() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = srv.URL

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))

	// Output:
	// Response code: 200
}

func ExampleMetricHandlers_Ping() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = srv.URL + "/ping"

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))

	// Output:
	// Response code: 200
}

func ExampleMetricHandlers_Updates() {
	cLog, err := logger.Build("debug")
	if err != nil {
		panic(err)
	}
	metricHandlers := handler.NewMetricHandlers(repository.NewMemory(), cLog)

	r := server.InitRouter(metricHandlers, cLog, "", "")
	srv := httptest.NewServer(r)
	defer srv.Close()

	data := `
	[
		{
			"id": "PollCount",
			"type": "counter",
			"delta": 1
		},
		{
			"id": "Alloc",
			"type": "gauge",
			"value": 5.18
		}
	]
	`
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = srv.URL + "/updates/"
	req.Body = data

	response, err := req.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("Response code: " + strconv.Itoa(response.StatusCode()))

	// Output:
	// Response code: 200
}
