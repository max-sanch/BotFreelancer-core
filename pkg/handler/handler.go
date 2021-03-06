package handler

import (
	"github.com/max-sanch/BotFreelancer-core/pkg/service"

	"github.com/gin-gonic/gin"
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
			channels.GET("/data", h.getTasksChannel)
			channels.POST("/channel", h.getChannel)
			channels.POST("/create", h.createChannel)
			channels.POST("/update", h.updateChannel)
			channels.POST("/delete", h.deleteChannel)
		}

		users := api.Group("/users")
		{
			users.GET("/data", h.getTasksUser)
			users.POST("/user", h.getUser)
			users.POST("/create", h.createUser)
			users.POST("/update", h.updateUser)
		}
	}

	return router
}
