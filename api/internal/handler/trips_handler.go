package handler

import (
	"api/voyago/internal/service"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	rand2 "math/rand"
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

	body, err := service.GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	inviteCode := generateInviteCode(srv.cfg.InviteCodeLength, srv.cfg.InviteCodeRunes)
	sql := `INSERT INTO trips (id, name, description, owner_id, start_date, end_date, invite_code) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = srv.pool.Exec(r.Context(), sql, uuid.New(), req.Name, req.Description, body.Id, req.StartDate, req.EndDate, inviteCode)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
		}
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func generateInviteCode(length int, letterRunes []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand2.Intn(len(letterRunes))]
	}
	return string(b)
}
