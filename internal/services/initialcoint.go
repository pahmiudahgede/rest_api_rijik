package services

import (
	"errors"

	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetPoints() ([]domain.Point, error) {
	return repositories.GetPoints()
}

func GetPointByID(id string) (domain.Point, error) {
	point, err := repositories.GetPointByID(id)
	if err != nil {
		return domain.Point{}, errors.New("point not found")
	}
	return point, nil
}

func CreatePoint(coinName string, valuePerUnit float64) (domain.Point, error) {

	newPoint := domain.Point{
		CoinName:     coinName,
		ValuePerUnit: valuePerUnit,
	}

	if err := repositories.CreatePoint(&newPoint); err != nil {
		return domain.Point{}, err
	}

	return newPoint, nil
}

func UpdatePoint(id, coinName string, valuePerUnit float64) (domain.Point, error) {

	point, err := repositories.GetPointByID(id)
	if err != nil {
		return domain.Point{}, errors.New("point not found")
	}

	point.CoinName = coinName
	point.ValuePerUnit = valuePerUnit

	if err := repositories.UpdatePoint(&point); err != nil {
		return domain.Point{}, err
	}

	return point, nil
}

func DeletePoint(id string) error {

	_, err := repositories.GetPointByID(id)
	if err != nil {
		return errors.New("point not found")
	}

	if err := repositories.DeletePoint(id); err != nil {
		return err
	}

	return nil
}
