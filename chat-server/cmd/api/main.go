package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/vavelour/chat/configs"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/repository/inmemorydb"
	"github.com/vavelour/chat/internal/repository/inmemorydb/repos"
	"github.com/vavelour/chat/internal/repository/postgres"
	repossql "github.com/vavelour/chat/internal/repository/postgres/repos"
	postgresdb "github.com/vavelour/chat/pkg/database_utils/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vavelour/chat/docs"
	"github.com/vavelour/chat/internal/handler"
	"github.com/vavelour/chat/internal/handler/middlewares"
	"github.com/vavelour/chat/internal/service"
	"github.com/vavelour/chat/pkg/http_utils/server"
)

const (
	basePath = "/chat/api"
)

type AuthRepository interface {
	InsertUser(username string, password string) error
	GetUser(username string) (entities.User, error)
}

type PublicRepository interface {
	InsertMessage(m entities.Message) error
	GetMessages(limit, offset int) ([]entities.Message, error)
}

type PrivateRepository interface {
	InsertMessage(m entities.Message) error
	GetMessages(sender, recipient string, limit, offset int) ([]entities.Message, error)
	GetUsers(user string) ([]string, error)
}

type AuthService interface {
	CreateUser(username, password string) (string, error)
	UserIdentity(usr interface{}) (string, error)
}

type IdentityService interface {
	Identify(next http.Handler) http.Handler
}

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Messenger

//	@host		localhost:8080
//	@BasePath	/chat/api

// @securityDefinitions.basic	BasicAuth
// @in							header
// @name						Authorization
func main() {
	cfg, err := configs.InitConfig()
	if err != nil {
		log.Println(err)
		return
	}

	var (
		authRepo     AuthRepository
		publicRepo   PublicRepository
		privateRepo  PrivateRepository
		authService  AuthService
		userIdentity IdentityService
		logInMW      func(next http.Handler) http.Handler
	)

	switch cfg.DB.Type {
	case "in_memory_db":
		db := inmemorydb.NewDB()
		authRepo = repos.NewAuthRepos(db)
		publicRepo = repos.NewPublicRepos(db)
		privateRepo = repos.NewPrivateRepos(db)
	case "postgres":
		db, err := postgres.NewSqlPostgresDB(postgresdb.SqlPostgresConfig{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			DBName:   cfg.DB.DBName,
			Password: cfg.DB.Password,
			SSLMode:  cfg.DB.SSLMode})
		if err != nil {
			log.Println(err)
			return
		}
		authRepo = repossql.NewAuthSqlRepos(db)
		publicRepo = repossql.NewPublicSqlRepos(db)
		privateRepo = repossql.NewPrivateSqlRepos(db)
	default:
		log.Println("в конфиге написана хуйня")
		return
	}

	validate := validator.New()

	switch cfg.Auth.Type {
	case "basic_auth":
		authService = service.NewAuthService(authRepo)
		userIdentity = middlewares.NewBasicUserIdentity(authService, validate)
		logInMW = userIdentity.Identify
	case "bearer_jwt":
		authService = service.NewJWTService(authRepo)
		userIdentity = middlewares.NewJWTUserIdentity(authService, validate)
		logInMW = userIdentity.Identify
	default:
		log.Println("в конфиге написана хуйня")
		return
	}

	authHandler := handler.NewAuthHandler(authService, validate)

	publicService := service.NewPublicService(publicRepo)
	publicHandler := handler.NewPublicHandler(publicService, validate)

	privateService := service.NewPrivateService(privateRepo)
	privateHandler := handler.NewPrivateHAndler(privateService, validate)

	mainRouter := chi.NewRouter()

	authHandler.AuthRoutes(mainRouter, middlewares.MyLogger, middlewares.MyRecoverer)
	publicHandler.PublicRoutes(mainRouter, logInMW, middlewares.MyLogger, middlewares.MyRecoverer)
	privateHandler.PrivateRoutes(mainRouter, logInMW, middlewares.MyLogger, middlewares.MyRecoverer)
	mainRouter.Get("/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	mainRouter.Mount(basePath, mainRouter)

	srv := server.NewServer(server.HttpServerConfig{
		Addr:           cfg.Server.Addr,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
	}, mainRouter)
	go func() {
		err := srv.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal()
		}
	}()

	log.Printf("Server start at port: %s", cfg.Server.Addr)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	if err := srv.Shutdown(context.Background()); err != nil {
		fmt.Printf("server stop: %s", err)
	}

	log.Println("Server stop.")
}
