package responsetime

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Handler will set `X-Response-Time` header in response.
func Handler(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	log.WithFields(log.Fields{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"latency":    time.Since(startTime).Microseconds(),
		"host":       c.ClientIP(),
		"status":     c.Writer.Status(),
		"user-agent": c.Request.UserAgent(),
	}).Info("")
}
