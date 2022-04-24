package handler

import (
	"github.com/gin-gonic/gin"
	core "github.com/max-sanch/BotFreelancer-core"
	"net/http"
)

func (h *Handler) getTasksChannel(c *gin.Context) {
	tasks := make([]core.ChannelTaskResponse, 0)
	tasks = append(tasks, core.ChannelTaskResponse{
		ApiId: 123456,
		ApiHash: "0123456789abcdef0123456789abcdef",
		Title: "Test",
		Body: "Body test",
		Url: "github.com/max-sanch/BotFreelancer-core",
	})

	c.JSON(http.StatusOK, core.ChannelTasksResponse{
		Tasks: tasks,
	})
}

func (h *Handler) getChannel(c *gin.Context) {
	var input core.ApiIdInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	channel, err := h.services.Channel.GetByApiId(input.ApiId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, channel)
}

func (h *Handler) createChannel(c *gin.Context) {
	var input core.ChannelInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Channel.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updateChannel(c *gin.Context) {
	var input core.ChannelInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Channel.Update(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) deleteChannel(c *gin.Context) {
	var input core.ApiIdInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Channel.Delete(input.ApiId); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}
