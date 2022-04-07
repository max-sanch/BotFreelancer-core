package handler

import (
	"github.com/gin-gonic/gin"
	core "github.com/max-sanch/BotFreelancer-core"
	"net/http"
)

func (h *Handler) getDataChannel(c *gin.Context) {
	tasks := make([]core.Task, 0)
	tasks = append(tasks, core.Task{
		ID: 123456,
		APIHash: "0123456789abcdef0123456789abcdef",
		Title: "Test",
		Body: "Body test",
		Url: "github.com/max-sanch/BotFreelancer-core",
	})

	c.JSON(http.StatusOK, core.ChannelRequest{
		Tasks: tasks,
	})
}

func (h *Handler) createChannel(c *gin.Context) {

}

func (h *Handler) updateChannel(c *gin.Context) {

}

func (h *Handler) deleteChannel(c *gin.Context) {

}
