package database

import (
	"context"
	"fmt"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"

	"gorm.io/gorm"
)

// CourseRepoImpl 课程仓储实现
type CourseRepoImpl struct {
	db *gorm.DB
}

// NewCourseRepo 创建课程仓储
func NewCourseRepo(db *gorm.DB) repository.ICourseRepo {
	return &CourseRepoImpl{db: db}
}

func (r *CourseRepoImpl) Create(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *CourseRepoImpl) GetByID(ctx context.Context, id int) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).First(&course, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepoImpl) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Course{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("course not found")
	}
	return nil
}

func (r *CourseRepoImpl) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Course{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("course not found")
	}
	return nil
}

func (r *CourseRepoImpl) List(ctx context.Context, offset, limit int) ([]*model.Course, error) {
	var courses []*model.Course
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&courses).Error
	return courses, err
}

func (r *CourseRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Course{}).Count(&count)
	return count, result.Error
}

// BindRepoImpl 绑定仓储实现
type BindRepoImpl struct {
	db *gorm.DB
}

// NewBindRepo 创建绑定仓储
func NewBindRepo(db *gorm.DB) repository.IBindRepo {
	return &BindRepoImpl{db: db}
}

func (r *BindRepoImpl) Create(ctx context.Context, bind *model.Bind) error {
	return r.db.WithContext(ctx).Create(bind).Error
}

func (r *BindRepoImpl) GetByTeacherID(ctx context.Context, teacherID int) ([]*model.Course, error) {
	var courses []*model.Course
	err := r.db.WithContext(ctx).
		Table("course").
		Joins("JOIN bind ON course.id = bind.course_id").
		Where("bind.teacher_id = ?", teacherID).
		Find(&courses).Error
	return courses, err
}

func (r *BindRepoImpl) GetByCourseID(ctx context.Context, courseID int) (*int, error) {
	var bind model.Bind
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).First(&bind).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &bind.TeacherID, nil
}

func (r *BindRepoImpl) DeleteByCourseID(ctx context.Context, courseID int) error {
	result := r.db.WithContext(ctx).Where("course_id = ?", courseID).Delete(&model.Bind{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *BindRepoImpl) DeleteByTeacherID(ctx context.Context, teacherID int) error {
	result := r.db.WithContext(ctx).Where("teacher_id = ?", teacherID).Delete(&model.Bind{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// ChoiceRepoImpl 选课仓储实现
type ChoiceRepoImpl struct {
	db *gorm.DB
}

// NewChoiceRepo 创建选课仓储
func NewChoiceRepo(db *gorm.DB) repository.IChoiceRepo {
	return &ChoiceRepoImpl{db: db}
}

func (r *ChoiceRepoImpl) Create(ctx context.Context, choice *model.Choice) error {
	return r.db.WithContext(ctx).Create(choice).Error
}

func (r *ChoiceRepoImpl) Delete(ctx context.Context, studentID, courseID int) error {
	result := r.db.WithContext(ctx).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Delete(&model.Choice{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ChoiceRepoImpl) GetByStudentID(ctx context.Context, studentID int) ([]*model.Course, error) {
	var courses []*model.Course
	err := r.db.WithContext(ctx).
		Table("course").
		Joins("JOIN choice ON course.id = choice.course_id").
		Where("choice.student_id = ?", studentID).
		Find(&courses).Error
	return courses, err
}

func (r *ChoiceRepoImpl) GetByCourseID(ctx context.Context, courseID int) ([]int, error) {
	var studentIDs []int
	err := r.db.WithContext(ctx).
		Model(&model.Choice{}).
		Where("course_id = ?", courseID).
		Pluck("student_id", &studentIDs).Error
	return studentIDs, err
}

func (r *ChoiceRepoImpl) Exists(ctx context.Context, studentID, courseID int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Choice{}).
		Where("student_id = ? AND course_id = ?", studentID, courseID).
		Count(&count).Error
	return count > 0, err
}

func (r *ChoiceRepoImpl) CountByCourseID(ctx context.Context, courseID int) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Choice{}).
		Where("course_id = ?", courseID).
		Count(&count).Error
	return int(count), err
}
