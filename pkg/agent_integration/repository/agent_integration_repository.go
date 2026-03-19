package repository

import (
	"context"
	"evo-ai-core-service/pkg/agent_integration/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AgentIntegrationRepository interface {
	Upsert(ctx context.Context, integration model.AgentIntegration) (*model.AgentIntegration, error)
	GetByAccountAgentAndProvider(ctx context.Context, accountID, agentID uuid.UUID, provider string) (*model.AgentIntegration, error)
	ListByAccountAndAgent(ctx context.Context, accountID, agentID uuid.UUID) ([]*model.AgentIntegration, error)
	Delete(ctx context.Context, accountID, agentID uuid.UUID, provider string) error
}

type agentIntegrationRepository struct {
	db *gorm.DB
}

func NewAgentIntegrationRepository(db *gorm.DB) AgentIntegrationRepository {
	return &agentIntegrationRepository{db: db}
}

func (r *agentIntegrationRepository) Upsert(ctx context.Context, integration model.AgentIntegration) (*model.AgentIntegration, error) {
	integration.UpdatedAt = time.Now()

	// Use GORM's upsert functionality
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_id"}, {Name: "agent_id"}, {Name: "provider"}},
		DoUpdates: clause.AssignmentColumns([]string{"config", "updated_at"}),
	}).Create(&integration).Error

	if err != nil {
		return nil, err
	}

	return &integration, nil
}

func (r *agentIntegrationRepository) GetByAccountAgentAndProvider(ctx context.Context, accountID, agentID uuid.UUID, provider string) (*model.AgentIntegration, error) {
	var integration model.AgentIntegration

	err := r.db.WithContext(ctx).
		Where("account_id = ? AND agent_id = ? AND provider = ?", accountID, agentID, provider).
		First(&integration).Error

	if err != nil {
		return nil, err
	}

	return &integration, nil
}

func (r *agentIntegrationRepository) ListByAccountAndAgent(ctx context.Context, accountID, agentID uuid.UUID) ([]*model.AgentIntegration, error) {
	var integrations []*model.AgentIntegration

	err := r.db.WithContext(ctx).
		Where("account_id = ? AND agent_id = ?", accountID, agentID).
		Find(&integrations).Error

	if err != nil {
		return nil, err
	}

	return integrations, nil
}

func (r *agentIntegrationRepository) Delete(ctx context.Context, accountID, agentID uuid.UUID, provider string) error {
	return r.db.WithContext(ctx).
		Where("account_id = ? AND agent_id = ? AND provider = ?", accountID, agentID, provider).
		Delete(&model.AgentIntegration{}).Error
}
