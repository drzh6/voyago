package handler

import (
	"encoding/json"
	"fmt"
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

func (srv *Service) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = srv.pool.Exec(r.Context(), `INSERT INTO users (id, login, email, password, name, surname) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)`, req.Login, req.Email, req.Password, req.Name, req.Surname)
	if err != nil {
		log.Println("failed to insert user:", err)
		return
	}

	fmt.Println(req.Login)
	w.Write([]byte(req.Login))
}
