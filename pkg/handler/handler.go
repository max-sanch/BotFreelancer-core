package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/max-sanch/BotFreelancer-core/pkg/service"
	"github.com/spf13/viper"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	if viper.GetString("releaseMode") == "True" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	api := router.Group("/api")
	{
		channels := api.Group("/channels")
		{
			channels.GET("/data")
			channels.POST("/create")
			channels.POST("/update")
			channels.POST("/delete")
		}

		users := api.Group("/users")
		{
			users.GET("/data")
			users.POST("/create")
			users.POST("/update")
		}
	}

	return router
}
