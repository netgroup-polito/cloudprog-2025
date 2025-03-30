package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "endpoint"},
)

func RegisterRequestsCounter(r *gin.Engine) {
	if err := prometheus.Register(requestsCounter); err != nil {
		panic(err)
	}

	r.Use(func(c *gin.Context) {
		requestsCounter.WithLabelValues(c.Request.Method, c.FullPath()).Inc()
		c.Next()
	})

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
