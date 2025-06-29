package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type TripAddRequestBody struct {
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	IsPublic      bool      `json:"is_public"`
	CoverImageUrl string    `json:"cover_image_url"`
}

func (srv Service) AddTripHandler(w http.ResponseWriter, r *http.Request) {
	var req TripAddRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = srv.pool.Exec(r.Context(), `INSERT INTO trips (
                   name,
                   description,
                   start_date,
                   end_date,
                   status,
                   is_public,
                   cover_image_url,)`)
}
