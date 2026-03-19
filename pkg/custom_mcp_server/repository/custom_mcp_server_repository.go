package repository

import (
	"context"
	"evo-ai-core-service/pkg/custom_mcp_server/model"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CustomMcpServerRepository interface {
	Create(ctx context.Context, customMcpServer model.CustomMcpServer) (*model.CustomMcpServer, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.CustomMcpServer, error)
	ListByAccountID(ctx context.Context, request model.CustomMcpServerListRequest) ([]*model.CustomMcpServer, error)
	CountByAccountID(ctx context.Context, request model.CustomMcpServerListRequest) (int64, error)
	Update(ctx context.Context, customMcpServer *model.CustomMcpServer, id uuid.UUID) (*model.CustomMcpServer, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetByAgentConfig(ctx context.Context, accountID uuid.UUID, serverIDs []uuid.UUID) ([]*model.CustomMcpServer, error)
	GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.CustomMcpServer, error)
}

type customMcpServerRepository struct {
	db *gorm.DB
}

func NewCustomMcpServerRepository(db *gorm.DB) CustomMcpServerRepository {
	return &customMcpServerRepository{db: db}
}

func (r *customMcpServerRepository) Create(ctx context.Context, customMcpServer model.CustomMcpServer) (*model.CustomMcpServer, error) {
	if err := r.db.WithContext(ctx).Create(&customMcpServer).Error; err != nil {
		return nil, err
	}

	return &customMcpServer, nil
}

func (r *customMcpServerRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.CustomMcpServer, error) {
	var customMcpServer model.CustomMcpServer

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&customMcpServer).Error; err != nil {
		return nil, err
	}

	return &customMcpServer, nil
}

func (r *customMcpServerRepository) ListByAccountID(ctx context.Context, request model.CustomMcpServerListRequest) ([]*model.CustomMcpServer, error) {
	var customMcpServer []*model.CustomMcpServer

	query := r.db.WithContext(ctx).Where("account_id = ?", request.AccountID)

	if request.Search != "" {
		query = query.Where("name ILIKE ?", "%"+request.Search+"%")
	}

	if request.Tags != "" {
		query = query.Where("tags && ?", pq.Array(strings.Split(request.Tags, ",")))
	}

	if err := query.Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&customMcpServer).Error; err != nil {
		return []*model.CustomMcpServer{}, err
	}

	return customMcpServer, nil
}

func (r *customMcpServerRepository) CountByAccountID(ctx context.Context, request model.CustomMcpServerListRequest) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.CustomMcpServer{}).Where("account_id = ?", request.AccountID)

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

func (r *customMcpServerRepository) Update(ctx context.Context, customMcpServer *model.CustomMcpServer, id uuid.UUID) (*model.CustomMcpServer, error) {
	customMcpServer.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(customMcpServer).First(&customMcpServer).Error; err != nil {
		return nil, err
	}

	return customMcpServer, nil
}

func (r *customMcpServerRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := r.db.WithContext(ctx).Model(&model.CustomMcpServer{}).Where("id = ?", id).Delete(&model.CustomMcpServer{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *customMcpServerRepository) GetByAgentConfig(ctx context.Context, accountID uuid.UUID, serverIDs []uuid.UUID) ([]*model.CustomMcpServer, error) {
	var customMcpServer []*model.CustomMcpServer

	if err := r.db.WithContext(ctx).Where("account_id = ? AND id IN (?)", accountID, serverIDs).Find(&customMcpServer).Error; err != nil {
		return nil, err
	}

	return customMcpServer, nil
}

func (r *customMcpServerRepository) GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.CustomMcpServer, error) {
	var customMcpServer model.CustomMcpServer

	if err := r.db.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&customMcpServer).Error; err != nil {
		return nil, err
	}

	return &customMcpServer, nil
}
