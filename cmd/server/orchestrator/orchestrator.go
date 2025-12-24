package orchestrator

import (
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/cmd/server/middleware"
	adminApp "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/application"
	adminHttp "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/http"
	adminPersist "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/admin/infrastructure/persistence"
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
	// Admin handlers
	AdminUser        *adminHttp.UserHandler
	AdminBusiness    *adminHttp.BusinessHandler
	AdminChurch      *adminHttp.ChurchHandler
	AdminIndustry    *adminHttp.IndustryHandler
	AdminFieldOfWork *adminHttp.FieldOfWorkHandler
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
	// ## Admin
	addressPersistence := adminPersist.NewAddressPersistence(o.db)
	churchPersistence := adminPersist.NewChurchPersistence(o.db)
	industryPersistence := adminPersist.NewIndustryPersistence(o.db)
	fieldOfWorkPersistence := adminPersist.NewFieldOfWorkPersistence(o.db)

	// # Application
	// ## User
	authService := application.NewAuthService(o.log, o.cfg, o.cache, o.queue, o.tokenManager, userPersistence)
	userService := application.NewUserService(o.log, userPersistence, notificationPreferencesPersistence, jobProfilePersistence, addressPersistence)
	// ## Entrepreneur
	businessService := entrepreneurApp.NewBusinessService(o.log, o.cache, businessPersistence)
	productService := entrepreneurApp.NewProductService(o.log, productPersistence, businessPersistence)
	serviceService := entrepreneurApp.NewServiceService(o.log, servicePersistence, businessPersistence)
	jobService := entrepreneurApp.NewJobService(o.log, jobPersistence, businessPersistence)
	// ## Admin
	churchService := adminApp.NewChurchService(o.log, churchPersistence, addressPersistence)
	industryService := adminApp.NewIndustryService(o.log, industryPersistence)
	fieldOfWorkService := adminApp.NewFieldOfWorkService(o.log, fieldOfWorkPersistence)

	// # HTTP
	// ## User
	authHandler := http.NewAuthHandler(o.log, o.cache, authService, userService)
	userHandler := http.NewUserHandler(o.log, userService)
	// ## Entrepreneur
	businessHandler := entrepreneurHttp.NewBusinessHandler(o.log, businessService)
	productHandler := entrepreneurHttp.NewProductHandler(o.log, productService)
	serviceHandler := entrepreneurHttp.NewServiceHandler(o.log, serviceService)
	jobHandler := entrepreneurHttp.NewJobHandler(o.log, jobService)
	// ## Admin
	adminUserHandler := adminHttp.NewUserHandler(o.log, userService)
	adminBusinessHandler := adminHttp.NewBusinessHandler(o.log, businessService)
	adminChurchHandler := adminHttp.NewChurchHandler(o.log, churchService)
	adminIndustryHandler := adminHttp.NewIndustryHandler(o.log, industryService)
	adminFieldOfWorkHandler := adminHttp.NewFieldOfWorkHandler(o.log, fieldOfWorkService)

	// # Middleware
	middleware := middleware.NewMiddleware(userPersistence, o.tokenManager)

	return &Symphony{
		Auth:             authHandler,
		User:             userHandler,
		Business:         businessHandler,
		Product:          productHandler,
		Service:          serviceHandler,
		Job:              jobHandler,
		AdminUser:        adminUserHandler,
		AdminBusiness:    adminBusinessHandler,
		AdminChurch:      adminChurchHandler,
		AdminIndustry:    adminIndustryHandler,
		AdminFieldOfWork: adminFieldOfWorkHandler,
		Middleware:       middleware,
	}
}
