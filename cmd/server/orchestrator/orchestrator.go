package orchestrator

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/middleware"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/application"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/http"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/infrastructure/persistence"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/helper/auth"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/storage"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Symphony struct {
	Auth       *http.AuthHandler
	User       *http.UserHandler
	Middleware *middleware.Middleware
}

type Orchestrator struct {
	cfg          config.Config
	log          *zap.SugaredLogger
	db           *sqlx.DB
	cache        storage.CacheStorage
	tokenManager *auth.TokenManager
}

func New(cfg config.Config, log *zap.SugaredLogger, db *sqlx.DB, redis storage.CacheStorage, tokenManager *auth.TokenManager) *Orchestrator {
	return &Orchestrator{
		cfg:          cfg,
		log:          log,
		db:           db,
		cache:        redis,
		tokenManager: tokenManager,
	}
}

func (o *Orchestrator) Compose() *Symphony {
	// Persistence
	userPersistence := persistence.NewUserPersistence(o.db)
	notificationPreferencesPersistence := persistence.NewNotificationPreferencesPersistence(o.db)
	jobProfilePersistence := persistence.NewJobProfilePersistence(o.db)
	// Application
	authService := application.NewAuthService(o.log, o.tokenManager, userPersistence)
	userService := application.NewUserService(o.log, userPersistence, notificationPreferencesPersistence, jobProfilePersistence)
	// HTTP
	authHandler := http.NewAuthHandler(o.log, o.cache, authService, userService)
	userHandler := http.NewUserHandler(o.log, userService)

	// Middleware
	middleware := middleware.NewMiddleware(userPersistence, o.tokenManager)

	return &Symphony{
		Auth:       authHandler,
		User:       userHandler,
		Middleware: middleware,
	}
}
