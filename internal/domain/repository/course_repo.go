package repository

import (
	"context"

	"course_select/internal/domain/model"
)

// ICourseRepo 课程仓储接口
type ICourseRepo interface {
	Create(ctx context.Context, course *model.Course) error
	GetByID(ctx context.Context, id int) (*model.Course, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*model.Course, error)
	Count(ctx context.Context) (int64, error)
}

// IBindRepo 绑定仓储接口
type IBindRepo interface {
	Create(ctx context.Context, bind *model.Bind) error
	GetByTeacherID(ctx context.Context, teacherID int) ([]*model.Course, error)
	GetByCourseID(ctx context.Context, courseID int) (*int, error) // 返回教师ID
	DeleteByCourseID(ctx context.Context, courseID int) error
	DeleteByTeacherID(ctx context.Context, teacherID int) error
}

// IChoiceRepo 选课仓储接口
type IChoiceRepo interface {
	Create(ctx context.Context, choice *model.Choice) error
	Delete(ctx context.Context, studentID, courseID int) error
	GetByStudentID(ctx context.Context, studentID int) ([]*model.Course, error)
	GetByCourseID(ctx context.Context, courseID int) ([]int, error) // 返回学生ID列表
	Exists(ctx context.Context, studentID, courseID int) (bool, error)
	CountByCourseID(ctx context.Context, courseID int) (int, error)
}
