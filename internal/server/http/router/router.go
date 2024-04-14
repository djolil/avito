package router

import (
	"avito/internal/server/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type BannerHandler interface {
	GetByTagAndFeature(ctx *gin.Context)
	GetManyByTagOrFeature(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type UserHandler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type HttpRouter struct {
	router *gin.Engine
}

func NewHttpRouter() *HttpRouter {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(middleware.ErrorHandler)

	return &HttpRouter{router: router}
}

func (r *HttpRouter) Register(bannerHandler BannerHandler, userHandler UserHandler, authMiddleware gin.HandlerFunc) {
	userBannerRouter := r.router.Group("/user_banner")
	userBannerRouter.Use(authMiddleware)
	userBannerRouter.GET("", bannerHandler.GetByTagAndFeature) // ?tag_id=<integer>&feature_id=<integer>&use_last_revision=false

	bannerRouter := r.router.Group("/banner")
	bannerRouter.Use(authMiddleware)
	bannerRouter.Use(middleware.IsAdmin)
	bannerRouter.GET("", bannerHandler.GetManyByTagOrFeature) // ?feature_id=<integer>&tag_id=<integer>&limit=<integer>&offset=<integer>
	bannerRouter.POST("", bannerHandler.Create)
	bannerRouter.PATCH("/:id", bannerHandler.Update)
	bannerRouter.DELETE("/:id", bannerHandler.Delete)

	userRouter := r.router.Group("/user")
	userRouter.POST("/register", userHandler.Register)
	userRouter.POST("/login", userHandler.Login)
}

type ConfigHTTPServer struct {
	Address string `yaml:"address" env-default:"localhost:8080"`
}

func (r *HttpRouter) Run(cfg *ConfigHTTPServer) error {
	if err := r.router.Run(cfg.Address); err != nil {
		return err
	}
	return nil
}
