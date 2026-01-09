package router

import (
	"github.com/gin-gonic/gin"

	"course_select/internal/interface/api/handler"
	"course_select/internal/interface/api/middleware"
)

// Router 路由配置
type Router struct {
	authHandler     *handler.AuthHandler
	memberHandler   *handler.MemberHandler
	courseHandler   *handler.CourseHandler
	authMiddleware  *middleware.AuthMiddleware
	limiterMiddleware *middleware.LimiterMiddleware
}

// NewRouter 创建路由
func NewRouter(
	authHandler *handler.AuthHandler,
	memberHandler *handler.MemberHandler,
	courseHandler *handler.CourseHandler,
	authMiddleware *middleware.AuthMiddleware,
	limiterMiddleware *middleware.LimiterMiddleware,
) *Router {
	return &Router{
		authHandler:      authHandler,
		memberHandler:    memberHandler,
		courseHandler:    courseHandler,
		authMiddleware:   authMiddleware,
		limiterMiddleware: limiterMiddleware,
	}
}

// RegisterRoutes 注册路由
func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// 全局限流
	engine.Use(r.limiterMiddleware.GlobalLimit(4000))

	// 日志和恢复中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// API v1 分组
	v1 := engine.Group("/api/v1")
	{
		// 认证路由
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/logout", r.authHandler.Logout)
			auth.GET("/whoami", r.authMiddleware.RequireAuth(), r.authHandler.WhoAmI)
		}

		// 成员管理路由
		member := v1.Group("/member")
		{
			member.GET("", r.memberHandler.GetMember)
			member.GET("/list", r.memberHandler.GetMemberList)
			// 需要管理员权限
			member.POST("/create", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.CreateMember)
			member.POST("/update", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.UpdateMember)
			member.POST("/delete", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.memberHandler.DeleteMember)
		}

		// 课程管理路由
		course := v1.Group("/course")
		{
			course.GET("/get", r.courseHandler.GetCourse)
			course.POST("/create", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.CreateCourse)
			course.POST("/schedule", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.ScheduleCourse)
		}

		// 教师管理路由
		teacher := v1.Group("/teacher")
		{
			teacher.GET("/get_course", r.courseHandler.GetTeacherCourses)
			teacher.POST("/bind_course", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.BindCourse)
			teacher.POST("/unbind_course", r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin(), r.courseHandler.UnbindCourse)
		}

		// 学生选课路由
		student := v1.Group("/student")
		{
			student.POST("/book_course", r.courseHandler.BookCourse)
			student.GET("/course", r.courseHandler.GetStudentCourses)
		}
	}

	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
