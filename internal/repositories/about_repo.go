package repositories

import (
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type AboutRepository interface {
	CreateAbout(about *model.About) error
	CreateAboutDetail(aboutDetail *model.AboutDetail) error
	GetAllAbout() ([]model.About, error)
	GetAboutByID(id string) (*model.About, error)
	GetAboutDetailByID(id string) (*model.AboutDetail, error)
	UpdateAbout(id string, about *model.About) (*model.About, error)
	UpdateAboutDetail(id string, aboutDetail *model.AboutDetail) (*model.AboutDetail, error)
	DeleteAbout(id string) error
	DeleteAboutDetail(id string) error
}

type aboutRepository struct {
	DB *gorm.DB
}

func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepository{DB: db}
}

func (r *aboutRepository) CreateAbout(about *model.About) error {
	if err := r.DB.Create(&about).Error; err != nil {
		return fmt.Errorf("failed to create About: %v", err)
	}
	return nil
}

func (r *aboutRepository) CreateAboutDetail(aboutDetail *model.AboutDetail) error {
	if err := r.DB.Create(&aboutDetail).Error; err != nil {
		return fmt.Errorf("failed to create AboutDetail: %v", err)
	}
	return nil
}

func (r *aboutRepository) GetAllAbout() ([]model.About, error) {
	var abouts []model.About
	if err := r.DB.Find(&abouts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all About records: %v", err)
	}
	return abouts, nil
}

func (r *aboutRepository) GetAboutByID(id string) (*model.About, error) {
	var about model.About
	if err := r.DB.Preload("AboutDetail").Where("id = ?", id).First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("about with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch About by ID: %v", err)
	}
	return &about, nil
}

func (r *aboutRepository) GetAboutDetailByID(id string) (*model.AboutDetail, error) {
	var aboutDetail model.AboutDetail
	if err := r.DB.Where("id = ?", id).First(&aboutDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("aboutdetail with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to fetch About by ID: %v", err)
	}
	return &aboutDetail, nil
}

func (r *aboutRepository) UpdateAbout(id string, about *model.About) (*model.About, error) {
	if err := r.DB.Model(&about).Where("id = ?", id).Updates(about).Error; err != nil {
		return nil, fmt.Errorf("failed to update About: %v", err)
	}
	return about, nil
}

func (r *aboutRepository) UpdateAboutDetail(id string, aboutDetail *model.AboutDetail) (*model.AboutDetail, error) {
	if err := r.DB.Model(&aboutDetail).Where("id = ?", id).Updates(aboutDetail).Error; err != nil {
		return nil, fmt.Errorf("failed to update AboutDetail: %v", err)
	}
	return aboutDetail, nil
}

func (r *aboutRepository) DeleteAbout(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&model.About{}).Error; err != nil {
		return fmt.Errorf("failed to delete About: %v", err)
	}
	return nil
}

func (r *aboutRepository) DeleteAboutDetail(id string) error {
	if err := r.DB.Where("id = ?", id).Delete(&model.AboutDetail{}).Error; err != nil {
		return fmt.Errorf("failed to delete AboutDetail: %v", err)
	}
	return nil
}
