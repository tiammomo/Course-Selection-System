package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP 请求总数
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// HTTP 请求延迟
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 活跃连接数
	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	// 选课成功数
	bookingSuccessTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "booking_success_total",
			Help: "Total number of successful course bookings",
		},
	)

	// 选课失败数
	bookingFailTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "booking_fail_total",
			Help: "Total number of failed course bookings",
		},
		[]string{"reason"},
	)

	// 课程容量
	courseCapacity = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "course_capacity_remaining",
			Help: "Remaining capacity for each course",
		},
		[]string{"course_id"},
	)
)

// RecordRequest 记录 HTTP 请求
func RecordRequest(method, path, status string, duration float64) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration)
}

// IncActiveConnections 增加活跃连接数
func IncActiveConnections() {
	activeConnections.Inc()
}

// DecActiveConnections 减少活跃连接数
func DecActiveConnections() {
	activeConnections.Dec()
}

// IncBookingSuccess 增加选课成功数
func IncBookingSuccess() {
	bookingSuccessTotal.Inc()
}

// IncBookingFail 增加选课失败数
func IncBookingFail(reason string) {
	bookingFailTotal.WithLabelValues(reason).Inc()
}

// SetCourseCapacity 设置课程容量
func SetCourseCapacity(courseID string, capacity int) {
	courseCapacity.WithLabelValues(courseID).Set(float64(capacity))
}
