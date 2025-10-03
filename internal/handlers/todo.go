package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmadjafari86/go-todo-list/internal/models"
	"github.com/ahmadjafari86/go-todo-list/internal/service"
	"github.com/ahmadjafari86/go-todo-list/internal/validation"
	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	svc service.TodoService
}

func NewTodoHandler(svc service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var payload models.Todo
	if err := c.ShouldBindJSON(&payload); err != nil {
		validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	ownerID := getUserIDFromContext(c)
	if ownerID == 0 {
		validation.RespondProblem(c, http.StatusUnauthorized, "Unauthorized", "missing user id")
		return
	}
	if err := h.svc.CreateTodo(&payload, ownerID); err != nil {
		validation.RespondProblem(c, http.StatusBadRequest, "Create Failed", err.Error())
		return
	}
	c.JSON(http.StatusCreated, payload)
}

func (h *TodoHandler) ListTodos(c *gin.Context) {
	ownerID := getUserIDFromContext(c)
	if ownerID == 0 {
		validation.RespondProblem(c, http.StatusUnauthorized, "Unauthorized", "missing user id")
		return
	}
	todos, err := h.svc.ListTodos(ownerID)
	if err != nil {
		validation.RespondProblem(c, http.StatusInternalServerError, "Server Error", err.Error())
		return
	}
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ownerID := getUserIDFromContext(c)
	todo, err := h.svc.GetTodo(uint(id), ownerID)
	if err != nil || todo == nil {
		validation.RespondProblem(c, http.StatusNotFound, "Not Found", "todo not found")
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload models.Todo
	if err := c.ShouldBindJSON(&payload); err != nil {
		validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	payload.ID = uint(id)
	ownerID := getUserIDFromContext(c)
	if err := h.svc.UpdateTodo(&payload, ownerID); err != nil {
		validation.RespondProblem(c, http.StatusBadRequest, "Update Failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, payload)
}

func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ownerID := getUserIDFromContext(c)
	todo, err := h.svc.ToggleComplete(uint(id), ownerID)
	if err != nil || todo == nil {
		validation.RespondProblem(c, http.StatusNotFound, "Not Found", "todo not found")
		return
	}
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	ownerID := getUserIDFromContext(c)
	if err := h.svc.DeleteTodo(uint(id), ownerID); err != nil {
		validation.RespondProblem(c, http.StatusNotFound, "Not Found", "todo not found")
		return
	}
	c.Status(http.StatusNoContent)
}

func getUserIDFromContext(c *gin.Context) uint {
	if v, ok := c.Get("user_id"); ok {
		if idStr, ok2 := v.(string); ok2 {
			if id, err := strconv.Atoi(idStr); err == nil {
				return uint(id)
			}
		}
	}
	return 0
}
