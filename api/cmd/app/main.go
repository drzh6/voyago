package main

import (
	"api/voyago/internal/config"
	"api/voyago/internal/handler"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	srv := handler.NewService()
	handler.RegisterRoutes(mux, srv)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe("localhost:"+cfg.Port, mux))
}
