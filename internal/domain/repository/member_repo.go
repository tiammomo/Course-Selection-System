package repository

import (
	"context"

	"course_select/internal/domain/model"
)

// IMemberRepo 成员仓储接口
type IMemberRepo interface {
	Create(ctx context.Context, member *model.Member) error
	GetByID(ctx context.Context, id int) (*model.Member, error)
	GetByUsername(ctx context.Context, username string) (*model.Member, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*model.Member, error)
	Count(ctx context.Context) (int64, error)
}
