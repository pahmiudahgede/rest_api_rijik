package services

import (
	"fmt"
	"time"

	"github.com/pahmiudahgede/senggoldong/dto"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/model"
	"github.com/pahmiudahgede/senggoldong/utils"
)

type AddressService interface {
	CreateAddress(userID string, request dto.CreateAddressDTO) (*dto.AddressResponseDTO, error)
	GetAddressByUserID(userID string) ([]dto.AddressResponseDTO, error)
	GetAddressByID(id string) (*dto.AddressResponseDTO, error)
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

func (s *addressService) GetAddressByUserID(userID string) ([]dto.AddressResponseDTO, error) {

	cacheKey := fmt.Sprintf("user:%s:addresses", userID)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		var addresses []dto.AddressResponseDTO
		if data, ok := cachedData["data"].([]interface{}); ok {
			for _, item := range data {
				addressData, ok := item.(map[string]interface{})
				if ok {
					addresses = append(addresses, dto.AddressResponseDTO{
						UserID:     addressData["user_id"].(string),
						ID:         addressData["address_id"].(string),
						Province:   addressData["province"].(string),
						Regency:    addressData["regency"].(string),
						District:   addressData["district"].(string),
						Village:    addressData["village"].(string),
						PostalCode: addressData["postalCode"].(string),
						Detail:     addressData["detail"].(string),
						Geography:  addressData["geography"].(string),
						CreatedAt:  addressData["createdAt"].(string),
						UpdatedAt:  addressData["updatedAt"].(string),
					})
				}
			}
			return addresses, nil
		}
	}

	addresses, err := s.AddressRepo.FindAddressByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %v", err)
	}

	var addressDTOs []dto.AddressResponseDTO
	for _, address := range addresses {
		createdAt, _ := utils.FormatDateToIndonesianFormat(address.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

		addressDTOs = append(addressDTOs, dto.AddressResponseDTO{
			UserID:     address.UserID,
			ID:         address.ID,
			Province:   address.Province,
			Regency:    address.Regency,
			District:   address.District,
			Village:    address.Village,
			PostalCode: address.PostalCode,
			Detail:     address.Detail,
			Geography:  address.Geography,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	cacheData := map[string]interface{}{
		"data": addressDTOs,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching addresses to Redis: %v\n", err)
	}

	return addressDTOs, nil
}

func (s *addressService) GetAddressByID(id string) (*dto.AddressResponseDTO, error) {

	cacheKey := fmt.Sprintf("address:%s", id)
	cachedData, err := utils.GetJSONData(cacheKey)
	if err == nil && cachedData != nil {
		addressData, ok := cachedData["data"].(map[string]interface{})
		if ok {
			address := dto.AddressResponseDTO{
				UserID:     addressData["user_id"].(string),
				ID:         addressData["address_id"].(string),
				Province:   addressData["province"].(string),
				Regency:    addressData["regency"].(string),
				District:   addressData["district"].(string),
				Village:    addressData["village"].(string),
				PostalCode: addressData["postalCode"].(string),
				Detail:     addressData["detail"].(string),
				Geography:  addressData["geography"].(string),
				CreatedAt:  addressData["createdAt"].(string),
				UpdatedAt:  addressData["updatedAt"].(string),
			}
			return &address, nil
		}
	}

	address, err := s.AddressRepo.FindAddressByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address: %v", err)
	}

	createdAt, _ := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressDTO := &dto.AddressResponseDTO{
		UserID:     address.UserID,
		ID:         address.ID,
		Province:   address.Province,
		Regency:    address.Regency,
		District:   address.District,
		Village:    address.Village,
		PostalCode: address.PostalCode,
		Detail:     address.Detail,
		Geography:  address.Geography,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}

	cacheData := map[string]interface{}{
		"data": addressDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching address to Redis: %v\n", err)
	}

	return addressDTO, nil
}
