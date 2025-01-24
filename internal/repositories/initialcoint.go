package repositories

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/domain"
)

type PointRepository struct{}

func NewPointRepository() *PointRepository {
	return &PointRepository{}
}

func (r *PointRepository) GetAll() ([]domain.Point, error) {
	var points []domain.Point
	err := config.DB.Find(&points).Error
	if err != nil {
		return nil, errors.New("failed to fetch points from database")
	}
	return points, nil
}

func (r *PointRepository) GetByID(id string) (*domain.Point, error) {
	var point domain.Point
	err := config.DB.First(&point, "id = ?", id).Error
	if err != nil {
		return nil, errors.New("point not found")
	}
	return &point, nil
}

func (r *PointRepository) Create(point *domain.Point) error {
	return config.DB.Create(point).Error
}

func (r *PointRepository) Update(point *domain.Point) error {
	return config.DB.Save(point).Error
}

func (r *PointRepository) Delete(point *domain.Point) error {
	return config.DB.Delete(point).Error
}