package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AsaHero/whereismycity/delivery/api"
	"github.com/AsaHero/whereismycity/delivery/api/dto/models"
	"github.com/AsaHero/whereismycity/delivery/api/handlers"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/embeddings"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/repository/locations"
	users_repo "github.com/AsaHero/whereismycity/internal/infrasturcture/repository/users"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/transliterator"
	"github.com/AsaHero/whereismycity/internal/infrasturcture/typesense"
	"github.com/AsaHero/whereismycity/internal/service/auth"
	"github.com/AsaHero/whereismycity/internal/service/search"
	"github.com/AsaHero/whereismycity/internal/service/users"
	"github.com/AsaHero/whereismycity/pkg/bot"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/AsaHero/whereismycity/pkg/database/postgres"
	"github.com/AsaHero/whereismycity/pkg/logger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type App struct {
	server *http.Server
	config *config.Config
	logger *logrus.Logger
	bot    *bot.Bot
	db     *gorm.DB
}

func New(cfg *config.Config) (*App, error) {
	// Init logger
	logger := logger.Init(cfg, cfg.APP+".log")

	// Init database
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	// Init bot
	bot, err := bot.New(cfg)
	if err != nil {
		return nil, err
	}

	// Inin redis, kafka, grpc, rebbitmq, etc.

	return &App{
		config: cfg,
		logger: logger,
		bot:    bot,
		db:     db,
	}, nil
}

func (a *App) Start() error {
	// Init context timeout duration
	comtextDuration, err := time.ParseDuration(a.config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("failed to parse timeout duration: %w", err)
	}

	// Inin embeddings client
	embeddingsClient, err := embeddings.New(a.config)
	if err != nil {
		return fmt.Errorf("failed to init embeddings client: %w", err)
	}

	// Inin typesense clinet
	typesenseClient, err := typesense.New(a.config)
	if err != nil {
		return fmt.Errorf("failed to init typesense client: %w", err)
	}

	// Init transliterator client
	transliteratorClient, err := transliterator.New(a.config)
	if err != nil {
		return fmt.Errorf("failed to init transliterator client: %w", err)
	}

	// Init repo
	userRepo := users_repo.New(a.db)
	locationsRepo := locations.New(a.db)

	// Init service
	authService := auth.New(comtextDuration, userRepo)
	userService := users.New(comtextDuration, userRepo)
	searchService := search.New(comtextDuration, locationsRepo, embeddingsClient, typesenseClient, transliteratorClient)

	// Init gin router
	apiRouter := api.NewRouter(a.config, &handlers.HandlerOptions{
		Bot:           a.bot,
		AuthService:   authService,
		UserService:   userService,
		SearchService: searchService,
	})

	a.bot.SendContacts(context.Background(), models.SendContactsRequest{
		Name:    "Asa",
		Email:   "6zDz6@example.com",
		Company: "Whereismycity",
		Message: "Hello",
	})

	// Init http server
	a.server, err = api.NewServer(a.config, apiRouter)
	if err != nil {
		return fmt.Errorf("failed to init http server: %w", err)
	}
	fmt.Println("Listen: ", "address", a.config.Server.Host+a.config.Server.Port)
	return a.server.ListenAndServe()
}

func (a *App) Stop() {
	a.server.Shutdown(context.Background())

	sqlDB, _ := a.db.DB()

	sqlDB.Close()

	a.logger.Writer().Close()
}
