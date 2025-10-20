package router

import (
	"net/http"

	guestHandler "rawuh-service/internal/guest/handler"

	"github.com/gorilla/mux"
)

func NewRouter(h *guestHandler.GuestHandler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/events/{event_id}/guests", h.ListGuests).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/guests", h.AddGuest).Methods(http.MethodPost)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", h.UpdateGuestByID).Methods(http.MethodPut)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", h.GetGuestByID).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", h.GetGuestByID).Methods(http.MethodDelete)

	// r.HandleFunc("/guest/{event_id}/list", h.ListGuests).Methods(http.MethodGet)
	// r.HandleFunc("/guest/{event_id}/create", h.AddGuest).Methods(http.MethodPost)
	// r.HandleFunc("/guest/{event_id}/update/{guest_id}", h.UpdateGuestByID).Methods(http.MethodPost)
	// r.HandleFunc("/guest/{event_id}/{guest_id}", h.UpdateGuestByID).Methods(http.MethodPost)
	return r
}
