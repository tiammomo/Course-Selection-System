package service

import (
	"context"
	"strconv"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"
	"course_select/internal/infrastructure/encrypt"
	"course_select/internal/pkg/errcode"

	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
)

// AuthService 认证服务
type AuthService struct {
	memberRepo repository.IMemberRepo
	sessionKey string
	cookieName string
	expireHours int
}

// NewAuthService 创建认证服务
func NewAuthService(memberRepo repository.IMemberRepo, sessionKey, cookieName string, expireHours int) *AuthService {
	return &AuthService{
		memberRepo:  memberRepo,
		sessionKey: sessionKey,
		cookieName: cookieName,
		expireHours: expireHours,
	}
}

// Login 登录
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.Member, error) {
	member, err := s.memberRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errcode.UserNotExisted
	}
	if member.IsDeleted {
		return nil, errcode.UserHasDeleted
	}
	if !encrypt.ComparePassword(password, member.Password) {
		return nil, errcode.WrongPassword
	}
	return member, nil
}

// Logout 登出
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	// Session 管理由中间件处理
	return nil
}

// GetMemberBySession 从 Session 获取成员
func (s *AuthService) GetMemberBySession(ctx context.Context, session sessions.Session) *model.Member {
	v := session.Get(s.sessionKey)
	if v == nil {
		return nil
	}
	data, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	userIDVal, ok := data["user_id"]
	if !ok {
		return nil
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return nil
	}
	id := strToInt(userID)
	if id <= 0 {
		return nil
	}
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil
	}
	return member
}

// GenerateSessionID 生成 Session ID
func (s *AuthService) GenerateSessionID() string {
	return uuid.New().String()
}

// CreateSession 创建 Session
func (s *AuthService) CreateSession(session sessions.Session, member *model.Member) {
	session.Set(s.sessionKey, map[string]interface{}{
		"user_id":   strconv.Itoa(member.UserID),
		"nickname":  member.Nickname,
		"username":  member.Username,
		"user_type": int(member.UserType),
	})
	session.Options(sessions.Options{
		MaxAge:   s.expireHours * 3600,
		HttpOnly: true,
	})
}

// ValidateMember 验证成员是否存在且未删除
func (s *AuthService) ValidateMember(ctx context.Context, userID string) (*model.Member, error) {
	id := strToInt(userID)
	if id <= 0 {
		return nil, errcode.ParamInvalid.WithMsg("无效的用户ID")
	}
	member, err := s.memberRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, errcode.UserNotExisted
	}
	if member.IsDeleted {
		return nil, errcode.UserHasDeleted
	}
	return member, nil
}

// strToInt string 转 int
func strToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
