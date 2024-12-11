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

func UpdateTrashCategory(id, name string) (domain.TrashCategory, error) {
	category, err := repositories.GetTrashCategoryDetail(id)
	if err != nil {
		return domain.TrashCategory{}, err
	}
	category.Name = name
	if err := repositories.UpdateTrashCategory(&category); err != nil {
		return domain.TrashCategory{}, err
	}
	return category, nil
}

func UpdateTrashDetail(id, description string, price int) (domain.TrashDetail, error) {

	detail, err := repositories.GetTrashDetailByID(id)
	if err != nil {

		return domain.TrashDetail{}, err
	}

	detail.Description = description
	detail.Price = price

	if err := repositories.UpdateTrashDetail(&detail); err != nil {

		return domain.TrashDetail{}, err
	}

	return detail, nil
}
