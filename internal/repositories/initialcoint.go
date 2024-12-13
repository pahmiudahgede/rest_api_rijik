package repositories

import (
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

func GetPoints() ([]domain.Point, error) {
	var points []domain.Point
	if err := config.DB.Find(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}

func GetPointByID(id string) (domain.Point, error) {
	var point domain.Point
	if err := config.DB.Where("id = ?", id).First(&point).Error; err != nil {
		return point, err
	}
	return point, nil
}

func CreatePoint(point *domain.Point) error {

	if err := config.DB.Create(point).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePoint(point *domain.Point) error {
	if err := config.DB.Save(point).Error; err != nil {
		return err
	}
	return nil
}

func DeletePoint(id string) error {
	if err := config.DB.Where("id = ?", id).Delete(&domain.Point{}).Error; err != nil {
		return err
	}
	return nil
}