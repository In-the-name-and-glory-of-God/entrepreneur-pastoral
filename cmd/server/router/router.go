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
			r.Put("/register", srv.symphony.Auth.Register)
			r.Post("/login", srv.symphony.Auth.Login)
			r.Post("/password/reset", srv.symphony.Auth.ResetPassword)
			r.Patch("/password/reset/{id}/{token}", srv.symphony.Auth.ResetPassword)
			r.Patch("/email/verify/{token}", srv.symphony.Auth.VerifyEmail)
			r.Get("/refresh", srv.symphony.Auth.Refresh)
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(srv.symphony.Middleware.Authenticate)
			r.Get("/{id}", srv.symphony.User.GetByID)
			r.Put("/{id}", srv.symphony.User.Update)
			r.Post("/list", srv.symphony.User.List)
			// Flags handlers
			r.Route("/{id}/flag", func(r chi.Router) {
				r.Patch("/active", srv.symphony.User.SetIsActive)
				r.Patch("/catholic", srv.symphony.User.SetIsCatholic)
				r.Patch("/entrepreneur", srv.symphony.User.SetIsEntrepreneur)
			})
		})

		r.Route("/entrepreneur", func(r chi.Router) {
			r.Route("/business", func(r chi.Router) {
				// Public routes
				r.Post("/list", srv.symphony.Business.List)
				r.Get("/{id}", srv.symphony.Business.GetByID)

				// Authenticated routes
				r.Group(func(r chi.Router) {
					r.Use(srv.symphony.Middleware.Authenticate)
					r.Use(srv.symphony.Middleware.UserIsCatholic)
					r.Use(srv.symphony.Middleware.UserIsEntrepreneur)
					r.Post("/", srv.symphony.Business.Create)
					r.Put("/{id}", srv.symphony.Business.Update)
					r.Delete("/{id}", srv.symphony.Business.Delete)
				})
			})

			r.Route("/product", func(r chi.Router) {
				// Public routes
				r.Post("/list", srv.symphony.Product.List)
				r.Get("/{id}", srv.symphony.Product.GetByID)

				// Authenticated routes
				r.Group(func(r chi.Router) {
					r.Use(srv.symphony.Middleware.Authenticate)
					r.Use(srv.symphony.Middleware.UserIsCatholic)
					r.Use(srv.symphony.Middleware.UserIsEntrepreneur)
					r.Post("/", srv.symphony.Product.Create)
					r.Put("/{id}", srv.symphony.Product.Update)
					r.Delete("/{id}", srv.symphony.Product.Delete)
				})
			})

			r.Route("/service", func(r chi.Router) {
				// Public routes
				r.Post("/list", srv.symphony.Service.List)
				r.Get("/{id}", srv.symphony.Service.GetByID)

				// Authenticated routes
				r.Group(func(r chi.Router) {
					r.Use(srv.symphony.Middleware.Authenticate)
					r.Use(srv.symphony.Middleware.UserIsCatholic)
					r.Use(srv.symphony.Middleware.UserIsEntrepreneur)
					r.Post("/", srv.symphony.Service.Create)
					r.Put("/{id}", srv.symphony.Service.Update)
					r.Delete("/{id}", srv.symphony.Service.Delete)
				})
			})

			r.Route("/job", func(r chi.Router) {
				r.Use(srv.symphony.Middleware.Authenticate)
				r.Use(srv.symphony.Middleware.UserIsCatholic)
				r.Post("/list", srv.symphony.Job.List)
				r.Get("/{id}", srv.symphony.Job.GetByID)

				r.Group(func(r chi.Router) {
					r.Use(srv.symphony.Middleware.UserIsEntrepreneur)
					r.Post("/", srv.symphony.Job.Create)
					r.Put("/{id}", srv.symphony.Job.Update)
					r.Delete("/{id}", srv.symphony.Job.Delete)
				})
			})
		})
	})

	return r
}
