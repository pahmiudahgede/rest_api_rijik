package about

import (
	"context"
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type AboutRepository interface {
	CreateAbout(ctx context.Context, about *model.About) error
	CreateAboutDetail(ctx context.Context, aboutDetail *model.AboutDetail) error
	GetAllAbout(ctx context.Context) ([]model.About, error)
	GetAboutByID(ctx context.Context, id string) (*model.About, error)
	GetAboutByIDWithoutPrel(ctx context.Context, id string) (*model.About, error)
	GetAboutDetailByID(ctx context.Context, id string) (*model.AboutDetail, error)
	UpdateAbout(ctx context.Context, id string, about *model.About) (*model.About, error)
	UpdateAboutDetail(ctx context.Context, id string, aboutDetail *model.AboutDetail) (*model.AboutDetail, error)
	DeleteAbout(ctx context.Context, id string) error
	DeleteAboutDetail(ctx context.Context, id string) error
}

type aboutRepository struct {
	db *gorm.DB
}

func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepository{db}
}

func (r *aboutRepository) CreateAbout(ctx context.Context, about *model.About) error {
	if err := r.db.WithContext(ctx).Create(&about).Error; err != nil {
		return fmt.Errorf("failed to create About: %v", err)
	}
	return nil
}

func (r *aboutRepository) CreateAboutDetail(ctx context.Context, aboutDetail *model.AboutDetail) error {
	if err := r.db.WithContext(ctx).Create(&aboutDetail).Error; err != nil {
		return fmt.Errorf("failed to create AboutDetail: %v", err)
	}
	return nil
}

func (r *aboutRepository) GetAllAbout(ctx context.Context) ([]model.About, error) {
	var abouts []model.About
	if err := r.db.WithContext(ctx).Find(&abouts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all About records: %v", err)
	}
	return abouts, nil
}

func (r *aboutRepository) GetAboutByID(ctx context.Context, id string) (*model.About, error) {
	var about model.About
	if err := r.db.WithContext(ctx).Preload("AboutDetail").Where("id = ?", id).First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("about with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch About by ID: %v", err)
	}
	return &about, nil
}

func (r *aboutRepository) GetAboutByIDWithoutPrel(ctx context.Context, id string) (*model.About, error) {
	var about model.About
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("about with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch About by ID: %v", err)
	}
	return &about, nil
}

func (r *aboutRepository) GetAboutDetailByID(ctx context.Context, id string) (*model.AboutDetail, error) {
	var aboutDetail model.AboutDetail
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&aboutDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("aboutdetail with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch AboutDetail by ID: %v", err)
	}
	return &aboutDetail, nil
}

func (r *aboutRepository) UpdateAbout(ctx context.Context, id string, about *model.About) (*model.About, error) {
	if err := r.db.WithContext(ctx).Model(&about).Where("id = ?", id).Updates(about).Error; err != nil {
		return nil, fmt.Errorf("failed to update About: %v", err)
	}
	return about, nil
}

func (r *aboutRepository) UpdateAboutDetail(ctx context.Context, id string, aboutDetail *model.AboutDetail) (*model.AboutDetail, error) {
	if err := r.db.WithContext(ctx).Model(&aboutDetail).Where("id = ?", id).Updates(aboutDetail).Error; err != nil {
		return nil, fmt.Errorf("failed to update AboutDetail: %v", err)
	}
	return aboutDetail, nil
}

func (r *aboutRepository) DeleteAbout(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.About{}).Error; err != nil {
		return fmt.Errorf("failed to delete About: %v", err)
	}
	return nil
}

func (r *aboutRepository) DeleteAboutDetail(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.AboutDetail{}).Error; err != nil {
		return fmt.Errorf("failed to delete AboutDetail: %v", err)
	}
	return nil
}
