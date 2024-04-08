package retcode

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestRetcodeCounter *prometheus.CounterVec

func init() {
	RequestRetcodeCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "request_retcode_counter",
		Help: "each request retcode counter",
	}, []string{"path", "method", "code", "desc"})
}
