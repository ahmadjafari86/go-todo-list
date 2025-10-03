package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/ahmadjafari86/go-todo-list/internal/service"
	"github.com/ahmadjafari86/go-todo-list/internal/validation"
)

type registerPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	svc       service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc, validator: validator.New()}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Register user"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} validation.ProblemDetails
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var p registerPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", errs.Error())
			return
		}
		validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}

	u, err := h.svc.Register(p.Email, p.Password)
	if err != nil {
		validation.RespondProblem(c, http.StatusBadRequest, "Registration Failed", err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} validation.ProblemDetails
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var p loginPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", errs.Error())
			return
		}
		validation.RespondProblem(c, http.StatusBadRequest, "Invalid Request", err.Error())
		return
	}
	token, err := h.svc.Login(p.Email, p.Password)
	if err != nil {
		validation.RespondProblem(c, http.StatusUnauthorized, "Invalid Credentials", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
