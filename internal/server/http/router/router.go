package router

import (
	"avito/internal/server/http/middleware"
	"os"

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

func (r *HttpRouter) Register(bannerHandler BannerHandler) {
	userBannerRouter := r.router.Group("/user_banner")
	// ?tag_id=<integer>&feature_id=<integer>&use_last_revision=false
	userBannerRouter.GET("", bannerHandler.GetByTagAndFeature)

	bannerRouter := r.router.Group("/banner")
	// ?feature_id=<integer>&tag_id=<integer>&limit=<integer>&offset=<integer>
	bannerRouter.GET("", bannerHandler.GetManyByTagOrFeature)
	bannerRouter.POST("", bannerHandler.Create)
	bannerRouter.PATCH("/:id", bannerHandler.Update)
	bannerRouter.DELETE("/:id", bannerHandler.Delete)
}

func (r *HttpRouter) Run() error {
	if err := r.router.Run(os.Getenv("ADDRESS")); err != nil {
		return err
	}
	return nil
}
