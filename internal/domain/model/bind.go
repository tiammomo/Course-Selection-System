package model

import (
	"time"
)

// Bind 教师课程绑定实体
type Bind struct {
	TeacherID int `gorm:"primaryKey"`
	CourseID  int `gorm:"primaryKey"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Bind) TableName() string {
	return "bind"
}

// Choice 学生选课实体
type Choice struct {
	StudentID int `gorm:"primaryKey"`
	CourseID  int `gorm:"primaryKey"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Choice) TableName() string {
	return "choice"
}
