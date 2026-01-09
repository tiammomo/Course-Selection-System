package load

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Config 压测配置
type Config struct {
	BaseURL      string        // 服务基础地址
	Endpoint     string        // 测试端点
	Method       string        // HTTP 方法
	RequestBody  interface{}   // 请求体
	Workers      int           // 并发 worker 数量
	TotalRequests int          // 总请求数
	Timeout      time.Duration // 请求超时
}

// Result 压测结果
type Result struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	TotalDuration   time.Duration
	RequestsPerSec  float64
	MinLatency      time.Duration
	MaxLatency      time.Duration
	AvgLatency      time.Duration
	Percentiles     map[float64]time.Duration
}

// RunLoadTest 运行压测
func RunLoadTest(t *testing.T, cfg Config) *Result {
	t.Helper()

	var (
		totalRequests   int64
		successRequests int64
		failedRequests  int64
		minLatency      time.Duration = time.Hour
		maxLatency      time.Duration
		latencies       []time.Duration
	)

	// 序列化请求体
	var body []byte
	var err error
	if cfg.RequestBody != nil {
		body, err = json.Marshal(cfg.RequestBody)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	// 创建工作池
	var wg sync.WaitGroup
	workerCh := make(chan struct{}, cfg.Workers)  // 限制并发数

	start := time.Now()

	for i := 0; i < cfg.TotalRequests; i++ {
		workerCh <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			<-workerCh

			reqStart := time.Now()
			client := &http.Client{Timeout: cfg.Timeout}

			var req *http.Request
			if len(body) > 0 {
				req, err = http.NewRequest(cfg.Method, cfg.BaseURL+cfg.Endpoint, bytes.NewReader(body))
			} else {
				req, err = http.NewRequest(cfg.Method, cfg.BaseURL+cfg.Endpoint, nil)
			}

			if err != nil {
				atomic.AddInt64(&failedRequests, 1)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				atomic.AddInt64(&failedRequests, 1)
				return
			}
			defer resp.Body.Close()

			latency := time.Since(reqStart)
			atomic.AddInt64(&totalRequests, 1)

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				atomic.AddInt64(&successRequests, 1)
			} else {
				atomic.AddInt64(&failedRequests, 1)
			}

			// 记录延迟
			mu.Lock()
			latencies = append(latencies, latency)
			if latency < minLatency {
				minLatency = latency
			}
			if latency > maxLatency {
				maxLatency = latency
			}
			mu.Unlock()
		}()
	}

	wg.Wait()
	totalDuration := time.Since(start)

	// 计算统计数据
	avgLatency := sumLatencies(latencies) / time.Duration(len(latencies))

	result := &Result{
		TotalRequests:   totalRequests,
		SuccessRequests: successRequests,
		FailedRequests:  failedRequests,
		TotalDuration:   totalDuration,
		RequestsPerSec:  float64(totalRequests) / totalDuration.Seconds(),
		MinLatency:      minLatency,
		MaxLatency:      maxLatency,
		AvgLatency:      avgLatency,
		Percentiles:     calculatePercentiles(latencies, []float64{0.5, 0.9, 0.95, 0.99}),
	}

	return result
}

var mu sync.Mutex

func sumLatencies(latencies []time.Duration) time.Duration {
	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	return sum
}

func calculatePercentiles(latencies []time.Duration, percentiles []float64) map[float64]time.Duration {
	n := len(latencies)
	if n == 0 {
		result := make(map[float64]time.Duration)
		for _, p := range percentiles {
			result[p] = 0
		}
		return result
	}

	// 排序
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if latencies[i] > latencies[j] {
				latencies[i], latencies[j] = latencies[j], latencies[i]
			}
		}
	}

	result := make(map[float64]time.Duration)
	for _, p := range percentiles {
		index := int(float64(n) * p)
		if index >= n {
			index = n - 1
		}
		result[p] = latencies[index]
	}

	return result
}

// PrintResult 打印压测结果
func PrintResult(t *testing.T, result *Result) {
	t.Helper()

	fmt.Println("\n========== Load Test Results ==========")
	fmt.Printf("Total Requests:      %d\n", result.TotalRequests)
	fmt.Printf("Success Requests:    %d\n", result.SuccessRequests)
	fmt.Printf("Failed Requests:     %d\n", result.FailedRequests)
	fmt.Printf("Success Rate:        %.2f%%\n", float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("Total Duration:      %v\n", result.TotalDuration)
	fmt.Printf("Requests/sec (RPS):  %.2f\n", result.RequestsPerSec)
	fmt.Println("\nLatency:")
	fmt.Printf("  Min:    %v\n", result.MinLatency)
	fmt.Printf("  Avg:    %v\n", result.AvgLatency)
	fmt.Printf("  Max:    %v\n", result.MaxLatency)
	fmt.Println("\nPercentiles:")
	for p, latency := range result.Percentiles {
		fmt.Printf("  P%.0f:   %v\n", p*100, latency)
	}
	fmt.Println("========================================")
}

// ========== 选课接口压测 ==========

// BookCourseTest 选课接口压测
func TestBookCourseLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/student/book_course",
		Method:       "POST",
		RequestBody:  map[string]string{"student_id": "1", "course_id": "1"},
		Workers:      1000,
		TotalRequests: 20000,
		Timeout:      5 * time.Second,
	}

	// 预热
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("POST", cfg.BaseURL+"/api/v1/auth/login", bytes.NewReader([]byte(`{"username":"test","password":"test"}`)))
	req.Header.Set("Content-Type", "application/json")
	client.Do(req)

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)

	// 验证结果
	if result.SuccessRequests < int64(cfg.TotalRequests)*99/100 {
		t.Errorf("Success rate too low: %.2f%%", float64(result.SuccessRequests)/float64(cfg.TotalRequests)*100)
	}
}

