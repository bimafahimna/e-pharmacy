package middleware

import (
	"fmt"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.Log
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		if lastErr := c.Errors.Last(); lastErr != nil {
			log.WithFields(map[string]any{
				"METHOD":    reqMethod,
				"URI":       reqUri,
				"STATUS":    statusCode,
				"LATENCY":   latencyTime,
				"CLIENT_IP": clientIP,
			}).Error(c.Errors)
			fmt.Printf("\nERRORS: %+v", c.Errors)
			return
		}

		if reqUri == "/metrics" {
			log.WithFields(map[string]any{
				"METHOD":    reqMethod,
				"URI":       reqUri,
				"STATUS":    statusCode,
				"LATENCY":   latencyTime,
				"CLIENT_IP": clientIP,
			}).Debugf("REQUEST %s %s SUCCESS", reqMethod, reqUri)
			return
		}

		log.WithFields(map[string]any{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
		}).Infof("REQUEST %s %s SUCCESS", reqMethod, reqUri)
	}
}
