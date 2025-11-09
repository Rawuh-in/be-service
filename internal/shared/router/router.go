package router

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	authHandler "rawuh-service/internal/auth/handler"
	eventHandler "rawuh-service/internal/event/handler"
	guestHandler "rawuh-service/internal/guest/handler"
	projectHandler "rawuh-service/internal/project/handler"
	"rawuh-service/internal/shared/middleware"
	redisPkg "rawuh-service/internal/shared/redis"
	userHandler "rawuh-service/internal/user/handler"

	docs "rawuh-service/docs"

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

	// Allow overriding the swagger spec URL via env var SWAGGER_URL
	// default points to the local handler below: /swagger/doc.json
	swaggerURL := utils.GetEnv("SWAGGER_URL", "/swagger/doc.json")

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

		// If SWAGGER_URL is an absolute URL (e.g. https://myhost/path/swagger/doc.json)
		// use it to override host/basePath/schemes so the UI examples point to the right public API.
		swaggerURLEnv := swaggerURL
		swaggerIsAbs := false
		if u, err := url.Parse(swaggerURLEnv); err == nil && u.Scheme != "" && u.Host != "" {
			swaggerIsAbs = true
			docObj["host"] = u.Host
			// If the swagger doc path is nested (e.g. /api/swagger/doc.json), set basePath to the prefix
			bp := strings.TrimSuffix(u.Path, "/swagger/doc.json")
			if bp == "" {
				bp = "/"
			}
			docObj["basePath"] = bp
			// If SWAGGER_SCHEMES not set explicitly, infer from URL scheme
			if schemesEnv := utils.GetEnv("SWAGGER_SCHEMES", ""); schemesEnv == "" && u.Scheme != "" {
				docObj["schemes"] = []string{u.Scheme}
			}
		}

		// Only apply docs.SwaggerInfo overrides when swagger URL isn't an absolute external URL
		if !swaggerIsAbs && docs.SwaggerInfo != nil {
			if hostVal := docs.SwaggerInfo.Host; hostVal != "" {
				host := strings.TrimSpace(hostVal)
				// If host contains a scheme (user supplied it incorrectly), parse and extract host
				if strings.Contains(host, "://") {
					if u, err := url.Parse(host); err == nil {
						if u.Host != "" {
							host = u.Host
						}
						// If SWAGGER_SCHEMES not provided explicitly, infer from the parsed URL
						if schemesEnv := utils.GetEnv("SWAGGER_SCHEMES", ""); schemesEnv == "" && u.Scheme != "" {
							docObj["schemes"] = []string{u.Scheme}
						}
					} else {
						// fallback: strip common prefixes
						host = strings.ReplaceAll(host, "http://", "")
						host = strings.ReplaceAll(host, "https://", "")
						if idx := strings.Index(host, "://"); idx != -1 {
							host = host[idx+3:]
						}
					}
				}
				docObj["host"] = host
			}
			if bp := docs.SwaggerInfo.BasePath; bp != "" {
				docObj["basePath"] = bp
			}
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

		// Optionally override schemes (e.g. https) via SWAGGER_SCHEMES env (comma-separated)
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
		}

		out, err := json.Marshal(docObj)
		if err != nil {
			_, _ = w.Write([]byte(doc))
			return
		}
		_, _ = w.Write(out)
	})
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(swaggerURL)))

	return r
}
