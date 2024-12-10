package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetTrashCategories() ([]domain.TrashCategory, error) {

	return repositories.GetTrashCategories()
}

func GetTrashCategoryDetail(id string) (domain.TrashCategory, error) {
	return repositories.GetTrashCategoryDetail(id)
}

func CreateTrashCategory(name string) (domain.TrashCategory, error) {
	category := domain.TrashCategory{Name: name}

	err := repositories.CreateTrashCategory(&category)
	if err != nil {
		return category, err
	}

	return category, nil
}

func CreateTrashDetail(categoryID, description string, price int) (domain.TrashDetail, error) {
	detail := domain.TrashDetail{
		CategoryID:  categoryID,
		Description: description,
		Price:       price,
	}

	err := repositories.CreateTrashDetail(&detail)
	if err != nil {
		return detail, err
	}

	return detail, nil
}
