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

// CreateTodo godoc
// @Summary Create a todo
// @Description Create a new todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.CreateTodoRequest true "Todo"
// @Success 201 {object} models.Todo
// @Failure 400 {object} validation.ProblemDetails
// @Failure 401 {object} validation.ProblemDetails
// @Router /api/todos [post]
// @Security BearerAuth
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

// ListTodos godoc
// @Summary List todos
// @Description Get all todos for the authenticated user
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 401 {object} validation.ProblemDetails
// @Router /api/todos [get]
// @Security BearerAuth
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

// GetTodo godoc
// @Summary Get a todo
// @Description Get a todo by ID (must belong to the authenticated user)
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 401 {object} validation.ProblemDetails
// @Failure 404 {object} validation.ProblemDetails
// @Router /api/todos/{id} [get]
// @Security BearerAuth
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

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update an existing todo (must belong to the authenticated user)
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body models.UpdateTodoRequest true "Updated Todo"
// @Success 200 {object} models.Todo
// @Failure 400 {object} validation.ProblemDetails
// @Failure 401 {object} validation.ProblemDetails
// @Failure 404 {object} validation.ProblemDetails
// @Router /api/todos/{id} [put]
// @Security BearerAuth
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

// ToggleComplete godoc
// @Summary Toggle todo completion
// @Description Mark a todo as complete/incomplete
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 401 {object} validation.ProblemDetails
// @Failure 404 {object} validation.ProblemDetails
// @Router /api/todos/{id}/complete [patch]
// @Security BearerAuth
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

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo (must belong to the authenticated user)
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204
// @Failure 401 {object} validation.ProblemDetails
// @Failure 404 {object} validation.ProblemDetails
// @Router /api/todos/{id} [delete]
// @Security BearerAuth
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
