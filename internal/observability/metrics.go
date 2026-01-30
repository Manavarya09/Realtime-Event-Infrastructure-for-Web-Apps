package observability
package observability

import (































}	return promhttp.Handler()func MetricsHandler() http.Handler {}	prometheus.MustRegister(EventsProcessed, RequestDuration)func init() {)	)		[]string{"method", "path", "status"},		},			Buckets: prometheus.DefBuckets,			Help:    "HTTP request duration in seconds",			Name:    "http_request_duration_seconds",		prometheus.HistogramOpts{	RequestDuration = prometheus.NewHistogramVec(	)		[]string{"event_name", "status"},		},			Help: "Total number of events processed",			Name: "events_processed_total",		prometheus.CounterOpts{	EventsProcessed = prometheus.NewCounterVec(var ()	"github.com/prometheus/client_golang/prometheus/promhttp"	"github.com/prometheus/client_golang/prometheus"	"net/http"