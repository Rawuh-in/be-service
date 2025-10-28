package router

import (
	"net/http"

	eventHandler "rawuh-service/internal/event/handler"
	guestHandler "rawuh-service/internal/guest/handler"
	projectHandler "rawuh-service/internal/project/handler"

	"github.com/gorilla/mux"
)

func NewRouter(g *guestHandler.GuestHandler, e *eventHandler.EventHandler, p *projectHandler.ProjectHandler) http.Handler {
	r := mux.NewRouter()

	// PROJECT ROUTES
	r.HandleFunc("/project/list", p.ListProject).Methods(http.MethodGet)
	r.HandleFunc("/project", p.CreateProject).Methods(http.MethodPost)
	r.HandleFunc("/project/{project_id}", p.UpdateProject).Methods(http.MethodPut)
	r.HandleFunc("/project/{project_id}", p.DeleteProject).Methods(http.MethodDelete)
	r.HandleFunc("/project/{project_id}", p.DetailProject).Methods(http.MethodGet)

	// EVENT ROUTES
	r.HandleFunc("/{project_id}/events", e.AddEvent).Methods(http.MethodPost)
	r.HandleFunc("/{project_id}/events/list", e.ListEvent).Methods(http.MethodGet)
	r.HandleFunc("/{project_id}/events/{event_id}", e.DetailEvent).Methods(http.MethodGet)
	r.HandleFunc("/{project_id}/events/{event_id}", e.UpdateEvent).Methods(http.MethodPut)
	r.HandleFunc("/{project_id}/events/{event_id}", e.DeleteEvent).Methods(http.MethodDelete)

	// GUEST ROUTES
	r.HandleFunc("/{project_id}/events/{event_id}/guests/list", g.ListGuests).Methods(http.MethodGet)
	r.HandleFunc("/{project_id}/events/{event_id}/guests", g.AddGuest).Methods(http.MethodPost)
	r.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.UpdateGuestByID).Methods(http.MethodPut)
	r.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.GetGuestByID).Methods(http.MethodGet)
	r.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.DeleteGuestByID).Methods(http.MethodDelete)

	return r
}
