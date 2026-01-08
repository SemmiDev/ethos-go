package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/semmidev/ethos-go/internal/common/decorator"
)

// PrometheusMetricsClient implements decorator.MetricsClient using Prometheus
type PrometheusMetricsClient struct {
	mu       sync.RWMutex
	counters map[string]prometheus.Counter
}

// NewPrometheusMetricsClient creates a Prometheus-backed metrics client
func NewPrometheusMetricsClient() *PrometheusMetricsClient {
	return &PrometheusMetricsClient{
		counters: make(map[string]prometheus.Counter),
	}
}

// Ensure PrometheusMetricsClient implements decorator.MetricsClient
var _ decorator.MetricsClient = (*PrometheusMetricsClient)(nil)

// Inc increments a counter by the specified value
func (c *PrometheusMetricsClient) Inc(key string, value int) {
	c.mu.RLock()
	counter, exists := c.counters[key]
	c.mu.RUnlock()

	if !exists {
		c.mu.Lock()
		// Double-check after acquiring write lock
		if counter, exists = c.counters[key]; !exists {
			counter = promauto.NewCounter(prometheus.CounterOpts{
				Name: sanitizeMetricName(key),
				Help: "Auto-generated counter for " + key,
			})
			c.counters[key] = counter
		}
		c.mu.Unlock()
	}

	counter.Add(float64(value))
}

// sanitizeMetricName converts arbitrary strings to valid Prometheus metric names
func sanitizeMetricName(name string) string {
	result := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9' && i > 0) || c == '_' {
			result = append(result, c)
		} else if c == '.' || c == '-' || c == ' ' {
			result = append(result, '_')
		}
	}
	return string(result)
}
