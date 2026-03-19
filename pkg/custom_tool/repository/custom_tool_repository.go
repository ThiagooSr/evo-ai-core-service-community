package repository

import (
	"context"
	"evo-ai-core-service/pkg/custom_tool/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CustomToolRepository interface {
	Create(ctx context.Context, customTool model.CustomTool) (*model.CustomTool, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.CustomTool, error)
	ListByAccountID(ctx context.Context, request model.CustomToolListRequest) ([]*model.CustomTool, error)
	CountByAccountID(ctx context.Context, request model.CustomToolListRequest) (int64, error)
	GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.CustomTool, error)
	Update(ctx context.Context, customTool *model.CustomTool, id uuid.UUID) (*model.CustomTool, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type customToolRepository struct {
	db *gorm.DB
}

func NewCustomToolRepository(db *gorm.DB) CustomToolRepository {
	return &customToolRepository{db: db}
}

func (r *customToolRepository) Create(ctx context.Context, customTool model.CustomTool) (*model.CustomTool, error) {
	if err := r.db.WithContext(ctx).Create(&customTool).Error; err != nil {
		return nil, err
	}

	return &customTool, nil
}

func (r *customToolRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CustomTool, error) {
	var customTool model.CustomTool

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&customTool).Error; err != nil {
		return nil, err
	}

	return &customTool, nil
}

func (r *customToolRepository) ListByAccountID(ctx context.Context, request model.CustomToolListRequest) ([]*model.CustomTool, error) {
	var customTool []*model.CustomTool

	query := r.db.WithContext(ctx).Where("account_id = ?", request.AccountID)

	if request.Search != "" {
		query = query.Where("name ILIKE ?", "%"+request.Search+"%")
	}

	if request.Tags != "" {
		query = query.Where("tags && ?", pq.Array(strings.Split(request.Tags, ",")))
	}

	if err := query.Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&customTool).Error; err != nil {
		return []*model.CustomTool{}, err
	}

	return customTool, nil
}

func (r *customToolRepository) CountByAccountID(ctx context.Context, request model.CustomToolListRequest) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.CustomTool{}).Where("account_id = ?", request.AccountID)

	if request.Search != "" {
		query = query.Where("name ILIKE ?", "%"+request.Search+"%")
	}

	if request.Tags != "" {
		query = query.Where("tags && ?", pq.Array(strings.Split(request.Tags, ",")))
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *customToolRepository) GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.CustomTool, error) {
	var customTool model.CustomTool

	if err := r.db.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&customTool).Error; err != nil {
		return nil, err
	}

	return &customTool, nil
}

func (r *customToolRepository) Update(ctx context.Context, customTool *model.CustomTool, id uuid.UUID) (*model.CustomTool, error) {
	customTool.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(customTool).First(&customTool).Error; err != nil {
		return nil, err
	}

	return customTool, nil
}

func (r *customToolRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := r.db.WithContext(ctx).Model(&model.CustomTool{}).Where("id = ?", id).Delete(&model.CustomTool{}).Error; err != nil {
		return false, err
	}

	return true, nil
}
