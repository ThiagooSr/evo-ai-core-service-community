package repository

import (
	"context"
	"evo-ai-core-service/pkg/folder/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FolderRepository interface {
	Create(ctx context.Context, folder model.Folder) (*model.Folder, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error)
	ListByAccountID(ctx context.Context, accountId uuid.UUID, page int, pageSize int) ([]*model.Folder, error)
	CountByAccountID(ctx context.Context, accountId uuid.UUID) (int64, error)
	Update(ctx context.Context, folder *model.Folder, id uuid.UUID) (*model.Folder, error)
	Delete(ctx context.Context, id uuid.UUID) (bool, error)
	GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.Folder, error)
}

type folderRepository struct {
	db *gorm.DB
}

func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

func (r *folderRepository) Create(ctx context.Context, folder model.Folder) (*model.Folder, error) {
	if err := r.db.WithContext(ctx).Create(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *folderRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Folder, error) {
	var folder model.Folder

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}

func (r *folderRepository) ListByAccountID(ctx context.Context, accountId uuid.UUID, page int, pageSize int) ([]*model.Folder, error) {
	var folders []*model.Folder

	query := r.db.WithContext(ctx).Where("account_id = ?", accountId)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&folders).Error; err != nil {
		return []*model.Folder{}, err
	}

	return folders, nil
}

func (r *folderRepository) CountByAccountID(ctx context.Context, accountId uuid.UUID) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Folder{}).Where("account_id = ?", accountId).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *folderRepository) Update(ctx context.Context, folder *model.Folder, id uuid.UUID) (*model.Folder, error) {
	folder.UpdatedAt = time.Now()
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(folder).First(&folder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *folderRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	if err := r.db.WithContext(ctx).Model(&model.Folder{}).Where("id = ?", id).Delete(&model.Folder{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *folderRepository) GetByIDAndAccountID(ctx context.Context, id uuid.UUID, accountID uuid.UUID) (*model.Folder, error) {
	var folder model.Folder

	if err := r.db.WithContext(ctx).Where("id = ? AND account_id = ?", id, accountID).First(&folder).Error; err != nil {
		return nil, err
	}

	return &folder, nil
}
