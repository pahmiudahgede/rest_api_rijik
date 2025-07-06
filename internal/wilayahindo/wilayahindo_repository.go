package wilayahindo

import (
	"context"
	"errors"
	"fmt"
	"rijig/model"

	"gorm.io/gorm"
)

type WilayahIndonesiaRepository interface {
	ImportProvinces(ctx context.Context, provinces []model.Province) error
	ImportRegencies(ctx context.Context, regencies []model.Regency) error
	ImportDistricts(ctx context.Context, districts []model.District) error
	ImportVillages(ctx context.Context, villages []model.Village) error

	FindAllProvinces(ctx context.Context, page, limit int) ([]model.Province, int, error)
	FindProvinceByID(ctx context.Context, id string, page, limit int) (*model.Province, int, error)

	FindAllRegencies(ctx context.Context, page, limit int) ([]model.Regency, int, error)
	FindRegencyByID(ctx context.Context, id string, page, limit int) (*model.Regency, int, error)

	FindAllDistricts(ctx context.Context, page, limit int) ([]model.District, int, error)
	FindDistrictByID(ctx context.Context, id string, page, limit int) (*model.District, int, error)

	FindAllVillages(ctx context.Context, page, limit int) ([]model.Village, int, error)
	FindVillageByID(ctx context.Context, id string) (*model.Village, error)
}

type wilayahIndonesiaRepository struct {
	DB *gorm.DB
}

func NewWilayahIndonesiaRepository(db *gorm.DB) WilayahIndonesiaRepository {
	return &wilayahIndonesiaRepository{DB: db}
}

