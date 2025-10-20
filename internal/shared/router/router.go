package router

import (
	"net/http"

	guestHandler "rawuh-service/internal/guest/handler"

	"github.com/gorilla/mux"
)

func NewRouter(h *guestHandler.GuestHandler) http.Handler {
	r := mux.NewRouter()
	// r.HandleFunc("/product/create", h.AddProduct).Methods(http.MethodPost)
	// r.HandleFunc("/product/list", h.ListProducts).Methods(http.MethodGet)
	r.HandleFunc("/guest/list", h.ListGuests).Methods(http.MethodGet)
	r.HandleFunc("/guest/create", h.AddGuest).Methods(http.MethodPost)
	return r
}
