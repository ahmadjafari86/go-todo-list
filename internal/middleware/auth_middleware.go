package middleware

import (
	"net/http"
	"strings"

	"github.com/ahmadjafari86/go-todo-list/internal/service"
	"github.com/ahmadjafari86/go-todo-list/internal/validation"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authSvc service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			validation.RespondProblem(c, http.StatusUnauthorized, "Unauthorized", "missing authorization header")
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			validation.RespondProblem(c, http.StatusUnauthorized, "Unauthorized", "invalid authorization header format")
			return
		}
		token := parts[1]
		claims, err := authSvc.ParseToken(token)
		if err != nil {
			validation.RespondProblem(c, http.StatusUnauthorized, "Unauthorized", "invalid token")
			return
		}
		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
