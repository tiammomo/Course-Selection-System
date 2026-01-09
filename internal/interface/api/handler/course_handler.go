package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"course_select/internal/application/dto"
	appService "course_select/internal/application/service"
	"course_select/internal/domain/model"
	domainService "course_select/internal/domain/service"
	"course_select/internal/pkg/errcode"
	"course_select/internal/pkg/response"
)

// CourseHandler 课程处理器
type CourseHandler struct {
	courseService       *domainService.CourseService
	scheduleService     *domainService.ScheduleService
	selectionAppService *appService.SelectionAppService
}

// NewCourseHandler 创建课程处理器
func NewCourseHandler(
	courseService *domainService.CourseService,
	scheduleService *domainService.ScheduleService,
	selectionAppService *appService.SelectionAppService,
) *CourseHandler {
	return &CourseHandler{
		courseService:       courseService,
		scheduleService:     scheduleService,
		selectionAppService: selectionAppService,
	}
}

// CreateCourse 创建课程
// @Summary 创建课程
// @Description 创建新课程
// @Tags course
// @Accept json
// @Produce json
// @Param request body model.CreateCourseRequest true "创建课程请求"
// @Success 200 {object} response.Response
// @Router /course/create [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	var req model.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	course, err := h.courseService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(map[string]string{
		"course_id": strconv.Itoa(course.CourseID),
	}))
}

// GetCourse 获取课程信息
// @Summary 获取课程信息
// @Description 根据课程ID获取课程信息
// @Tags course
// @Produce json
// @Param course_id query string true "课程ID"
// @Success 200 {object} response.Response
// @Router /course/get [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	courseID := c.Query("course_id")
	if courseID == "" {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg("course_id 不能为空")))
		return
	}

	course, err := h.courseService.Get(c.Request.Context(), courseID)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(course.ToResponse()))
}

// BindCourse 绑定课程到教师
// @Summary 绑定课程到教师
// @Description 将课程绑定到指定教师
// @Tags teacher
// @Accept json
// @Produce json
// @Param request body model.BindCourseRequest true "绑定请求"
// @Success 200 {object} response.Response
// @Router /teacher/bind_course [post]
func (h *CourseHandler) BindCourse(c *gin.Context) {
	var req model.BindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := h.courseService.BindCourse(c.Request.Context(), req.CourseID, req.TeacherID); err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(nil))
}

// UnbindCourse 解绑课程
// @Summary 解绑课程
// @Description 将课程从教师解绑
// @Tags teacher
// @Accept json
// @Produce json
// @Param request body model.UnbindCourseRequest true "解绑请求"
// @Success 200 {object} response.Response
// @Router /teacher/unbind_course [post]
func (h *CourseHandler) UnbindCourse(c *gin.Context) {
	var req model.UnbindCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := h.courseService.UnbindCourse(c.Request.Context(), req.CourseID, req.TeacherID); err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(nil))
}

// GetTeacherCourses 获取教师的课程列表
// @Summary 获取教师的课程列表
// @Description 获取当前登录教师的所有课程
// @Tags teacher
// @Produce json
// @Success 200 {object} response.Response
// @Router /teacher/get_course [get]
func (h *CourseHandler) GetTeacherCourses(c *gin.Context) {
	teacherID, ok := GetUserIDFromSession(c)
	if !ok {
		c.JSON(200, response.Fail(errcode.LoginRequired))
		return
	}

	courses, err := h.courseService.GetTeacherCourses(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	var courseList []*model.CourseResponse
	for _, c := range courses {
		courseList = append(courseList, c.ToResponse())
	}

	c.JSON(200, response.Success(map[string]interface{}{
		"course_list": courseList,
	}))
}

// ScheduleCourse 排课
// @Summary 排课
// @Description 使用二分图匹配算法自动排课
// @Tags course
// @Accept json
// @Produce json
// @Param request body model.ScheduleCourseRequest true "排课请求"
// @Success 200 {object} response.Response
// @Router /course/schedule [post]
func (h *CourseHandler) ScheduleCourse(c *gin.Context) {
	var req model.ScheduleCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	result, err := h.scheduleService.Schedule(c.Request.Context(), req.TeacherCourseRelationShip)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(map[string]interface{}{
		"assignments": result,
	}))
}

// BookCourse 学生选课
// @Summary 学生选课
// @Description 学生选择课程
// @Tags student
// @Accept json
// @Produce json
// @Param request body dto.BookCourseRequest true "选课请求"
// @Success 200 {object} response.Response
// @Router /student/book_course [post]
func (h *CourseHandler) BookCourse(c *gin.Context) {
	var req dto.BookCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, response.Fail(errcode.ParamInvalid.WithMsg(err.Error())))
		return
	}

	if err := h.selectionAppService.BookCourse(c.Request.Context(), &req); err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(nil))
}

// GetStudentCourses 获取学生课表
// @Summary 获取学生课表
// @Description 获取当前登录学生的所有课程
// @Tags student
// @Produce json
// @Success 200 {object} response.Response
// @Router /student/course [get]
func (h *CourseHandler) GetStudentCourses(c *gin.Context) {
	studentID, ok := GetUserIDFromSession(c)
	if !ok {
		c.JSON(200, response.Fail(errcode.LoginRequired))
		return
	}

	courses, err := h.selectionAppService.GetStudentCourses(c.Request.Context(), studentID)
	if err != nil {
		c.JSON(200, response.FailWithError(err))
		return
	}

	c.JSON(200, response.Success(dto.GetStudentCourseResponse{
		CourseList: courses,
	}))
}
