package router

import (
	"net/http"

	eventHandler "rawuh-service/internal/event/handler"
	guestHandler "rawuh-service/internal/guest/handler"

	"github.com/gorilla/mux"
)

func NewRouter(g *guestHandler.GuestHandler, e *eventHandler.EventHandler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/events/{event_id}/guests/list", g.ListGuests).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/guests", g.AddGuest).Methods(http.MethodPost)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", g.UpdateGuestByID).Methods(http.MethodPut)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", g.GetGuestByID).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/guests/{guest_id}", g.DeleteGuestByID).Methods(http.MethodDelete)

	r.HandleFunc("/events", e.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/events/list/{user_id}", e.ListEvent).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/{user_id}", e.DetailEvent).Methods(http.MethodGet)
	r.HandleFunc("/events/{event_id}/{user_id}", e.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/events/{event_id}/{user_id}", e.DeleteEvent).Methods(http.MethodDelete)

	return r
}
