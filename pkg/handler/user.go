package handler

import (
	"github.com/gin-gonic/gin"
	core "github.com/max-sanch/BotFreelancer-core"
	"net/http"
)

func (h *Handler) getDataUser(c *gin.Context) {
	tasks := make([]core.UserTask, 0)
	tasks = append(tasks, core.UserTask{
		TGID: 123456,
		Title: "Test",
		Body: "Body test",
		Url: "github.com/max-sanch/BotFreelancer-core",
	})

	c.JSON(http.StatusOK, core.UserTaskResponse{
		Tasks: tasks,
	})
}

func (h *Handler) getUser(c *gin.Context) {
	var input core.UserAPIIDInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.GetUser(input.TGID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var input core.UserInput

	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateUser(input)
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
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.UpdateUser(input)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
