package model

import (
	"course_select/internal/pkg/errcode"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// UserType 用户类型
type UserType int

const (
	UserTypeAdmin   UserType = 1
	UserTypeStudent UserType = 2
	UserTypeTeacher UserType = 3
)

// Member 成员实体
type Member struct {
	gorm.Model
	UserID    int      `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Username  string   `gorm:"size:20;uniqueIndex;not null" json:"username"`
	Password  string   `gorm:"size:50;not null" json:"-"`
	Nickname  string   `gorm:"size:20" json:"nickname"`
	UserType  UserType `gorm:"not null" json:"user_type"`
	IsDeleted bool     `gorm:"default:false;index" json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Member) TableName() string {
	return "member"
}

// BeforeCreate 创建前加密密码
func (m *Member) BeforeCreate(tx *gorm.DB) error {
	// 密码加密由应用层处理
	return nil
}

// IsAdmin 是否管理员
func (m *Member) IsAdmin() bool {
	return m.UserType == UserTypeAdmin
}

// IsStudent 是否学生
func (m *Member) IsStudent() bool {
	return m.UserType == UserTypeStudent
}

// IsTeacher 是否教师
func (m *Member) IsTeacher() bool {
	return m.UserType == UserTypeTeacher
}

// ToResponse 转换为响应结构
func (m *Member) ToResponse() *MemberResponse {
	if m == nil {
		return nil
	}
	return &MemberResponse{
		UserID:   intToString(int(m.UserID)),
		Nickname: m.Nickname,
		Username: m.Username,
		UserType: m.UserType,
	}
}

// MemberResponse 成员响应
type MemberResponse struct {
	UserID   string   `json:"user_id"`
	Nickname string   `json:"nickname"`
	Username string   `json:"username"`
	UserType UserType `json:"user_type"`
}

// CreateMemberRequest 创建成员请求
type CreateMemberRequest struct {
	Nickname string   `json:"nickname" binding:"required,min=4,max=20"`
	Username string   `json:"username" binding:"required,min=8,max=20,alpha"`
	Password string   `json:"password" binding:"required,min=8,max=20"`
	UserType UserType `json:"user_type" binding:"required"`
}

// Validate 验证请求
func (r *CreateMemberRequest) Validate() error {
	if r.UserType != UserTypeAdmin && r.UserType != UserTypeStudent && r.UserType != UserTypeTeacher {
		return errcode.ParamInvalid.WithMsg("UserType 必须为 1(管理员)、2(学生) 或 3(教师)")
	}
	return nil
}

// UpdateMemberRequest 更新成员请求
type UpdateMemberRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Nickname string `json:"nickname" binding:"required,min=4,max=20"`
}

// DeleteMemberRequest 删除成员请求
type DeleteMemberRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

// intToString int 转 string
func intToString(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}
