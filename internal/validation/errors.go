package validation

import (
	"github.com/gin-gonic/gin"
)

// ProblemDetails defines RFC7807 Problem+JSON structure
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func NewProblem(status int, title string, detail string, instance string) ProblemDetails {
	return ProblemDetails{
		Type:     "about:blank",
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

func RespondProblem(c *gin.Context, status int, title, detail string) {
	problem := NewProblem(status, title, detail, c.Request.RequestURI)
	c.Header("Content-Type", "application/problem+json")
	c.JSON(status, problem)
	c.Abort()
}
