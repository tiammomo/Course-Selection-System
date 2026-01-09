package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"course_select/internal/config"
	"course_select/internal/infrastructure/database"
	"course_select/internal/infrastructure/mq"
	redisClient "course_select/internal/infrastructure/redis"
	"course_select/internal/interface/api/handler"
	"course_select/internal/interface/api/middleware"
	"course_select/internal/interface/api/router"
	"course_select/internal/pkg/logger"
	appService "course_select/internal/application/service"
	domainService "course_select/internal/domain/service"
)

func main() {
	// 1. 初始化配置
	if err := config.Init("config.yaml"); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}
	cfg := config.Get()

	// 2. 初始化日志
	if err := logger.Init(&logger.Config{
		Level:  cfg.Logging.Level,
		Format: cfg.Logging.Format,
		Output: cfg.Logging.Output,
		Path:   cfg.Logging.Path,
	}); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	// 3. 初始化数据库
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", logger.Err(err))
	}
	defer database.Close()

	// 4. 初始化 Redis
	redisCli, err := redisClient.New(&cfg.Redis)
	if err != nil {
		logger.Fatal("Failed to init redis", logger.Err(err))
	}
	defer redisCli.Close()

	// 5. 初始化 RocketMQ (可选，如果没有配置则跳过)
	var mqCli *mq.Client
	if cfg.RocketMQ.NameServer != "" {
		mqCli, err = mq.New(&cfg.RocketMQ)
		if err != nil {
			logger.Error("Failed to init rocketmq", logger.Err(err))
			// 不致命，继续运行
		} else {
			defer mqCli.Close()
		}
	}

	// 6. 初始化仓储
	memberRepo := database.NewMemberRepo(database.Get())
	courseRepo := database.NewCourseRepo(database.Get())
	bindRepo := database.NewBindRepo(database.Get())
	choiceRepo := database.NewChoiceRepo(database.Get())

	// 7. 初始化服务
	authService := domainService.NewAuthService(memberRepo, cfg.Auth.SessionKey, cfg.Auth.CookieName, cfg.Auth.SessionExpireHours)
	memberService := domainService.NewMemberService(memberRepo)
	courseService := domainService.NewCourseService(courseRepo, bindRepo, choiceRepo)
	scheduleService := domainService.NewScheduleService(courseRepo, bindRepo)

	// 8. 初始化应用服务
	selectionAppService := appService.NewSelectionAppService(
		courseRepo,
		choiceRepo,
		bindRepo,
		redisCli,
		mqCli,
		nil, // 限流器在中间件中处理
	)

	// 9. 初始化中间件
	authMiddleware := middleware.NewAuthMiddleware(authService, cfg.Auth.SessionKey)
	limiterMiddleware := middleware.NewLimiterMiddleware(cfg.RateLimit.QPS, cfg.RateLimit.Burst)
	loggerMiddleware := middleware.NewLoggerMiddleware()

	// 10. 初始化 Handler
	authHandler := handler.NewAuthHandler(authService, cfg.Auth.SessionKey, cfg.Auth.CookieName)
	memberHandler := handler.NewMemberHandler(memberService)
	courseHandler := handler.NewCourseHandler(courseService, scheduleService, selectionAppService)

	// 11. 初始化路由
	router := router.NewRouter(authHandler, memberHandler, courseHandler, authMiddleware, limiterMiddleware)

	// 12. 初始化 Gin
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(loggerMiddleware.ZapLogger())
	engine.Use(loggerMiddleware.Recovery())

	// Session 中间件
	store, err := redis.NewStoreWithPool(
		redisCli.Pool(),
		[]byte(cfg.Auth.SessionKey),
	)
	if err != nil {
		logger.Fatal("Failed to init session store", logger.Err(err))
	}
	engine.Use(sessions.Sessions("camp-session", store))

	// 注册路由
	router.RegisterRoutes(engine)

	// 13. 启动 Prometheus 指标端点
	if cfg.Metrics.Enabled {
		engine.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// 14. 启动服务
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	// 优雅关闭
	go func() {
		logger.Info("Starting server", logger.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Err(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", logger.Err(err))
	}

	logger.Info("Server exiting")
}
