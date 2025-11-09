// AUTH ROUTES
package router

import (
	"encoding/json"
	"net/http"
	"strings"

	authHandler "rawuh-service/internal/auth/handler"
	eventHandler "rawuh-service/internal/event/handler"
	guestHandler "rawuh-service/internal/guest/handler"
	projectHandler "rawuh-service/internal/project/handler"
	"rawuh-service/internal/shared/middleware"
	redisPkg "rawuh-service/internal/shared/redis"
	userHandler "rawuh-service/internal/user/handler"

	"rawuh-service/internal/shared/lib/utils"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"

	"github.com/gorilla/mux"
)

func NewRouter(g *guestHandler.GuestHandler, e *eventHandler.EventHandler, p *projectHandler.ProjectHandler, u *userHandler.UserHandler, a *authHandler.AuthHandler, rdb *redisPkg.Redis) http.Handler {
	r := mux.NewRouter()
	// Apply CORS middleware first so preflight and headers are set globally.
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.AuthMiddleware(rdb))

	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.RequireAuth)

	// PROJECT ROUTES (protected)
	protected.HandleFunc("/project/list", p.ListProject).Methods(http.MethodGet)
	protected.HandleFunc("/project", p.CreateProject).Methods(http.MethodPost)
	protected.HandleFunc("/project/{project_id}", p.UpdateProject).Methods(http.MethodPut)
	protected.HandleFunc("/project/{project_id}", p.DeleteProject).Methods(http.MethodDelete)
	protected.HandleFunc("/project/{project_id}", p.DetailProject).Methods(http.MethodGet)

	// EVENT ROUTES (protected)
	protected.HandleFunc("/{project_id}/events", e.AddEvent).Methods(http.MethodPost)
	protected.HandleFunc("/{project_id}/events/list", e.ListEvent).Methods(http.MethodGet)
	protected.HandleFunc("/{project_id}/events/{event_id}", e.DetailEvent).Methods(http.MethodGet)
	protected.HandleFunc("/{project_id}/events/{event_id}", e.UpdateEvent).Methods(http.MethodPut)
	protected.HandleFunc("/{project_id}/events/{event_id}", e.DeleteEvent).Methods(http.MethodDelete)

	// GUEST ROUTES (protected)
	protected.HandleFunc("/{project_id}/events/{event_id}/guests/list", g.ListGuests).Methods(http.MethodGet)
	protected.HandleFunc("/{project_id}/events/{event_id}/guests", g.AddGuest).Methods(http.MethodPost)
	protected.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.UpdateGuestByID).Methods(http.MethodPut)
	protected.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.GetGuestByID).Methods(http.MethodGet)
	protected.HandleFunc("/{project_id}/events/{event_id}/guests/{guest_id}", g.DeleteGuestByID).Methods(http.MethodDelete)

	// USER ROUTES (protected)
	protected.HandleFunc("/users/list", u.ListUsers).Methods(http.MethodGet)
	protected.HandleFunc("/users", u.AddUser).Methods(http.MethodPost)
	protected.HandleFunc("/users/{user_id}", u.UpdateUserByID).Methods(http.MethodPut)
	protected.HandleFunc("/users/{user_id}", u.GetUserByID).Methods(http.MethodGet)
	protected.HandleFunc("/users/{user_id}", u.DeleteUserByID).Methods(http.MethodDelete)

	// AUTH ROUTES
	r.HandleFunc("/login", a.Login).Methods(http.MethodPost)
	protected.HandleFunc("/auth/me", a.TokenInfo).Methods(http.MethodGet)

	r.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		doc, err := swag.ReadDoc("swagger")
		if err != nil {
			http.Error(w, "failed to read swagger doc", http.StatusInternalServerError)
			return
		}
		var docObj map[string]interface{}
		if err := json.Unmarshal([]byte(doc), &docObj); err != nil {
			_, _ = w.Write([]byte(doc))
			return
		}

		// Host: prefer SWAGGER_HOST env, then X-Forwarded-Host, then request Host
		host := utils.GetEnv("SWAGGER_HOST", "")
		if host == "" {
			if xf := r.Header.Get("X-Forwarded-Host"); xf != "" {
				host = xf
			} else {
				host = r.Host
			}
		}
		docObj["host"] = host

		// Schemes: prefer SWAGGER_SCHEMES env; otherwise infer from X-Forwarded-Proto or request TLS
		if schemesEnv := utils.GetEnv("SWAGGER_SCHEMES", ""); schemesEnv != "" {
			schemes := []string{}
			for _, s := range strings.Split(schemesEnv, ",") {
				if v := strings.TrimSpace(s); v != "" {
					schemes = append(schemes, v)
				}
			}
			if len(schemes) > 0 {
				docObj["schemes"] = schemes
			}
		} else {
			scheme := r.Header.Get("X-Forwarded-Proto")
			if scheme == "" {
				if r.TLS != nil {
					scheme = "https"
				} else {
					scheme = "http"
				}
			}
			docObj["schemes"] = []string{scheme}
		}

		secDef := map[string]interface{}{
			"Bearer": map[string]interface{}{
				"type": "apiKey",
				"name": "Authorization",
				"in":   "header",
			},
		}
		docObj["securityDefinitions"] = secDef
		docObj["security"] = []interface{}{map[string]interface{}{"Bearer": []interface{}{}}}

		out, err := json.Marshal(docObj)
		if err != nil {
			_, _ = w.Write([]byte(doc))
			return
		}
		_, _ = w.Write(out)
	})
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
