package identitycart

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rijig/model"

	"gorm.io/gorm"
)

type IdentityCardRepository interface {
	CreateIdentityCard(ctx context.Context, identityCard *model.IdentityCard) (*model.IdentityCard, error)
	GetIdentityCardByID(ctx context.Context, id string) (*model.IdentityCard, error)
	GetIdentityCardsByUserID(ctx context.Context, userID string) ([]model.IdentityCard, error)
	UpdateIdentityCard(ctx context.Context, identity *model.IdentityCard) error
}

type identityCardRepository struct {
	db *gorm.DB
}

func NewIdentityCardRepository(db *gorm.DB) IdentityCardRepository {
	return &identityCardRepository{
		db: db,
	}
}

func (r *identityCardRepository) CreateIdentityCard(ctx context.Context, identityCard *model.IdentityCard) (*model.IdentityCard, error) {
	if err := r.db.WithContext(ctx).Create(identityCard).Error; err != nil {
		log.Printf("Error creating identity card: %v", err)
		return nil, fmt.Errorf("failed to create identity card: %w", err)
	}
	return identityCard, nil
}

func (r *identityCardRepository) GetIdentityCardByID(ctx context.Context, id string) (*model.IdentityCard, error) {
	var identityCard model.IdentityCard
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&identityCard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("identity card not found with id %s", id)
		}
		log.Printf("Error fetching identity card by ID: %v", err)
		return nil, fmt.Errorf("error fetching identity card by ID: %w", err)
	}
	return &identityCard, nil
}

func (r *identityCardRepository) GetIdentityCardsByUserID(ctx context.Context, userID string) ([]model.IdentityCard, error) {
	var identityCards []model.IdentityCard
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&identityCards).Error; err != nil {
		log.Printf("Error fetching identity cards by userID: %v", err)
		return nil, fmt.Errorf("error fetching identity cards by userID: %w", err)
	}
	return identityCards, nil
}

func (r *identityCardRepository) UpdateIdentityCard(ctx context.Context, identity *model.IdentityCard) error {
	return r.db.WithContext(ctx).
		Model(&model.IdentityCard{}).
		Where("user_id = ?", identity.UserID).
		Updates(identity).Error
}
