package orchestrator

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/middleware"
	entrepreneurApp "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/application"
	entrepreneurHttp "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/http"
	entrepreneurPersist "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/entrepreneur/infrastructure/persistence"
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
	Business   *entrepreneurHttp.BusinessHandler
	Product    *entrepreneurHttp.ProductHandler
	Service    *entrepreneurHttp.ServiceHandler
	Job        *entrepreneurHttp.JobHandler
	Middleware *middleware.Middleware
}

type Orchestrator struct {
	cfg          config.Config
	log          *zap.SugaredLogger
	db           *sqlx.DB
	cache        storage.CacheStorage
	queue        storage.QueueStorage
	tokenManager *auth.TokenManager
}

func New(cfg config.Config, log *zap.SugaredLogger, db *sqlx.DB, redis storage.CacheStorage, queue storage.QueueStorage, tokenManager *auth.TokenManager) *Orchestrator {
	return &Orchestrator{
		cfg:          cfg,
		log:          log,
		db:           db,
		cache:        redis,
		queue:        queue,
		tokenManager: tokenManager,
	}
}

func (o *Orchestrator) Compose() *Symphony {
	// # Persistence
	// ## User
	userPersistence := persistence.NewUserPersistence(o.db)
	notificationPreferencesPersistence := persistence.NewNotificationPreferencesPersistence(o.db)
	jobProfilePersistence := persistence.NewJobProfilePersistence(o.db)
	// ## Entrepreneur
	businessPersistence := entrepreneurPersist.NewBusinessPersistence(o.db)
	productPersistence := entrepreneurPersist.NewProductPersistence(o.db)
	servicePersistence := entrepreneurPersist.NewServicePersistence(o.db)
	jobPersistence := entrepreneurPersist.NewJobPersistence(o.db)

	// # Application
	// ## User
	authService := application.NewAuthService(o.log, o.tokenManager, userPersistence)
	userService := application.NewUserService(o.log, userPersistence, notificationPreferencesPersistence, jobProfilePersistence)
	// ## Entrepreneur
	businessService := entrepreneurApp.NewBusinessService(o.log, businessPersistence)
	productService := entrepreneurApp.NewProductService(o.log, productPersistence, businessPersistence)
	serviceService := entrepreneurApp.NewServiceService(o.log, servicePersistence, businessPersistence)
	jobService := entrepreneurApp.NewJobService(o.log, jobPersistence, businessPersistence)

	// # HTTP
	// ## User
	authHandler := http.NewAuthHandler(o.log, o.cache, authService, userService)
	userHandler := http.NewUserHandler(o.log, userService)
	// ## Entrepreneur
	businessHandler := entrepreneurHttp.NewBusinessHandler(o.log, businessService)
	productHandler := entrepreneurHttp.NewProductHandler(o.log, productService)
	serviceHandler := entrepreneurHttp.NewServiceHandler(o.log, serviceService)
	jobHandler := entrepreneurHttp.NewJobHandler(o.log, jobService)

	// # Middleware
	middleware := middleware.NewMiddleware(userPersistence, o.tokenManager)

	return &Symphony{
		Auth:       authHandler,
		User:       userHandler,
		Business:   businessHandler,
		Product:    productHandler,
		Service:    serviceHandler,
		Job:        jobHandler,
		Middleware: middleware,
	}
}
