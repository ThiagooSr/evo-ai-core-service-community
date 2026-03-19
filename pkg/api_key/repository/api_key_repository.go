package repository

import (
	"context"
	"evo-ai-core-service/pkg/api_key/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey model.ApiKey) (*model.ApiKey, error)
	GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.ApiKey, error)
	ListByAccountID(ctx context.Context, request model.ApiKeyListRequest) ([]*model.ApiKey, error)
	CountByAccountID(ctx context.Context, accountID uuid.UUID, active string) (int64, error)
	Update(ctx context.Context, apiKey *model.ApiKey, id uuid.UUID) (*model.ApiKey, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(ctx context.Context, apiKey model.ApiKey) (*model.ApiKey, error) {
	if err := r.db.WithContext(ctx).Create(&apiKey).Error; err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *apiKeyRepository) GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.ApiKey, error) {
	var apiKey model.ApiKey

	if err := r.db.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&apiKey).Error; err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (r *apiKeyRepository) ListByAccountID(ctx context.Context, request model.ApiKeyListRequest) ([]*model.ApiKey, error) {
	var apiKeys []*model.ApiKey

	query := r.db.WithContext(ctx).Where("account_id = ?", request.AccountID)

	// Filter by active status - default to active only
	if request.Active != "" {
		query = query.Where("is_active = ?", request.Active)
	} else {
		// Default: show only active API keys
		query = query.Where("is_active = ?", true)
	}

	if err := query.Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&apiKeys).Error; err != nil {
		return []*model.ApiKey{}, err
	}

	return apiKeys, nil
}

func (r *apiKeyRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID, active string) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&model.ApiKey{}).Where("account_id = ?", accountID)

	// Filter by active status - default to active only
	if active != "" {
		query = query.Where("is_active = ?", active)
	} else {
		// Default: count only active API keys
		query = query.Where("is_active = ?", true)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *apiKeyRepository) Update(ctx context.Context, apiKey *model.ApiKey, id uuid.UUID) (*model.ApiKey, error) {
	apiKey.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(apiKey).First(&apiKey).Error; err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (r *apiKeyRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := r.db.WithContext(ctx).Model(&model.ApiKey{}).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return false, err
	}

	return true, nil
}
