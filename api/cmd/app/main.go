package main

import (
	"api/voyago/internal/config"
	"api/voyago/internal/db"
	"api/voyago/internal/handler"
	"context"
	"github.com/rs/cors"
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
	srv := handler.NewService(pool, cfg)
	handler.RegisterRoutes(mux, srv)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler(mux)

	log.Println("Server running on port", cfg.Port)
	log.Fatal(http.ListenAndServe("localhost:"+cfg.Port, corsHandler))
}
