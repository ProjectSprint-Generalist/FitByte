package middleware

import (
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)


func Logger() gin.HandlerFunc {
	return logger.SetLogger(
		logger.WithLogger(func(c *gin.Context, l zerolog.Logger) zerolog.Logger {
			return l.Output(gin.DefaultWriter).With().
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Dur("latency", time.Since(time.Now())).
				Logger()
		}),
	)
}
