package handler

import (
	"api/voyago/internal/chat"
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

	mux.HandleFunc(`GET /api/trip/get_all`, srv.GetUserListTripsHandler)
	mux.HandleFunc(`GET /api/trip/{trip_id}`, srv.GetUserTripHandler)

	mux.HandleFunc(`PUT /api/trip/{trip_id}/update_trip`, srv.UpdateUserTripHandler)

	mux.HandleFunc(`PATCH /api/trip/{trip_id}/complete`, srv.CompleteUserTripHandler)

	mux.HandleFunc(`DELETE /api/trip/{trip_id}/delete`, srv.DeleteUserTripHandler)
	//User
	mux.HandleFunc(`DELETE /api/user/delete`, srv.DeleteUserHandler)

	//WebSockets
	hub := chat.NewHub()
	go hub.Run()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { chat.ServeWs(hub, w, r) })
}
