package handler

import (
	"fmt"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, s *Service) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("POST /api/registration", s.RegisterHandler)
}
