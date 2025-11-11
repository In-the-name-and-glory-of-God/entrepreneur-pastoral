package router

import (
	"net/http"
	"time"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/orchestrator"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httprateredis "github.com/go-chi/httprate-redis"
	"github.com/redis/go-redis/v9"
)

type ServerRouter struct {
	cfg      config.Config
	symphony *orchestrator.Symphony
}

func NewServerRouter(cfg config.Config, symphony *orchestrator.Symphony) *ServerRouter {
	return &ServerRouter{
		cfg:      cfg,
		symphony: symphony,
	}
}

func (srv *ServerRouter) Mount(client *redis.Client) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))
	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   srv.cfg.API.CORS.AllowedOrigins,
		AllowedMethods:   srv.cfg.API.CORS.AllowedMethods,
		AllowedHeaders:   srv.cfg.API.CORS.AllowedHeaders,
		ExposedHeaders:   srv.cfg.API.CORS.ExposedHeaders,
		AllowCredentials: srv.cfg.API.CORS.AllowCredentials,
		MaxAge:           srv.cfg.API.CORS.MaxAge, // Maximum value not ignored by any of major browsers
	}))
	// Rate Limiting: 10 Reqs/Min for each endpoint.
	// IP stored in Redis with key httprate:[ipaddr]
	r.Use(httprate.Limit(
		srv.cfg.API.RateLimiter.RequestLimit,
		srv.cfg.API.RateLimiter.WindowLength,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
		httprateredis.WithRedisLimitCounter(&httprateredis.Config{
			Client: client,
		}),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			srv.symphony.Auth.RegisterRoutes(r)
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(srv.symphony.Middleware.Authenticate)
			srv.symphony.User.RegisterRoutes(r)
		})
	})

	return r
}
