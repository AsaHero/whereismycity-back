package handlers

import (
	"github.com/AsaHero/whereismycity/delivery/api/validation"
	"github.com/AsaHero/whereismycity/internal/service/auth"
	"github.com/AsaHero/whereismycity/internal/service/search"
	"github.com/AsaHero/whereismycity/pkg/config"
)

type HandlerOptions struct {
	AuthService   auth.AuthService
	SearchService search.Service
}

type Handler struct {
	config        *config.Config
	validator     *validation.Validator
	searchService search.Service
}

func New(cfg *config.Config, validator *validation.Validator, opt *HandlerOptions) *Handler {
	return &Handler{
		config:        cfg,
		validator:     validator,
		searchService: opt.SearchService,
	}
}
