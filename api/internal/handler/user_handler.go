package handler

import (
	"log"
	"net/http"
)

func (srv Service) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := GetInfoFromCookie(r, srv.AccessTokenName, srv.cfg.JWTKey)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in DeleteUserHandler in GetInfoFromCookie: %v ", err)
		return
	}
	sql := `DELETE FROM users WHERE id = $1`
	_, err = srv.pool.Exec(r.Context(), sql, body.Id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("\n Error in DeleteUserHandler in Exec: %v ", err)
		return
	}
	w.Write([]byte("User Deleted Successfully"))
}
