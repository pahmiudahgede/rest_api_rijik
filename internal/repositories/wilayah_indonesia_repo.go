package repositories

import (
	"github.com/pahmiudahgede/senggoldong/model"
	"gorm.io/gorm"
)

type WilayahIndonesiaRepository interface {
	ImportProvinces(provinces []model.Province) error
	ImportRegencies(regencies []model.Regency) error
	ImportDistricts(districts []model.District) error
	ImportVillages(villages []model.Village) error
	// ================================================== //
	FindAllProvinces(page, limit int) ([]model.Province, int, error)
	FindProvinceByID(id string) (*model.Province, error)

	FindAllRegencies(page, limit int) ([]model.Regency, int, error)
	FindRegencyByID(id string) (*model.Regency, error)

	FindAllDistricts(page, limit int) ([]model.District, int, error)
	FindDistrictByID(id string) (*model.District, error)

	FindAllVillages(page, limit int) ([]model.Village, int, error)
	FindVillageByID(id string) (*model.Village, error)
}

type wilayahIndonesiaRepository struct {
	DB *gorm.DB
}

func NewWilayahIndonesiaRepository(db *gorm.DB) WilayahIndonesiaRepository {
	return &wilayahIndonesiaRepository{DB: db}
}

func (r *wilayahIndonesiaRepository) ImportProvinces(provinces []model.Province) error {
	for _, province := range provinces {
		if err := r.DB.Create(&province).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportRegencies(regencies []model.Regency) error {
	for _, regency := range regencies {
		if err := r.DB.Create(&regency).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportDistricts(districts []model.District) error {
	for _, district := range districts {
		if err := r.DB.Create(&district).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportVillages(villages []model.Village) error {
	for _, village := range villages {
		if err := r.DB.Create(&village).Error; err != nil {
			return err
		}
	}
	return nil
}

/*
|	============================================================	|
|	============================================================	|
*/

// FindAllProvinces with Pagination
func (r *wilayahIndonesiaRepository) FindAllProvinces(page, limit int) ([]model.Province, int, error) {
	var provinces []model.Province
	var total int64

	// Count total provinces
	err := r.DB.Model(&model.Province{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get provinces with pagination
	err = r.DB.Offset((page - 1) * limit).Limit(limit).Find(&provinces).Error
	if err != nil {
		return nil, 0, err
	}

	return provinces, int(total), nil
}

// FindProvinceByID
func (r *wilayahIndonesiaRepository) FindProvinceByID(id string) (*model.Province, error) {
	var province model.Province
	err := r.DB.Preload("Regencies").Where("id = ?", id).First(&province).Error
	if err != nil {
		return nil, err
	}
	return &province, nil
}

// FindAllRegencies with Pagination
func (r *wilayahIndonesiaRepository) FindAllRegencies(page, limit int) ([]model.Regency, int, error) {
	var regencies []model.Regency
	var total int64

	// Count total regencies
	err := r.DB.Model(&model.Regency{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get regencies with pagination
	err = r.DB.Offset((page - 1) * limit).Limit(limit).Preload("Districts").Find(&regencies).Error
	if err != nil {
		return nil, 0, err
	}

	return regencies, int(total), nil
}

// FindRegencyByID
func (r *wilayahIndonesiaRepository) FindRegencyByID(id string) (*model.Regency, error) {
	var regency model.Regency
	err := r.DB.Preload("Districts").Where("id = ?", id).First(&regency).Error
	if err != nil {
		return nil, err
	}
	return &regency, nil
}

// FindAllDistricts with Pagination
func (r *wilayahIndonesiaRepository) FindAllDistricts(page, limit int) ([]model.District, int, error) {
	var districts []model.District
	var total int64

	// Count total districts
	err := r.DB.Model(&model.District{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get districts with pagination
	err = r.DB.Offset((page - 1) * limit).Limit(limit).Preload("Villages").Find(&districts).Error
	if err != nil {
		return nil, 0, err
	}

	return districts, int(total), nil
}

// FindDistrictByID
func (r *wilayahIndonesiaRepository) FindDistrictByID(id string) (*model.District, error) {
	var district model.District
	err := r.DB.Preload("Villages").Where("id = ?", id).First(&district).Error
	if err != nil {
		return nil, err
	}
	return &district, nil
}

// FindAllVillages with Pagination
func (r *wilayahIndonesiaRepository) FindAllVillages(page, limit int) ([]model.Village, int, error) {
	var villages []model.Village
	var total int64

	// Count total villages
	err := r.DB.Model(&model.Village{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get villages with pagination
	err = r.DB.Offset((page - 1) * limit).Limit(limit).Find(&villages).Error
	if err != nil {
		return nil, 0, err
	}

	return villages, int(total), nil
}

// FindVillageByID
func (r *wilayahIndonesiaRepository) FindVillageByID(id string) (*model.Village, error) {
	var village model.Village
	err := r.DB.Where("id = ?", id).First(&village).Error
	if err != nil {
		return nil, err
	}
	return &village, nil
}