// GetStudentCourseTest 获取课表接口压测
func TestGetStudentCourseLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/student/course?student_id=1",
		Method:       "GET",
		Workers:      1000,
		TotalRequests: 20000,
		Timeout:      5 * time.Second,
	}

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)

	if result.SuccessRequests < int64(cfg.TotalRequests)*99/100 {
		t.Errorf("Success rate too low: %.2f%%", float64(result.SuccessRequests)/float64(cfg.TotalRequests)*100)
	}
}

// ========== 认证接口压测 ==========

// LoginTest 登录接口压测
func TestLoginLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/auth/login",
		Method:       "POST",
		RequestBody:  map[string]string{"username": "testuser", "password": "testpass"},
		Workers:      500,
		TotalRequests: 5000,
		Timeout:      5 * time.Second,
	}

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)
}

// ========== 成员管理接口压测 ==========

// GetMemberListTest 成员列表压测
func TestGetMemberListLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/member/list?offset=0&limit=20",
		Method:       "GET",
		Workers:      500,
		TotalRequests: 5000,
		Timeout:      5 * time.Second,
	}

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)
}

// ========== 课程管理接口压测 ==========

// CreateCourseTest 创建课程压测
func TestCreateCourseLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/course/create",
		Method:       "POST",
		RequestBody:  map[string]interface{}{"name": "TestCourse", "cap": 100},
		Workers:      100,
		TotalRequests: 1000,
		Timeout:      5 * time.Second,
	}

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)
}

// GetCourseTest 获取课程压测
func TestGetCourseLoad(t *testing.T) {
	cfg := Config{
		BaseURL:      "http://localhost:8080",
		Endpoint:     "/api/v1/course/get?course_id=1",
		Method:       "GET",
		Workers:      500,
		TotalRequests: 5000,
		Timeout:      5 * time.Second,
	}

	result := RunLoadTest(t, cfg)
	PrintResult(t, result)
}

// ========== 基准测试 ==========

func BenchmarkBookCourse(b *testing.B) {
	client := &http.Client{Timeout: 5 * time.Second}
	url := "http://localhost:8080/api/v1/student/book_course"
	body := []byte(`{"student_id":"1","course_id":"1"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		client.Do(req)
	}
}

// ConcurrentBookCourseTest 并发选课测试（详细统计）
func TestConcurrentBookCourse(t *testing.T) {
	const (
		numWorkers     = 1000
		requestsPerWorker = 20
	)

	client := &http.Client{Timeout: 5 * time.Second}
	url := "http://localhost:8080/api/v1/student/book_course"
	body := []byte(`{"student_id":"1","course_id":"1"}`)

	var (
		success int64
		fail    int64
		mu      sync.Mutex
		latencies []time.Duration
	)

	var wg sync.WaitGroup
	sem := make(chan struct{}, numWorkers)

	start := time.Now()

	for i := 0; i < numWorkers; i++ {
		for j := 0; j < requestsPerWorker; j++ {
			sem <- struct{}{}
			wg.Add(1)

			go func() {
				defer wg.Done()
				<-sem

				reqStart := time.Now()
				req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					atomic.AddInt64(&fail, 1)
					return
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()

				latency := time.Since(reqStart)
				mu.Lock()
				latencies = append(latencies, latency)
				mu.Unlock()

				if resp.StatusCode == 0 {
					atomic.AddInt64(&fail, 1)
				} else {
					atomic.AddInt64(&success, 1)
				}
			}()
		}
	}

	wg.Wait()
	duration := time.Since(start)

	total := success + fail
	t.Logf("\n========== Concurrent Test Results ==========")
	t.Logf("Workers:            %d", numWorkers)
	t.Logf("Requests/Worker:    %d", requestsPerWorker)
	t.Logf("Total Requests:     %d", total)
	t.Logf("Success:            %d", success)
	t.Logf("Failed:             %d", fail)
	t.Logf("Duration:           %v", duration)
	t.Logf("RPS:                %.2f", float64(total)/duration.Seconds())
	t.Logf("Success Rate:       %.2f%%", float64(success)/float64(total)*100)
	t.Logf("============================================")
}
