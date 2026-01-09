package service

import (
	"context"
	"strconv"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"
	"course_select/internal/pkg/errcode"
)

// CourseService 课程服务
type CourseService struct {
	courseRepo repository.ICourseRepo
	bindRepo   repository.IBindRepo
	choiceRepo repository.IChoiceRepo
}

// NewCourseService 创建课程服务
func NewCourseService(courseRepo repository.ICourseRepo, bindRepo repository.IBindRepo, choiceRepo repository.IChoiceRepo) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		bindRepo:   bindRepo,
		choiceRepo: choiceRepo,
	}
}

// Create 创建课程
func (s *CourseService) Create(ctx context.Context, req *model.CreateCourseRequest) (*model.Course, error) {
	course := &model.Course{
		Name:        req.Name,
		Capacity:    req.Cap,
		CapSelected: 0,
		TeacherID:   nil,
	}

	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// Get 获取课程
func (s *CourseService) Get(ctx context.Context, courseID string) (*model.Course, error) {
	id, err := strconv.Atoi(courseID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}

	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errcode.CourseNotExisted
	}

	return course, nil
}

// List 获取课程列表
func (s *CourseService) List(ctx context.Context, offset, limit int) ([]*model.Course, error) {
	return s.courseRepo.List(ctx, offset, limit)
}

// BindCourse 绑定课程到教师
func (s *CourseService) BindCourse(ctx context.Context, courseID, teacherID string) error {
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		return errcode.ParamInvalid
	}
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// 检查课程是否存在
	course, err := s.courseRepo.GetByID(ctx, cID)
	if err != nil {
		return err
	}
	if course == nil {
		return errcode.CourseNotExisted
	}
	if course.TeacherID != nil {
		return errcode.CourseHasBound
	}

	// 检查课程是否已被其他教师绑定
	existingTeacher, err := s.bindRepo.GetByCourseID(ctx, cID)
	if err != nil {
		return err
	}
	if existingTeacher != nil {
		return errcode.CourseHasBound
	}

	// 创建绑定
	bind := &model.Bind{
		TeacherID: tID,
		CourseID:  cID,
	}
	if err := s.bindRepo.Create(ctx, bind); err != nil {
		return err
	}

	// 更新课程的教师ID
	return s.courseRepo.Update(ctx, cID, map[string]interface{}{
		"teacher_id": tID,
	})
}

// UnbindCourse 解绑课程
func (s *CourseService) UnbindCourse(ctx context.Context, courseID, teacherID string) error {
	cID, err := strconv.Atoi(courseID)
	if err != nil {
		return errcode.ParamInvalid
	}
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return errcode.ParamInvalid
	}

	// 检查绑定是否存在
	teacher, err := s.bindRepo.GetByCourseID(ctx, cID)
	if err != nil {
		return err
	}
	if teacher == nil {
		return errcode.CourseNotBind
	}
	if *teacher != tID {
		return errcode.PermDenied
	}

	// 删除绑定
	if err := s.bindRepo.DeleteByCourseID(ctx, cID); err != nil {
		return err
	}

	// 更新课程的教师ID为nil
	return s.courseRepo.Update(ctx, cID, map[string]interface{}{
		"teacher_id": nil,
	})
}

// GetTeacherCourses 获取教师的课程列表
func (s *CourseService) GetTeacherCourses(ctx context.Context, teacherID string) ([]*model.Course, error) {
	tID, err := strconv.Atoi(teacherID)
	if err != nil {
		return nil, errcode.ParamInvalid
	}
	return s.bindRepo.GetByTeacherID(ctx, tID)
}

// IsCourseExist 检查课程是否存在
func (s *CourseService) IsCourseExist(ctx context.Context, courseID string) (bool, error) {
	id, err := strconv.Atoi(courseID)
	if err != nil {
		return false, nil
	}
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	return course != nil, nil
}
