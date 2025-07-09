package handler

import (
	"fmt"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, srv *Service) {
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})
	//Auth
	mux.HandleFunc("POST /api/user/registration", srv.RegisterHandler)
	mux.HandleFunc("POST /api/user/login", srv.LoginHandler)
	//Trips
	mux.HandleFunc(`POST /api/trip/add`, srv.CreateTripHandler)

	mux.HandleFunc(`GET /api/trip`, srv.GetUserListTripsHandler)
	mux.HandleFunc(`GET /api/trip/get_all`, srv.GetUserTripHandler)

	mux.HandleFunc(`UPDATE /api/trip/update_trip`, srv.UpdateUserTripHandler)
	mux.HandleFunc(`UPDATE /api/trip/complete`, srv.CompleteUserTripHandler)

	mux.HandleFunc(`DELETE /api/trip/delete`, srv.DeleteUserTripHandler)
	//User
	mux.HandleFunc(`POST /api/user/delete`, srv.DeleteUserTripHandler)
}
