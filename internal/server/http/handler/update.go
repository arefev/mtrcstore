package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/http/middleware"
	"github.com/arefev/mtrcstore/internal/server/service"
)

func UpdateHandler(mux *http.ServeMux) {
	mux.Handle("/update/", middleware.Post(http.HandlerFunc(update)))
}

func update(w http.ResponseWriter, r *http.Request) {

	parser := service.NewUrlParser(r.URL)
	params, err := parser.Exec()

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUrl):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	fmt.Printf("Params is %+v\n", params)
	w.Write([]byte("Metrics are updated!"))
}
