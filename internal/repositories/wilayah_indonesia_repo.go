package repositories

import (
	"rijig/model"

	"gorm.io/gorm"
)

type WilayahIndonesiaRepository interface {
	ImportProvinces(provinces []model.Province) error
	ImportRegencies(regencies []model.Regency) error
	ImportDistricts(districts []model.District) error
	ImportVillages(villages []model.Village) error

	FindAllProvinces(page, limit int) ([]model.Province, int, error)
	FindProvinceByID(id string, page, limit int) (*model.Province, int, error)

	FindAllRegencies(page, limit int) ([]model.Regency, int, error)
	FindRegencyByID(id string, page, limit int) (*model.Regency, int, error)

	FindAllDistricts(page, limit int) ([]model.District, int, error)
	FindDistrictByID(id string, page, limit int) (*model.District, int, error)

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

func (r *wilayahIndonesiaRepository) FindAllProvinces(page, limit int) ([]model.Province, int, error) {
	var provinces []model.Province
	var total int64

	err := r.DB.Model(&model.Province{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&provinces).Error
		if err != nil {
			return nil, 0, err
		}
	} else {

		err := r.DB.Find(&provinces).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return provinces, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindProvinceByID(id string, page, limit int) (*model.Province, int, error) {
	var province model.Province

	err := r.DB.Preload("Regencies", func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {

			return db.Offset((page - 1) * limit).Limit(limit)
		}

		return db
	}).Where("id = ?", id).First(&province).Error
	if err != nil {
		return nil, 0, err
	}

	var totalRegencies int64
	r.DB.Model(&model.Regency{}).Where("province_id = ?", id).Count(&totalRegencies)

	return &province, int(totalRegencies), nil
}

func (r *wilayahIndonesiaRepository) FindAllRegencies(page, limit int) ([]model.Regency, int, error) {
	var regencies []model.Regency
	var total int64

	err := r.DB.Model(&model.Regency{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&regencies).Error
		if err != nil {
			return nil, 0, err
		}
	} else {

		err := r.DB.Find(&regencies).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return regencies, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindRegencyByID(id string, page, limit int) (*model.Regency, int, error) {
	var regency model.Regency

	err := r.DB.Preload("Districts", func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {
			return db.Offset((page - 1) * limit).Limit(limit)
		}
		return db
	}).Where("id = ?", id).First(&regency).Error

	if err != nil {
		return nil, 0, err
	}

	var totalDistricts int64
	err = r.DB.Model(&model.District{}).Where("regency_id = ?", id).Count(&totalDistricts).Error
	if err != nil {
		return nil, 0, err
	}

	return &regency, int(totalDistricts), nil
}

func (r *wilayahIndonesiaRepository) FindAllDistricts(page, limit int) ([]model.District, int, error) {
	var district []model.District
	var total int64

	err := r.DB.Model(&model.District{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&district).Error
		if err != nil {
			return nil, 0, err
		}
	} else {

		err := r.DB.Find(&district).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return district, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindDistrictByID(id string, page, limit int) (*model.District, int, error) {
	var district model.District

	err := r.DB.Preload("Villages", func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {

			return db.Offset((page - 1) * limit).Limit(limit)
		}

		return db
	}).Where("id = ?", id).First(&district).Error
	if err != nil {
		return nil, 0, err
	}

	var totalVillage int64
	r.DB.Model(&model.Village{}).Where("district_id = ?", id).Count(&totalVillage)

	return &district, int(totalVillage), nil
}

func (r *wilayahIndonesiaRepository) FindAllVillages(page, limit int) ([]model.Village, int, error) {
	var villages []model.Village
	var total int64

	err := r.DB.Model(&model.Village{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		err := r.DB.Offset((page - 1) * limit).Limit(limit).Find(&villages).Error
		if err != nil {
			return nil, 0, err
		}
	} else {

		err := r.DB.Find(&villages).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return villages, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindVillageByID(id string) (*model.Village, error) {
	var village model.Village
	err := r.DB.Where("id = ?", id).First(&village).Error
	if err != nil {
		return nil, err
	}
	return &village, nil
}
