package handler

import (
	"database/sql"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/cache"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/mq"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/bcrypt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/cronjob"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/jsontag"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/jwt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logistic"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/repository"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/usecase"
	"github.com/hibiken/asynq"
)

type HandlerOpts struct {
	*UserHandler
	*AdminHandler
	*PharmacistHandler
	*LocationHandler
	*PharmacyHandler
	*CacheHandler
	*LogisticHandler
}

func Init(db *sql.DB, config *config.Config) *HandlerOpts {
	jsontag.Register()

	store := repository.NewStore(db)
	bcrypt := bcrypt.NewBcryptProvider(config.App.BcryptCost)
	jwt := jwt.NewJwtProvider(config.Jwt)
	logistic := logistic.NewRajaOngkirProvider(config.Logistic)
	crond := cronjob.NewCronProvider(config.Cron)
	go crond.Start()

	producer := mq.NewTaskProducer(asynq.RedisClientOpt{
		Addr: config.Redis.ServerAddress,
	})
	cache := cache.NewCacheProvider(config.Cache)

	userUseCase := usecase.NewUserUseCase(store, bcrypt, jwt, producer, config.Google)
	adminUseCase := usecase.NewAdminUseCase(store, bcrypt, producer, crond)
	pharmacistUseCase := usecase.NewPharmacistUseCase(store, producer)
	locationUseCase := usecase.NewLocationUseCase(store)
	pharmacyUseCase := usecase.NewPharmacyUseCase(store, logistic)
	logisticUseCase := usecase.NewLogisticUseCase(store)

	return &HandlerOpts{
		UserHandler:       NewUserHandler(userUseCase, cache, config),
		AdminHandler:      NewAdminHandler(adminUseCase),
		PharmacistHandler: NewPharmacistHandler(pharmacistUseCase),
		LocationHandler:   NewLocationHandler(locationUseCase),
		PharmacyHandler:   NewPharmacyHandler(pharmacyUseCase),
		CacheHandler:      NewCacheHandler(cache),
		LogisticHandler:   NewLogisticHandler(logisticUseCase),
	}
}
