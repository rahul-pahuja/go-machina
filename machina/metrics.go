package machina

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all the Prometheus metrics for the FSM
type Metrics struct {
	TransitionsTotal     *prometheus.CounterVec
	TransitionErrors     *prometheus.CounterVec
	TransitionDuration   *prometheus.HistogramVec
	AutoTransitionsTotal *prometheus.CounterVec
}

// NewMetrics creates a new Metrics instance with all the required metrics
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		TransitionsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gomachina_transitions_total",
				Help: "Total number of state transitions",
			},
			[]string{"from_state", "to_state", "event"},
		),
		TransitionErrors: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gomachina_transition_errors_total",
				Help: "Total number of transition errors",
			},
			[]string{"from_state", "event", "error_type"},
		),
		TransitionDuration: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "gomachina_transition_duration_seconds",
				Help:    "Duration of state transitions in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"from_state", "to_state", "event"},
		),
		AutoTransitionsTotal: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "gomachina_auto_transitions_total",
				Help: "Total number of automatic transitions",
			},
			[]string{"from_state", "to_state", "event"},
		),
	}

	return m
}
