package handler

import (
	"api/voyago/internal/config"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
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

	accessTokenCookie, refreshTokenCookie, err := CreateTokens(id, srv.cfg)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Token Error: " + err.Error())
		return
	}

	_, err = srv.pool.Exec(r.Context(), `INSERT INTO "users" (id, login, email, password, name, surname) VALUES ($1, $2, $3, $4, $5, $6)`, id, req.Login, req.Email, hashPassword, req.Name, req.Surname)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" /*EBAL ROT*/ {
				http.Error(w, "User already exists", http.StatusConflict)
				return
			}
		}
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
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

	accessTokenCookie, refreshTokenCookie, err := CreateTokens(id, srv.cfg)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Token Error:" + err.Error())
		return
	}

	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
	w.Write([]byte(name))
}

func HashPassword(password []byte) ([]byte, error) {
	if password == nil || len(password) == 0 {
		return nil, errors.New("password is required")
	}
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return bytes, err
}

func CreateTokens(id uuid.UUID, config config.Config) (*http.Cookie, *http.Cookie, error) {
	token, err := CreateToken(id, []byte(config.JWTRefreshKey))
	if err != nil {
		log.Println("Create Token Error: " + err.Error())
		return nil, nil, err
	}

	refreshToken, err := CreateToken(id, []byte(config.JWTRefreshKey))
	if err != nil {
		log.Println("Create Token Error: " + err.Error())
		return nil, nil, err
	}

	JWTAccessTime, err := strconv.Atoi(config.JWTAccessTime)
	if err != nil {
		return nil, nil, err
	}

	JWTRefreshTime, err := strconv.Atoi(config.JWTRefreshTime)
	if err != nil {
		return nil, nil, err
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTAccessTime) * time.Minute),
	}

	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/refresh",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTRefreshTime) * time.Hour * 24),
	}
	return accessCookie, refreshCookie, nil
}

func CreateToken(id uuid.UUID, key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	return token.SignedString(key)
}
