package api

import (
	"github.com/AsaHero/whereismycity/delivery/api/docs"
	"github.com/AsaHero/whereismycity/delivery/api/handlers"
	"github.com/AsaHero/whereismycity/delivery/api/middlewares"
	"github.com/AsaHero/whereismycity/delivery/api/validation"
	"github.com/AsaHero/whereismycity/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//  @title      		Where Is My City
//  @version    		0.0.1
//  @description  		Documentation for "Where Is My City" API
//  @termsOfService  	http://swagger.io/terms/

// @securityDefinitions.basic 	BasicAuth
// @in              			header
// @name           				Authorization
// @description     			Basic Auth "Authorization: Basic <base64 encoded username:password>"

func NewRouter(cfg *config.Config, opt *handlers.HandlerOptions) *gin.Engine {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Init Validator
	validator := validation.NewValidator()

	// handlers initialization
	mainHandler := handlers.NewHandler(cfg, validator, opt)

	// Set base path /api/v1
	router := r.Group(middlewares.APIPrefix)

	// Protected routes
	api := router.Group("/", middlewares.BasicAuth(opt.AuthService))
	{
		api.GET("/search", mainHandler.Search)
	}

	// Swagger Route
	docs.SwaggerInfo.BasePath = middlewares.APIPrefix
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
