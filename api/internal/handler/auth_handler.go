package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type AuthHandlerRequest struct {
	Username string `json:"username"`
}

func (srv *Service) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req AuthHandlerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(req.Username)
	w.Write([]byte(req.Username))
}
