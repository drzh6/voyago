package handler

import (
	"api/voyago/internal/config"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	Password []byte `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type LoginHandlerRequest struct {
	Login    string `json:"login"`
	Password []byte `json:"password"`
}

func (srv *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	hashPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = srv.pool.Exec(r.Context(), `INSERT INTO users (id, login, email, password, name, surname) VALUES ($1, $2, $3, $4, $5, $6)`, id, req.Login, req.Email, hashPassword, req.Name, req.Surname)
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

	err = CreateTokens(w, id, srv.cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println(req.Login)
	w.Write([]byte("Register success!"))
}

func (srv *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var id uuid.UUID
	var name string
	var hashedPassword []byte

	err = srv.pool.QueryRow(r.Context(), `SELECT id, name, password FROM users WHERE login = $1`, req.Login).Scan(&id, &name, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Invalid user", http.StatusUnauthorized)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = CreateTokens(w, id, srv.cfg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(name))
}

func HashPassword(password []byte) ([]byte, error) {
	if password == nil || len(password) == 0 {
		return nil, errors.New("password is required")
	}
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return bytes, err
}

func CreateTokens(w http.ResponseWriter, id uuid.UUID, config config.Config) error {
	token, err := CreateToken(id, config.JWTKey)
	if err != nil {
		return err
	}

	refreshToken, err := CreateToken(id, config.JWTRefreshKey)
	if err != nil {
		return err
	}

	JWTAccessTime, err := strconv.Atoi(config.JWTAccessTime)
	if err != nil {
		return err
	}

	JWTRefreshTime, err := strconv.Atoi(config.JWTRefreshTime)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTAccessTime) * time.Minute),
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/api/refresh",
		HttpOnly: true,
		Secure:   false, // True в проде (Https)
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(JWTRefreshTime) * time.Hour * 24),
	})
	return nil
}

func CreateToken(id uuid.UUID, key string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})
	return token.SignedString(key)
}
