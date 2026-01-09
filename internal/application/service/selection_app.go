package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"course_select/internal/application/dto"
	"course_select/internal/domain/repository"
	mq "course_select/internal/infrastructure/mq"
	"course_select/internal/infrastructure/redis"
	"course_select/internal/pkg/errcode"

	"golang.org/x/time/rate"
)

// SelectionAppService 选课应用服务 (高并发场景)
type SelectionAppService struct {
	courseRepo repository.ICourseRepo
	choiceRepo repository.IChoiceRepo
	bindRepo   repository.IBindRepo
	redis      *redis.Client
	mq         *mq.Client
	limiter    *rate.Limiter
}

// NewSelectionAppService 创建选课应用服务
func NewSelectionAppService(
	courseRepo repository.ICourseRepo,
	choiceRepo repository.IChoiceRepo,
	bindRepo repository.IBindRepo,
	redis *redis.Client,
	mq *mq.Client,
	limiter *rate.Limiter,
) *SelectionAppService {
	return &SelectionAppService{
		courseRepo: courseRepo,
		choiceRepo: choiceRepo,
		bindRepo:   bindRepo,
		redis:      redis,
		mq:         mq,
		limiter:    limiter,
	}
}

// BookCourse 选课 (高并发优化)
func (s *SelectionAppService) BookCourse(ctx context.Context, req *dto.BookCourseRequest) error {
	// 1. 限流检查
	if err := s.limiter.Wait(ctx); err != nil {
		return errcode.UnknownError.WithMsg("请求过于频繁")
	}

	// 2. 解析 ID
	studentID, err := strconv.Atoi(req.StudentID)
	if err != nil {
		return errcode.ParamInvalid
	}
	courseID, err := strconv.Atoi(req.CourseID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// 3. 检查学生是否已选过该课程 (从 Redis)
	enrolled, err := s.redis.SIsMember(ctx, fmt.Sprintf("student:%d:courses", studentID), req.CourseID)
	if err != nil {
		return err
	}
	if enrolled {
		return errcode.RepeatRequest
	}

	// 4. 检查课程是否存在
	course, err := s.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errcode.CourseNotExisted
	}

	// 5. 检查课程容量 (Redis 原子操作)
	remaining, err := s.redis.HIncrBy(ctx, "course:capacity", req.CourseID, -1)
	if err != nil {
		return err
	}
	if remaining < 0 {
		// 回滚
		if _, rollbackErr := s.redis.HIncrBy(ctx, "course:capacity", req.CourseID, 1); rollbackErr != nil {
			return errcode.UnknownError.WithMsg("容量回滚失败")
		}
		return errcode.CourseNotAvailable
	}

	// 6. 发送异步消息到 MQ
	msg := &mq.BookingMessage{
		StudentID: req.StudentID,
		CourseID:  req.CourseID,
		Timestamp: time.Now(),
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return errcode.UnknownError.WithMsg("消息序列化失败")
	}

	// 直接写入 Redis 队列
	if _, err := s.redis.LPush(ctx, "booking:queue", string(body)); err != nil {
		// 回滚
		if _, rollbackErr := s.redis.HIncrBy(ctx, "course:capacity", req.CourseID, 1); rollbackErr != nil {
			return errcode.UnknownError.WithMsg("队列写入失败，回滚也失败")
		}
		return err
	}

	return nil
}

// GetStudentCourses 获取学生课表
func (s *SelectionAppService) GetStudentCourses(ctx context.Context, studentID string) ([]dto.CourseDTO, error) {
	id, err := strconv.Atoi(studentID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}

	// 从 Redis 获取学生选课列表
	courseIDs, err := s.redis.SMembers(ctx, fmt.Sprintf("student:%d:courses", id))
	if err != nil {
		return nil, err
	}

	var courses []dto.CourseDTO
	for _, courseID := range courseIDs {
		cID, err := strconv.Atoi(courseID)
		if err != nil {
			continue // 无效的 courseID，跳过
		}
		course, err := s.courseRepo.GetByID(ctx, cID)
		if err != nil || course == nil {
			continue
		}

		courses = append(courses, dto.CourseDTO{
			CourseID:  courseID,
			Name:      course.Name,
			TeacherID: intToStringPtr(course.TeacherID),
		})
	}

	return courses, nil
}

// intToStringPtr int 转 string 指针
func intToStringPtr(i *int) string {
	if i == nil {
		return ""
	}
	return strconv.Itoa(*i)
}
