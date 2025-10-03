package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        if raw != "" {
            path = path + "?" + raw
        }

        entry := log.WithFields(log.Fields{
            "status":   c.Writer.Status(),
            "method":   c.Request.Method,
            "path":     path,
            "duration": latency.String(),
            "client":   c.ClientIP(),
        })
        if len(c.Errors) > 0 {
            entry.Error(c.Errors.String())
        } else {
            entry.Info("request completed")
        }
    }
}
