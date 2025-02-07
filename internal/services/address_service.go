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
	UpdateAddress(id string, addressDTO dto.CreateAddressDTO) (*dto.AddressResponseDTO, error)
	DeleteAddress(id string) error
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

func (s *addressService) CreateAddress(userID string, addressDTO dto.CreateAddressDTO) (*dto.AddressResponseDTO, error) {

	province, _, err := s.WilayahRepo.FindProvinceByID(addressDTO.Province, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid province_id")
	}

	regency, _, err := s.WilayahRepo.FindRegencyByID(addressDTO.Regency, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid regency_id")
	}

	district, _, err := s.WilayahRepo.FindDistrictByID(addressDTO.District, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid district_id")
	}

	village, err := s.WilayahRepo.FindVillageByID(addressDTO.Village)
	if err != nil {
		return nil, fmt.Errorf("invalid village_id")
	}

	address := model.Address{
		UserID:     userID,
		Province:   province.Name,
		Regency:    regency.Name,
		District:   district.Name,
		Village:    village.Name,
		PostalCode: addressDTO.PostalCode,
		Detail:     addressDTO.Detail,
		Geography:  addressDTO.Geography,
	}

	err = s.AddressRepo.CreateAddress(&address)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %v", err)
	}

	userCacheKey := fmt.Sprintf("user:%s:addresses", userID)
	utils.DeleteData(userCacheKey)

	createdAt, _ := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressResponseDTO := &dto.AddressResponseDTO{
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

	cacheKey := fmt.Sprintf("address:%s", address.ID)
	cacheData := map[string]interface{}{
		"data": addressResponseDTO,
	}
	err = utils.SetJSONData(cacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching new address to Redis: %v\n", err)
	}

	addresses, err := s.AddressRepo.FindAddressByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated addresses for user: %v", err)
	}

	var addressDTOs []dto.AddressResponseDTO
	for _, addr := range addresses {
		createdAt, _ := utils.FormatDateToIndonesianFormat(addr.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(addr.UpdatedAt)

		addressDTOs = append(addressDTOs, dto.AddressResponseDTO{
			UserID:     addr.UserID,
			ID:         addr.ID,
			Province:   addr.Province,
			Regency:    addr.Regency,
			District:   addr.District,
			Village:    addr.Village,
			PostalCode: addr.PostalCode,
			Detail:     addr.Detail,
			Geography:  addr.Geography,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	cacheData = map[string]interface{}{
		"data": addressDTOs,
	}
	err = utils.SetJSONData(userCacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated user addresses to Redis: %v\n", err)
	}

	return addressResponseDTO, nil
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

func (s *addressService) UpdateAddress(id string, addressDTO dto.CreateAddressDTO) (*dto.AddressResponseDTO, error) {

	province, _, err := s.WilayahRepo.FindProvinceByID(addressDTO.Province, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid province_id")
	}

	regency, _, err := s.WilayahRepo.FindRegencyByID(addressDTO.Regency, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid regency_id")
	}

	district, _, err := s.WilayahRepo.FindDistrictByID(addressDTO.District, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("invalid district_id")
	}

	village, err := s.WilayahRepo.FindVillageByID(addressDTO.Village)
	if err != nil {
		return nil, fmt.Errorf("invalid village_id")
	}

	address, err := s.AddressRepo.FindAddressByID(id)
	if err != nil {
		return nil, fmt.Errorf("address not found: %v", err)
	}

	address.Province = province.Name
	address.Regency = regency.Name
	address.District = district.Name
	address.Village = village.Name
	address.PostalCode = addressDTO.PostalCode
	address.Detail = addressDTO.Detail
	address.Geography = addressDTO.Geography
	address.UpdatedAt = time.Now()

	err = s.AddressRepo.UpdateAddress(address)
	if err != nil {
		return nil, fmt.Errorf("failed to update address: %v", err)
	}

	addressCacheKey := fmt.Sprintf("address:%s", id)
	utils.DeleteData(addressCacheKey)

	userAddressesCacheKey := fmt.Sprintf("user:%s:addresses", address.UserID)
	utils.DeleteData(userAddressesCacheKey)

	createdAt, _ := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	addressResponseDTO := &dto.AddressResponseDTO{
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
		"data": addressResponseDTO,
	}
	err = utils.SetJSONData(addressCacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated address to Redis: %v\n", err)
	}

	addresses, err := s.AddressRepo.FindAddressByUserID(address.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated addresses for user: %v", err)
	}

	var addressDTOs []dto.AddressResponseDTO
	for _, addr := range addresses {
		createdAt, _ := utils.FormatDateToIndonesianFormat(addr.CreatedAt)
		updatedAt, _ := utils.FormatDateToIndonesianFormat(addr.UpdatedAt)

		addressDTOs = append(addressDTOs, dto.AddressResponseDTO{
			UserID:     addr.UserID,
			ID:         addr.ID,
			Province:   addr.Province,
			Regency:    addr.Regency,
			District:   addr.District,
			Village:    addr.Village,
			PostalCode: addr.PostalCode,
			Detail:     addr.Detail,
			Geography:  addr.Geography,
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
		})
	}

	cacheData = map[string]interface{}{
		"data": addressDTOs,
	}
	err = utils.SetJSONData(userAddressesCacheKey, cacheData, time.Hour*24)
	if err != nil {
		fmt.Printf("Error caching updated user addresses to Redis: %v\n", err)
	}

	return addressResponseDTO, nil
}

func (s *addressService) DeleteAddress(id string) error {

	address, err := s.AddressRepo.FindAddressByID(id)
	if err != nil {
		return fmt.Errorf("address not found: %v", err)
	}

	err = s.AddressRepo.DeleteAddress(id)
	if err != nil {
		return fmt.Errorf("failed to delete address: %v", err)
	}

	addressCacheKey := fmt.Sprintf("address:%s", id)
	err = utils.DeleteData(addressCacheKey)
	if err != nil {
		fmt.Printf("Error deleting address cache: %v\n", err)
	}

	userAddressesCacheKey := fmt.Sprintf("user:%s:addresses", address.UserID)
	err = utils.DeleteData(userAddressesCacheKey)
	if err != nil {
		fmt.Printf("Error deleting user addresses cache: %v\n", err)
	}

	return nil
}
