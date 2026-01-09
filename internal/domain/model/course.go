package model

import (
	"course_select/internal/pkg/errcode"
	"time"

	"gorm.io/gorm"
)

// Course 课程实体
type Course struct {
	gorm.Model
	CourseID    int    `gorm:"primaryKey;autoIncrement" json:"course_id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	Capacity    int    `gorm:"not null" json:"capacity"`      // 课程容量
	CapSelected int    `gorm:"default:0;not null" json:"cap_selected"` // 已选人数
	TeacherID   *int   `gorm:"default:null;index" json:"teacher_id"`   // 授课教师

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "course"
}

// ToResponse 转换为响应结构
func (c *Course) ToResponse() *CourseResponse {
	if c == nil {
		return nil
	}
	var teacherID string
	if c.TeacherID != nil {
		teacherID = intToString(*c.TeacherID)
	}
	return &CourseResponse{
		CourseID:  intToString(c.CourseID),
		Name:      c.Name,
		Capacity:  c.Capacity,
		TeacherID: teacherID,
	}
}

// CourseResponse 课程响应
type CourseResponse struct {
	CourseID  string `json:"course_id"`
	Name      string `json:"name"`
	Capacity  int    `json:"capacity"`
	TeacherID string `json:"teacher_id,omitempty"`
}

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
	Cap  int    `json:"cap" binding:"required,min=1"`
}

// Validate 验证请求
func (r *CreateCourseRequest) Validate() error {
	if r.Cap <= 0 {
		return errcode.ParamInvalid.WithMsg("课程容量必须大于 0")
	}
	return nil
}

// BindCourseRequest 绑定课程请求
type BindCourseRequest struct {
	CourseID  string `json:"course_id" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
}

// UnbindCourseRequest 解绑课程请求
type UnbindCourseRequest struct {
	CourseID  string `json:"course_id" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
}

// ScheduleCourseRequest 排课请求
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string `json:"teacher_course_relationship" binding:"required"`
}
