package handler

import (
	"github.com/gin-gonic/gin"
	core "github.com/max-sanch/BotFreelancer-core"
	"net/http"
)

func (h *Handler) getDataChannel(c *gin.Context) {
	tasks := make([]core.ChannelTask, 0)
	tasks = append(tasks, core.ChannelTask{
		APIID: 123456,
		APIHash: "0123456789abcdef0123456789abcdef",
		Title: "Test",
		Body: "Body test",
		Url: "github.com/max-sanch/BotFreelancer-core",
	})

	c.JSON(http.StatusOK, core.ChannelTaskResponse{
		Tasks: tasks,
	})
}

func (h *Handler) getChannel(c *gin.Context) {
	var input core.ChannelAPIIDInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	channel, err := h.services.GetChannel(input.APIID)
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

	id, err := h.services.CreateChannel(input)
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

	id, err := h.services.UpdateChannel(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) deleteChannel(c *gin.Context) {
	var input core.ChannelAPIIDInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DeleteChannel(input.APIID); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"status": "OK",
	})
}
