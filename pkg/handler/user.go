package handler

import (
	"github.com/gin-gonic/gin"
	core "github.com/max-sanch/BotFreelancer-core"
	"net/http"
)

func (h *Handler) getTasksUser(c *gin.Context) {
	tasks, err := h.services.User.GetTasks()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, core.UserTasksResponse{
		Tasks: tasks,
	})
}

func (h *Handler) getUser(c *gin.Context) {
	var input core.TgIdInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	user, err := h.services.User.GetByTgId(input.TgId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var input core.UserInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.User.Create(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) updateUser(c *gin.Context) {
	var input core.UserInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.User.Update(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
