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
	handler.RegisterRoutes(mux)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
