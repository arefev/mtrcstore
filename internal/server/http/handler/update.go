package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arefev/mtrcstore/internal/server/http/middleware"
	"github.com/arefev/mtrcstore/internal/server/repository"
	"github.com/arefev/mtrcstore/internal/server/service"
)

type UpdateHandler struct {
	Storage repository.Storage
}

func (h *UpdateHandler) Handle(mux *http.ServeMux) {
	mux.Handle("/update/", middleware.Post(http.HandlerFunc(h.update)))
}

func (h *UpdateHandler) update(w http.ResponseWriter, r *http.Request) {

	parser := service.NewUrlParser(r.URL)
	p, err := parser.Exec()

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidUrl):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	h.Storage.Save(p.Type, p.Name, p.Value)

	fmt.Printf("Params is %+v\n", p)
	fmt.Printf("Storage has %+v\n", h.Storage)
	w.Write([]byte("Metrics are updated!"))
}
