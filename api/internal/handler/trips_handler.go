package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	rand2 "math/rand"
	"net/http"
	"time"
)

func (srv Service) CreateTripHandler(w http.ResponseWriter, r *http.Request) {
	var req TripAddRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
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
				http.Error(w, "Trip already exists", http.StatusConflict)
				return
			}
		}
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (srv Service) GetUserListTripsHandler(w http.ResponseWriter, r *http.Request) {
	var req TripRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in GetUserListTripsHandlerGetUserListTripsHandler in GetInfoFromCookie: %v ", err)
		return
	}

	sql := `SELECT id, name, start_date, end_date, status FROM trips WHERE owner_id = $1 AND status != $2`
	rows, err := srv.pool.Query(r.Context(), sql, body.Id, "closed")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("\n Error in GetUserListTripsHandlerGetUserListTripsHandler in Query: %v ", err)
		return
	}

	jsonResult, err := RowsToJSON(rows)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("\n Error in GetUserListTripsHandlerGetUserListTripsHandler in RowsToJSON: %v ", err)
		return
	}
	w.Write(jsonResult)
}

func (srv Service) UpdateUserTripHandler(w http.ResponseWriter, r *http.Request) {
	var req TripRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in GetUserListTripsHandlerGetUserListTripsHandler in GetInfoFromCookie: %v ", err)
		return
	}

	sql := `UPDATE trips SET 
                 name = $1,
      			 description = $2,
      			 start_date = $3,
                 end_date = $4,
                 status = $5,
                 is_public = $6,
                 cover_image = $7,
                 update_at = $8
                 WHERE id = $9 AND owner_id = $10
`
	_, err = srv.pool.Exec(r.Context(), sql,
		req.Name,
		req.Description,
		req.StartDate,
		req.EndDate,
		req.Status,
		req.IsPublic,
		req.CoverImageUrl,
		time.Now(),
		req.Id,
		body.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in UpdateUserTripHandler in Exec: %v ", err)
		return
	}

	w.Write([]byte("Update succeeded"))
}

func (srv Service) GetUserTripHandler(w http.ResponseWriter, r *http.Request) {
	var req string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in GetUserTripHandler in GetInfoFromCookie: %v ", err)
		return
	}

	sql := `SELECT * FROM trips WHERE id = $1 and owner_id = $2`
	var result TripRequestBody
	err = srv.pool.QueryRow(r.Context(), sql, body.Id).Scan(result)
	if err != nil {
		log.Printf("\n Error in GetUserTripHandler in QueryRow: %v ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Printf("\n Error in GetUserTripHandler in QueryRow: %v ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.Write(jsonResult)
}

func (srv Service) DeleteUserTripHandler(w http.ResponseWriter, r *http.Request) {
	var req string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in DeleteUserTripHandler in GetInfoFromCookie: %v ", err)
		return
	}

	sql := `DELETE FROM trips WHERE id = $1 and owner_id = $2`
	_, err = srv.pool.Exec(r.Context(), sql, req, body.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in DeleteUserTripHandler in SQL: %v ", err)
		return
	}

	w.Write([]byte("Delete succeeded"))
}

func (srv Service) CompleteUserTripHandler(w http.ResponseWriter, r *http.Request) {
	var req string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in CompleteUserTripHandler in GetInfoFromCookie: %v ", err)
		return
	}

	sql := `UPDATE trips SET status = $1 WHERE id = $2 AND owner_id = $3`
	_, err = srv.pool.Exec(r.Context(), sql, "completed", req, body.Id)
}

func generateInviteCode(length int, letterRunes []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand2.Intn(len(letterRunes))]
	}
	return string(b)
}
