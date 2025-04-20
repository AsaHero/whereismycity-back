package handlers

import (
	"github.com/AsaHero/whereismycity/delivery/api/validation"
	"github.com/AsaHero/whereismycity/internal/service/auth"
	"github.com/AsaHero/whereismycity/internal/service/search"
	"github.com/AsaHero/whereismycity/internal/service/users"
	"github.com/AsaHero/whereismycity/pkg/bot"
	"github.com/AsaHero/whereismycity/pkg/config"
)

type HandlerOptions struct {
	Bot           *bot.Bot
	AuthService   auth.AuthService
	UserService   users.Service
	SearchService search.Service
}

type Handler struct {
	bot           *bot.Bot
	config        *config.Config
	validator     *validation.Validator
	searchService search.Service
	userService   users.Service
	authService   auth.AuthService
}

func New(cfg *config.Config, validator *validation.Validator, opt *HandlerOptions) *Handler {
	return &Handler{
		bot:           opt.Bot,
		config:        cfg,
		validator:     validator,
		searchService: opt.SearchService,
		userService:   opt.UserService,
		authService:   opt.AuthService,
	}
}
