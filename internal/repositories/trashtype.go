package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetTrashCategories() ([]domain.TrashCategory, error) {
	var categories []domain.TrashCategory
	if err := config.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func GetTrashCategoryDetail(id string) (domain.TrashCategory, error) {
	var category domain.TrashCategory
	if err := config.DB.Preload("Details").Where("id = ?", id).First(&category).Error; err != nil {
		return category, err
	}
	return category, nil
}

func GetTrashDetailByID(id string) (domain.TrashDetail, error) {
	var detail domain.TrashDetail

	if err := config.DB.Where("id = ?", id).First(&detail).Error; err != nil {
		return detail, err
	}
	return detail, nil
}

func CreateTrashCategory(category *domain.TrashCategory) error {
	if err := config.DB.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func CreateTrashDetail(detail *domain.TrashDetail) error {
	if err := config.DB.Create(detail).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTrashCategory(category *domain.TrashCategory) error {
	if err := config.DB.Save(category).Error; err != nil {
		return err
	}
	return nil
}

func UpdateTrashDetail(detail *domain.TrashDetail) error {
	if err := config.DB.Save(detail).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTrashCategory(id string) error {

	if err := config.DB.Where("category_id = ?", id).Delete(&domain.TrashDetail{}).Error; err != nil {
		return err
	}

	if err := config.DB.Where("id = ?", id).Delete(&domain.TrashCategory{}).Error; err != nil {
		return err
	}

	return nil
}

func DeleteTrashDetail(id string) error {

	if err := config.DB.Where("id = ?", id).Delete(&domain.TrashDetail{}).Error; err != nil {
		return err
	}

	return nil
}
