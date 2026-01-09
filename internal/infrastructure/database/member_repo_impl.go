package database

import (
	"context"
	"fmt"

	"course_select/internal/domain/model"
	"course_select/internal/domain/repository"

	"gorm.io/gorm"
)

// MemberRepoImpl 成员仓储实现
type MemberRepoImpl struct {
	db *gorm.DB
}

// NewMemberRepo 创建成员仓储
func NewMemberRepo(db *gorm.DB) repository.IMemberRepo {
	return &MemberRepoImpl{db: db}
}

func (r *MemberRepoImpl) Create(ctx context.Context, member *model.Member) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *MemberRepoImpl) GetByID(ctx context.Context, id int) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).First(&member, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepoImpl) GetByUsername(ctx context.Context, username string) (*model.Member, error) {
	var member model.Member
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&member).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepoImpl) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Member{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

func (r *MemberRepoImpl) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Model(&model.Member{}).Where("id = ?", id).Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("member not found")
	}
	return nil
}

func (r *MemberRepoImpl) List(ctx context.Context, offset, limit int) ([]*model.Member, error) {
	var members []*model.Member
	err := r.db.WithContext(ctx).
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Find(&members).Error
	return members, err
}

func (r *MemberRepoImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Member{}).Where("is_deleted = ?", false).Count(&count)
	return count, err.Error
}
