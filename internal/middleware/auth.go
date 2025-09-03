package middleware

import (
	"github.com/gin-gonic/gin"
)

// DummyAuth is a dummy authentication middleware that hardcodes userID=1
func DummyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}
}

