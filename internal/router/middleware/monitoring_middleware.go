package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/metrics"
	"github.com/gin-gonic/gin"
)

func Monitoring() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		metrics.ActiveConnections.Inc()
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()
		metrics.ActiveConnections.Dec()
		latencyTime := endTime.Sub(startTime).Seconds()
		reqMethod := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		metrics.LatencyHistogram.WithLabelValues(reqMethod, ctx.Request.RequestURI, http.StatusText(statusCode)).Observe(latencyTime)
		metrics.RequestCount.WithLabelValues(reqMethod, ctx.FullPath(), strconv.Itoa(ctx.Writer.Status())).Inc()
	}
}
