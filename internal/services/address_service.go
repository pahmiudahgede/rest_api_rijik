package services

import (
	"fmt"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type AddressService interface {
	CreateAddress(userID string, request dto.CreateAddressDTO) (*dto.AddressResponseDTO, error)
}

type addressService struct {
	AddressRepo repositories.AddressRepository
	WilayahRepo repositories.WilayahIndonesiaRepository
}

func NewAddressService(addressRepo repositories.AddressRepository, wilayahRepo repositories.WilayahIndonesiaRepository) AddressService {
	return &addressService{
		AddressRepo: addressRepo,
		WilayahRepo: wilayahRepo,
	}
}

func (s *addressService) CreateAddress(userID string, request dto.CreateAddressDTO) (*dto.AddressResponseDTO, error) {

	errors, valid := request.Validate()
	if !valid {
		return nil, fmt.Errorf("validation failed: %v", errors)
	}

	province, _, err := s.WilayahRepo.FindProvinceByID(request.Province, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid province_id")
	}

	regency, _, err := s.WilayahRepo.FindRegencyByID(request.Regency, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid regency_id")
	}

	district, _, err := s.WilayahRepo.FindDistrictByID(request.District, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid district_id")
	}

	village, err := s.WilayahRepo.FindVillageByID(request.Village)
	if err != nil {
		return nil, fmt.Errorf("invalid village_id")
	}

	newAddress := &model.Address{
		UserID:     userID,
		Province:   province.Name,
		Regency:    regency.Name,
		District:   district.Name,
		Village:    village.Name,
		PostalCode: request.PostalCode,
		Detail:     request.Detail,
		Geography:  request.Geography,
	}

	err = s.AddressRepo.CreateAddress(newAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create user address: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(newAddress.CreatedAt)

	addressResponse := &dto.AddressResponseDTO{
		UserID:     newAddress.UserID,
		ID:         newAddress.ID,
		Province:   newAddress.Province,
		Regency:    newAddress.Regency,
		District:   newAddress.District,
		Village:    newAddress.Village,
		PostalCode: newAddress.PostalCode,
		Detail:     newAddress.Detail,
		Geography:  newAddress.Geography,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	}

	return addressResponse, nil
}
