package middleware

import "github.com/prometheus/client_golang/prometheus"

var PrometheusCli *prometheus.CounterVec

func PrometheusInit() {
	// 初始化Prometheus指标
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "endpoint"})
	prometheus.MustRegister(counter)

	PrometheusCli = counter
}
