package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request details
		if raw != "" {
			path = path + "?" + raw
		}

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		logEvent := log.Info()

		if statusCode >= 400 {
			logEvent = log.Error().
				Int("status", statusCode).
				Str("error", c.Errors.String())
		}

		logEvent.
			Str("method", method).
			Str("path", path).
			Str("ip", clientIP).
			Int("status", statusCode).
			Dur("latency", latency).
			Msg("Request")
	}
}
