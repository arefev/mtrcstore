package handler

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func ExampleMetricHandlers_Update() {
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = "http://localhost:8080/update/gauge/Alloc/1.55"

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response status code %d", res.StatusCode())
	fmt.Println(response)
}

func ExampleMetricHandlers_Find() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080/value/counter/PollCount"

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response status code %d", res.StatusCode())
	fmt.Println(response)
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

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response status code %d", res.StatusCode())
	fmt.Println(response)
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

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response body %s", res.Body())
	fmt.Println(response)
}

func ExampleMetricHandlers_Get() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080"

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response body %s", res.Body())
	fmt.Println(response)
}

func ExampleMetricHandlers_Ping() {
	req := resty.New().R()
	req.Method = http.MethodGet
	req.URL = "http://localhost:8080/ping"

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response status code %d", res.StatusCode())
	fmt.Println(response)
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

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	response := fmt.Sprintf("Response status code %d", res.StatusCode())
	fmt.Println(response)
}
