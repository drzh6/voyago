package main

import (
	"api/voyago/internal/config"
	"api/voyago/internal/db"
	"api/voyago/internal/handler"
	"context"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv := handler.NewService(pool)
	handler.RegisterRoutes(mux, srv)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe("localhost:"+cfg.Port, mux))
}
