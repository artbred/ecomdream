package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

var (
	EndpointErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_endpoint_not_2xx_responses",
		Help: "Total number of not 2xx responses from api",
	}, []string{"error_code", "endpoint"})
)

func Start() {
	http.Handle("/metrics", promhttp.Handler())
	logrus.Error(http.ListenAndServe(":7100", nil))
}

