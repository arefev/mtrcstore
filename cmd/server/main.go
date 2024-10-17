package main

import (
	"net/http"

	"github.com/arefev/mtrcstore/internal/http/middleware"
)


func updateMetric(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("update Metric"))
}


func main() {

	mux := http.NewServeMux()
	mux.Handle("/update/counter/someMetric/123", middleware.PostOnly(http.HandlerFunc(updateMetric)))

	err := http.ListenAndServe(":8080", mux)
	if (err != nil) {
		panic(err)
	}

}
