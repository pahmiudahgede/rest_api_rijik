package repositories

import (
	"errors"
	"fmt"
	"log"
	"rijig/model"
	"rijig/utils"

	"gorm.io/gorm"
)

type CollectorRepository interface {
	FindActiveCollectors() ([]model.Collector, error)
	FindCollectorById(collector_id string) (*model.Collector, error)
	CreateCollector(collector *model.Collector) error
	UpdateCollector(userId string, jobStatus string) (*model.Collector, error)
	FindAllAutomaticMethodRequestWithDistance(requestMethod, statuspickup string, collectorLat, collectorLon float64, maxDistance float64) ([]model.RequestPickup, error)
}

type collectorRepository struct {
	DB *gorm.DB
}

func NewCollectorRepository(db *gorm.DB) CollectorRepository {
	return &collectorRepository{DB: db}
}

func (r *collectorRepository) FindActiveCollectors() ([]model.Collector, error) {
	var collectors []model.Collector

	err := r.DB.Where("job_status = ?", "active").First(&collectors).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active collectors: %v", err)
	}

	return collectors, nil
}

func (r *collectorRepository) FindCollectorById(collector_id string) (*model.Collector, error) {
	var collector model.Collector
	err := r.DB.Where("user_id = ?", collector_id).First(&collector).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching collector: %v", err)
	}
	fmt.Printf("menampilkan data collector %v", &collector)
	return &collector, nil
}

func (r *collectorRepository) CreateCollector(collector *model.Collector) error {
	if err := r.DB.Create(collector).Error; err != nil {
		return fmt.Errorf("failed to create collector: %v", err)
	}
	return nil
}

func (r *collectorRepository) UpdateCollector(userId string, jobStatus string) (*model.Collector, error) {
	var existingCollector model.Collector

	if err := r.DB.Where("user_id = ?", userId).First(&existingCollector).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("collector dengan user_id %s tidak ditemukan", userId)
		}
		log.Printf("Gagal mencari collector: %v", err)
		return nil, fmt.Errorf("gagal fetching collector: %w", err)
	}

	if jobStatus != "active" && jobStatus != "nonactive" {
		return nil, fmt.Errorf("invalid job status: %v", jobStatus)
	}

	if err := r.DB.Model(&existingCollector).Update("jobstatus", jobStatus).Error; err != nil {
		log.Printf("Gagal mengupdate data collector: %v", err)
		return nil, fmt.Errorf("gagal mengupdate job status untuk collector: %w", err)
	}

	return &existingCollector, nil
}

// #====experimen====#
func (r *collectorRepository) FindAllAutomaticMethodRequestWithDistance(requestMethod, statuspickup string, collectorLat, collectorLon float64, maxDistance float64) ([]model.RequestPickup, error) {
	var requests []model.RequestPickup

	err := r.DB.Preload("RequestItems").
		Where("request_method = ? AND status_pickup = ?", requestMethod, statuspickup).
		Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching request pickups with request_method '%s' and status '%s': %v", requestMethod, statuspickup, err)
	}

	var nearbyRequests []model.RequestPickup
	for _, request := range requests {
		address := request.Address

		requestCoord := utils.Coord{Lat: address.Latitude, Lon: address.Longitude}
		collectorCoord := utils.Coord{Lat: collectorLat, Lon: collectorLon}
		_, km := utils.Distance(requestCoord, collectorCoord)

		if km <= maxDistance {
			nearbyRequests = append(nearbyRequests, request)
		}
	}

	return nearbyRequests, nil
}
