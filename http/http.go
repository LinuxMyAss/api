package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"api"
	"api/logger"
	"http/routes"
	"net/http"
	"strings"
	"time"
)

// Config represents a configuration for a Server instance
type Config struct {
	Address string `json:"address"`
}

// Server represents a HTTP server
type Server struct {
	Config  Config
	Library *api.Library
	router  *chi.Mux
}

// New creates a new Server instance
func New(config Config, lib *api.Library) *Server {
	router := chi.NewRouter()

	corsOptions := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	router.Use(corsOptions.Handler)

	// Add our custom logger.
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			defer func() {
				var remoteAddr string

				if len(r.Header.Get("X-Forwarded-For")) > 0 {
					remoteAddr = r.Header.Get("X-Forwarded-For")
				} else {
					remoteAddr = r.RemoteAddr[0:strings.Index(r.RemoteAddr, ":")]
				}

				logger.Infof("[HTTP] %s - %s %s %d (%v)", remoteAddr, r.Method, r.RequestURI, ww.Status(), time.Now().Sub(start))
			}()

			next.ServeHTTP(ww, r)
		})
	})

	// Add the "GET /" route.
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{}"))
		if err != nil {
			logger.Errorw("[HTTP] Failed to write response.", logger.Err(err))
		}
	})

	// Add the "GET /user" route.
	routes.User(router, lib)
	// Add the "POST /user/login" route.
	routes.UserLogin(router, lib)
	// Add the "POST /user/reset" route.
	routes.UserReset(router, lib)
	// Add the "POST /user/register" route.
	routes.UserRegister(router, lib)
	// Add the "POST /user/register/confirm" route.
	routes.UserRegisterConfirm(router, lib)
	// Add the "GET /user/info" route.
	routes.UserInfo(router, lib)
	// Add the "GET /user/staff" route.
	routes.UserStaff(router, lib)
	// Add the "GET /user/{id}" route.
	routes.UserID(router, lib)
	// Add the "POST /user" route.
	routes.UserCreate(router, lib)
	// Add the "PUT /user/{id}" route.
	routes.UserUpdate(router, lib)

	// Add the "GET /token" route.
	routes.Token(router, lib)
	// Add the "GET /token/{id}" route.
	routes.TokenGet(router, lib)
	// Add the "DELETE /token/{id}" route.
	routes.TokenDelete(router, lib)

	// Add the "GET /group" route.
	routes.Group(router, lib)
	// Add the "GET /group/{id}" route.
	routes.GroupID(router, lib)
	// Add the "POST /group" route.
	routes.GroupCreate(router, lib)
	// Add the "PUT /group/{id}" route.
	routes.GroupUpdate(router, lib)
	// Add the "DELETE /group/{id}" route.
	routes.GroupDelete(router, lib)

	// Add the "GET /punishment/{id}" route.
	routes.PunishmentID(router, lib)
	// Add the "POST /punishment" route.
	routes.PunishmentCreate(router, lib)

	// Return a new Server instance.
	return &Server{
		Config:  config,
		Library: lib,
		router:  router,
	}
}

// Start starts the HTTP server
func (server *Server) Start() error {
	return http.ListenAndServe(server.Config.Address, server.router)
}