func (r *wilayahIndonesiaRepository) ImportProvinces(ctx context.Context, provinces []model.Province) error {
	if len(provinces) == 0 {
		return errors.New("no provinces to import")
	}

	if err := r.DB.WithContext(ctx).CreateInBatches(provinces, 100).Error; err != nil {
		return fmt.Errorf("failed to import provinces: %w", err)
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportRegencies(ctx context.Context, regencies []model.Regency) error {
	if len(regencies) == 0 {
		return errors.New("no regencies to import")
	}

	if err := r.DB.WithContext(ctx).CreateInBatches(regencies, 100).Error; err != nil {
		return fmt.Errorf("failed to import regencies: %w", err)
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportDistricts(ctx context.Context, districts []model.District) error {
	if len(districts) == 0 {
		return errors.New("no districts to import")
	}

	if err := r.DB.WithContext(ctx).CreateInBatches(districts, 100).Error; err != nil {
		return fmt.Errorf("failed to import districts: %w", err)
	}
	return nil
}

func (r *wilayahIndonesiaRepository) ImportVillages(ctx context.Context, villages []model.Village) error {
	if len(villages) == 0 {
		return errors.New("no villages to import")
	}

	if err := r.DB.WithContext(ctx).CreateInBatches(villages, 100).Error; err != nil {
		return fmt.Errorf("failed to import villages: %w", err)
	}
	return nil
}

func (r *wilayahIndonesiaRepository) FindAllProvinces(ctx context.Context, page, limit int) ([]model.Province, int, error) {
	var provinces []model.Province
	var total int64

	if err := r.DB.WithContext(ctx).Model(&model.Province{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count provinces: %w", err)
	}

	query := r.DB.WithContext(ctx)

	if page > 0 && limit > 0 {
		if page < 1 {
			return nil, 0, errors.New("page must be greater than 0")
		}
		if limit < 1 || limit > 1000 {
			return nil, 0, errors.New("limit must be between 1 and 1000")
		}
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	if err := query.Find(&provinces).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find provinces: %w", err)
	}

	return provinces, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindProvinceByID(ctx context.Context, id string, page, limit int) (*model.Province, int, error) {
	if id == "" {
		return nil, 0, errors.New("province ID cannot be empty")
	}

	var province model.Province

	preloadQuery := func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {
			if page < 1 {
				return db
			}
			if limit < 1 || limit > 1000 {
				return db
			}
			return db.Offset((page - 1) * limit).Limit(limit)
		}
		return db
	}

	if err := r.DB.WithContext(ctx).Preload("Regencies", preloadQuery).Where("id = ?", id).First(&province).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("province with ID %s not found", id)
		}
		return nil, 0, fmt.Errorf("failed to find province: %w", err)
	}

	var totalRegencies int64
	if err := r.DB.WithContext(ctx).Model(&model.Regency{}).Where("province_id = ?", id).Count(&totalRegencies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count regencies: %w", err)
	}

	return &province, int(totalRegencies), nil
}

func (r *wilayahIndonesiaRepository) FindAllRegencies(ctx context.Context, page, limit int) ([]model.Regency, int, error) {
	var regencies []model.Regency
	var total int64

	if err := r.DB.WithContext(ctx).Model(&model.Regency{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count regencies: %w", err)
	}

	query := r.DB.WithContext(ctx)

	if page > 0 && limit > 0 {
		if page < 1 {
			return nil, 0, errors.New("page must be greater than 0")
		}
		if limit < 1 || limit > 1000 {
			return nil, 0, errors.New("limit must be between 1 and 1000")
		}
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	if err := query.Find(&regencies).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find regencies: %w", err)
	}

	return regencies, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindRegencyByID(ctx context.Context, id string, page, limit int) (*model.Regency, int, error) {
	if id == "" {
		return nil, 0, errors.New("regency ID cannot be empty")
	}

	var regency model.Regency

	preloadQuery := func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {
			if page < 1 {
				return db
			}
			if limit < 1 || limit > 1000 {
				return db
			}
			return db.Offset((page - 1) * limit).Limit(limit)
		}
		return db
	}

	if err := r.DB.WithContext(ctx).Preload("Districts", preloadQuery).Where("id = ?", id).First(&regency).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("regency with ID %s not found", id)
		}
		return nil, 0, fmt.Errorf("failed to find regency: %w", err)
	}

	var totalDistricts int64
	if err := r.DB.WithContext(ctx).Model(&model.District{}).Where("regency_id = ?", id).Count(&totalDistricts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count districts: %w", err)
	}

	return &regency, int(totalDistricts), nil
}

func (r *wilayahIndonesiaRepository) FindAllDistricts(ctx context.Context, page, limit int) ([]model.District, int, error) {
	var districts []model.District
	var total int64

	if err := r.DB.WithContext(ctx).Model(&model.District{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count districts: %w", err)
	}

	query := r.DB.WithContext(ctx)

	if page > 0 && limit > 0 {
		if page < 1 {
			return nil, 0, errors.New("page must be greater than 0")
		}
		if limit < 1 || limit > 1000 {
			return nil, 0, errors.New("limit must be between 1 and 1000")
		}
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	if err := query.Find(&districts).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find districts: %w", err)
	}

	return districts, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindDistrictByID(ctx context.Context, id string, page, limit int) (*model.District, int, error) {
	if id == "" {
		return nil, 0, errors.New("district ID cannot be empty")
	}

	var district model.District

	preloadQuery := func(db *gorm.DB) *gorm.DB {
		if page > 0 && limit > 0 {
			if page < 1 {
				return db
			}
			if limit < 1 || limit > 1000 {
				return db
			}
			return db.Offset((page - 1) * limit).Limit(limit)
		}
		return db
	}

	if err := r.DB.WithContext(ctx).Preload("Villages", preloadQuery).Where("id = ?", id).First(&district).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, fmt.Errorf("district with ID %s not found", id)
		}
		return nil, 0, fmt.Errorf("failed to find district: %w", err)
	}

	var totalVillages int64
	if err := r.DB.WithContext(ctx).Model(&model.Village{}).Where("district_id = ?", id).Count(&totalVillages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count villages: %w", err)
	}

	return &district, int(totalVillages), nil
}

func (r *wilayahIndonesiaRepository) FindAllVillages(ctx context.Context, page, limit int) ([]model.Village, int, error) {
	var villages []model.Village
	var total int64

	if err := r.DB.WithContext(ctx).Model(&model.Village{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count villages: %w", err)
	}

	query := r.DB.WithContext(ctx)

	if page > 0 && limit > 0 {
		if page < 1 {
			return nil, 0, errors.New("page must be greater than 0")
		}
		if limit < 1 || limit > 1000 {
			return nil, 0, errors.New("limit must be between 1 and 1000")
		}
		query = query.Offset((page - 1) * limit).Limit(limit)
	}

	if err := query.Find(&villages).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find villages: %w", err)
	}

	return villages, int(total), nil
}

func (r *wilayahIndonesiaRepository) FindVillageByID(ctx context.Context, id string) (*model.Village, error) {
	if id == "" {
		return nil, errors.New("village ID cannot be empty")
	}

	var village model.Village
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&village).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("village with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to find village: %w", err)
	}

	return &village, nil
}
