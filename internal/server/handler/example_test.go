package handler

import (
	"net/http"

	"github.com/go-resty/resty/v2"
)

func ExampleMetricHandlers_Update() {
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = "http://localhost:8080/update/gauge/Alloc/1.55"

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_Find() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080/value/counter/PollCount"

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_UpdateJSON() {
	data := `
	{
		"id": "PollCount",
		"type": "counter",
		"delta": 4
	}
	`
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = "http://localhost:8080/update/"
	req.Body = data

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_FindJSON() {
	data := `
	{
		"id": "PollCount",
		"type": "counter"
	}
	`
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = "http://localhost:8080/value/"
	req.Body = data

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_Get() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080"

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_Ping() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080/ping"

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}

func ExampleMetricHandlers_Updates() {
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
	req.URL = "http://localhost:8080/updates/"
	req.Body = data

	_, err := req.Send()
	if err != nil {
		panic(err)
	}

	// Output:
}
