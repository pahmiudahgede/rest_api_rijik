package repositories

import (
	"errors"
	"fmt"
	"log"
	"rijig/model"

	"gorm.io/gorm"
)

type IdentityCardRepository interface {
	CreateIdentityCard(identityCard *model.IdentityCard) (*model.IdentityCard, error)
	GetIdentityCardByID(id string) (*model.IdentityCard, error)
	GetIdentityCardsByUserID(userID string) ([]model.IdentityCard, error)
	UpdateIdentityCard(id string, updatedCard *model.IdentityCard) (*model.IdentityCard, error)
	DeleteIdentityCard(id string) error
}

type identityCardRepository struct {
	db *gorm.DB
}

func NewIdentityCardRepository(db *gorm.DB) IdentityCardRepository {
	return &identityCardRepository{
		db: db,
	}
}

func (r *identityCardRepository) CreateIdentityCard(identityCard *model.IdentityCard) (*model.IdentityCard, error) {
	if err := r.db.Create(identityCard).Error; err != nil {
		log.Printf("Error creating identity card: %v", err)
		return nil, fmt.Errorf("failed to create identity card: %w", err)
	}
	return identityCard, nil
}

func (r *identityCardRepository) GetIdentityCardByID(id string) (*model.IdentityCard, error) {
	var identityCard model.IdentityCard
	if err := r.db.Where("id = ?", id).First(&identityCard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("identity card not found with id %s", id)
		}
		log.Printf("Error fetching identity card by ID: %v", err)
		return nil, fmt.Errorf("error fetching identity card by ID: %w", err)
	}
	return &identityCard, nil
}

func (r *identityCardRepository) GetIdentityCardsByUserID(userID string) ([]model.IdentityCard, error) {
	var identityCards []model.IdentityCard
	if err := r.db.Where("user_id = ?", userID).Find(&identityCards).Error; err != nil {
		log.Printf("Error fetching identity cards by userID: %v", err)
		return nil, fmt.Errorf("error fetching identity cards by userID: %w", err)
	}
	return identityCards, nil
}

func (r *identityCardRepository) UpdateIdentityCard(id string, updatedCard *model.IdentityCard) (*model.IdentityCard, error) {
	var existingCard model.IdentityCard
	if err := r.db.Where("id = ?", id).First(&existingCard).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("identity card with ID %s not found", id)
		}
		log.Printf("Error fetching identity card for update: %v", err)
		return nil, fmt.Errorf("error fetching identity card for update: %w", err)
	}

	if err := r.db.Save(&existingCard).Error; err != nil {
		log.Printf("Error updating identity card: %v", err)
		return nil, fmt.Errorf("failed to update identity card: %w", err)
	}
	return &existingCard, nil
}

func (r *identityCardRepository) DeleteIdentityCard(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&model.IdentityCard{}).Error; err != nil {
		log.Printf("Error deleting identity card: %v", err)
		return fmt.Errorf("failed to delete identity card: %w", err)
	}
	return nil
}
