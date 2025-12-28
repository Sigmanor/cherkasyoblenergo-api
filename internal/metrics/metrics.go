package metrics

import (
	"expvar"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	requestCount   = expvar.NewInt("requests_total")
	errorCount     = expvar.NewInt("errors_total")
	totalLatencyNs int64
	requestsForAvg int64
)

func init() {
	expvar.Publish("avg_latency_ms", expvar.Func(func() any {
		reqs := atomic.LoadInt64(&requestsForAvg)
		if reqs == 0 {
			return 0
		}
		totalNs := atomic.LoadInt64(&totalLatencyNs)
		return float64(totalNs) / float64(reqs) / 1e6
	}))
}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()

		latency := time.Since(start).Nanoseconds()
		atomic.AddInt64(&totalLatencyNs, latency)
		atomic.AddInt64(&requestsForAvg, 1)

		requestCount.Add(1)
		if err != nil || c.Response().StatusCode() >= 500 {
			errorCount.Add(1)
		}

		return err
	}
}
