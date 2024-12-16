package services

import (
	"github.com/pahmiudahgede/senggoldong/domain"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
)

func GetRequestPickupsByUser(userID string) ([]domain.RequestPickup, error) {

	return repositories.GetRequestPickupsByUser(userID)
}
