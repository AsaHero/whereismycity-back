package api

import (
	"github.com/AsaHero/whereismycity/delivery/api/docs"
	"github.com/AsaHero/whereismycity/delivery/api/handlers"
	"github.com/AsaHero/whereismycity/delivery/api/middlewares"
	"github.com/AsaHero/whereismycity/delivery/api/validation"
	"github.com/AsaHero/whereismycity/internal/entity"
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
// @securityDefinitions.apikey 	ApiKeyAuth
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
	mainHandler := handlers.New(cfg, validator, opt)

	// Set base path /api/v1
	router := r.Group(middlewares.APIPrefix)

	// Public routes
	public := router.Group("/")
	{
		public.POST("/auth/register", mainHandler.Register)
		public.POST("/auth/login", mainHandler.Login)
		public.POST("/auth/refresh", mainHandler.RefreshToken)
		public.GET("/demo", mainHandler.Search)
		public.POST("/contacts/send", mainHandler.SendContacts)
	}

	// Bearer protected routes
	bearerProtected := router.Group("/", middlewares.BearerAuth(cfg.Token.Secret))
	{
		bearerProtected.GET("/profile", mainHandler.GetProfile)
		bearerProtected.PATCH("/profile", mainHandler.PatchProfile)
	}

	// Basic protected routes
	basicProtected := router.Group("/", middlewares.BasicAuth(opt.AuthService))
	{
		basicProtected.GET("/search", mainHandler.Search)
	}

	adminApi := router.Group("/admin", middlewares.BasicAuth(opt.AuthService), middlewares.RoleRequired(opt.AuthService, string(entity.UserRoleAdmin)))
	{
		// Users
		adminApi.POST("/users", mainHandler.CreateUser)
		adminApi.GET("/users/search", mainHandler.SearchUsers)
		adminApi.GET("/users/:id", mainHandler.GetUser)
		adminApi.PATCH("/users/:id", mainHandler.PatchUser)
		adminApi.DELETE("/users/:id", mainHandler.DeleteUser)

		// Statistics
		// adminApi.GET("/statistics", mainHandler.GetStatistics)
	}

	// Swagger Route
	docs.SwaggerInfo.BasePath = middlewares.APIPrefix
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
