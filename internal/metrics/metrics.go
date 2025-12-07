package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of HTTP requests processed.",
	}, []string{"path", "method"})
)

func InitMetrics() {
	prometheus.MustRegister(RequestCount)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
