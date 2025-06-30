package handler

import (
	"api/voyago/internal/service"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type RegisterHandlerRequest struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type LoginHandlerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (srv *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	hashPassword, err := HashPassword([]byte(req.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	jwtSrv := service.GetNewCookies(id, srv.cfg)
	if jwtSrv.Error != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Token Error:" + err.Error())
		return
	}

	_, err = srv.pool.Exec(r.Context(), `INSERT INTO "users" (id, login, email, password, name, surname) VALUES ($1, $2, $3, $4, $5, $6)`, id, req.Login, req.Email, hashPassword, req.Name, req.Surname)
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

	http.SetCookie(w, jwtSrv.AccessCookie)
	http.SetCookie(w, jwtSrv.RefreshCookie)
	log.Println("Register success!")
	w.Write([]byte("Register success!"))
}

func (srv *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Login Request error: " + err.Error())
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var id uuid.UUID
	var name string
	var hashedPassword []byte

	err = srv.pool.QueryRow(r.Context(), `SELECT id, name, password FROM users WHERE login = $1`, req.Login).Scan(&id, &name, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("SQL Query error: " + err.Error())
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password))
	if err != nil {
		log.Println("Hash password error" + err.Error())

		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	jwtSrv := service.GetNewCookies(id, srv.cfg)
	if jwtSrv.Error != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Token Error:" + err.Error())
		return
	}

	http.SetCookie(w, jwtSrv.AccessCookie)
	http.SetCookie(w, jwtSrv.RefreshCookie)
	w.Write([]byte(name))
}

func HashPassword(password []byte) ([]byte, error) {
	if password == nil || len(password) == 0 {
		return nil, errors.New("password is required")
	}
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return bytes, err
}
