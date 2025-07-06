package address

import (
	"context"
	"errors"
	"fmt"
	"time"

	"rijig/internal/wilayahindo"
	"rijig/model"
	"rijig/utils"
)

const (
	cacheTTL = time.Hour * 24

	userAddressesCacheKeyPattern = "user:%s:addresses"
	addressCacheKeyPattern       = "address:%s"
)

type AddressService interface {
	CreateAddress(ctx context.Context, userID string, request CreateAddressDTO) (*AddressResponseDTO, error)
	GetAddressByUserID(ctx context.Context, userID string) ([]AddressResponseDTO, error)
	GetAddressByID(ctx context.Context, userID, id string) (*AddressResponseDTO, error)
	UpdateAddress(ctx context.Context, userID, id string, addressDTO CreateAddressDTO) (*AddressResponseDTO, error)
	DeleteAddress(ctx context.Context, userID, id string) error
}

type addressService struct {
	addressRepo AddressRepository
	wilayahRepo wilayahindo.WilayahIndonesiaRepository
}

func NewAddressService(addressRepo AddressRepository, wilayahRepo wilayahindo.WilayahIndonesiaRepository) AddressService {
	return &addressService{
		addressRepo: addressRepo,
		wilayahRepo: wilayahRepo,
	}
}

func (s *addressService) validateWilayahIDs(ctx context.Context, addressDTO CreateAddressDTO) (string, string, string, string, error) {

	province, _, err := s.wilayahRepo.FindProvinceByID(ctx, addressDTO.Province, 0, 0)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid province_id: %w", err)
	}

	regency, _, err := s.wilayahRepo.FindRegencyByID(ctx, addressDTO.Regency, 0, 0)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid regency_id: %w", err)
	}

	district, _, err := s.wilayahRepo.FindDistrictByID(ctx, addressDTO.District, 0, 0)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid district_id: %w", err)
	}

	village, err := s.wilayahRepo.FindVillageByID(ctx, addressDTO.Village)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid village_id: %w", err)
	}

	return province.Name, regency.Name, district.Name, village.Name, nil
}

func (s *addressService) mapToResponseDTO(address *model.Address) *AddressResponseDTO {
	createdAt, _ := utils.FormatDateToIndonesianFormat(address.CreatedAt)
	updatedAt, _ := utils.FormatDateToIndonesianFormat(address.UpdatedAt)

	return &AddressResponseDTO{
		UserID:     address.UserID,
		ID:         address.ID,
		Province:   address.Province,
		Regency:    address.Regency,
		District:   address.District,
		Village:    address.Village,
		PostalCode: address.PostalCode,
		Detail:     address.Detail,
		Latitude:   address.Latitude,
		Longitude:  address.Longitude,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func (s *addressService) invalidateAddressCaches(userID, addressID string) {
	if addressID != "" {
		addressCacheKey := fmt.Sprintf(addressCacheKeyPattern, addressID)
		if err := utils.DeleteCache(addressCacheKey); err != nil {
			fmt.Printf("Error deleting address cache: %v\n", err)
		}
	}

	userCacheKey := fmt.Sprintf(userAddressesCacheKeyPattern, userID)
	if err := utils.DeleteCache(userCacheKey); err != nil {
		fmt.Printf("Error deleting user addresses cache: %v\n", err)
	}
}

func (s *addressService) cacheAddress(addressDTO *AddressResponseDTO) {
	cacheKey := fmt.Sprintf(addressCacheKeyPattern, addressDTO.ID)
	if err := utils.SetCache(cacheKey, addressDTO, cacheTTL); err != nil {
		fmt.Printf("Error caching address to Redis: %v\n", err)
	}
}

func (s *addressService) cacheUserAddresses(ctx context.Context, userID string) ([]AddressResponseDTO, error) {
	addresses, err := s.addressRepo.FindAddressByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}

	var addressDTOs []AddressResponseDTO
	for _, address := range addresses {
		addressDTOs = append(addressDTOs, *s.mapToResponseDTO(&address))
	}

	cacheKey := fmt.Sprintf(userAddressesCacheKeyPattern, userID)
	if err := utils.SetCache(cacheKey, addressDTOs, cacheTTL); err != nil {
		fmt.Printf("Error caching addresses to Redis: %v\n", err)
	}

	return addressDTOs, nil
}

func (s *addressService) checkAddressOwnership(ctx context.Context, userID, addressID string) (*model.Address, error) {
	address, err := s.addressRepo.FindAddressByID(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("address not found: %w", err)
	}

	if address.UserID != userID {
		return nil, errors.New("you are not authorized to access this address")
	}

	return address, nil
}

func (s *addressService) CreateAddress(ctx context.Context, userID string, addressDTO CreateAddressDTO) (*AddressResponseDTO, error) {

	provinceName, regencyName, districtName, villageName, err := s.validateWilayahIDs(ctx, addressDTO)
	if err != nil {
		return nil, err
	}

	address := model.Address{
		UserID:     userID,
		Province:   provinceName,
		Regency:    regencyName,
		District:   districtName,
		Village:    villageName,
		PostalCode: addressDTO.PostalCode,
		Detail:     addressDTO.Detail,
		Latitude:   addressDTO.Latitude,
		Longitude:  addressDTO.Longitude,
	}

	if err := s.addressRepo.CreateAddress(ctx, &address); err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	responseDTO := s.mapToResponseDTO(&address)

	s.cacheAddress(responseDTO)
	s.invalidateAddressCaches(userID, "")

	return responseDTO, nil
}

func (s *addressService) GetAddressByUserID(ctx context.Context, userID string) ([]AddressResponseDTO, error) {

	cacheKey := fmt.Sprintf(userAddressesCacheKeyPattern, userID)
	var cachedAddresses []AddressResponseDTO

	if err := utils.GetCache(cacheKey, &cachedAddresses); err == nil {
		return cachedAddresses, nil
	}

	return s.cacheUserAddresses(ctx, userID)
}

func (s *addressService) GetAddressByID(ctx context.Context, userID, id string) (*AddressResponseDTO, error) {

	address, err := s.checkAddressOwnership(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf(addressCacheKeyPattern, id)
	var cachedAddress AddressResponseDTO

	if err := utils.GetCache(cacheKey, &cachedAddress); err == nil {
		return &cachedAddress, nil
	}

	responseDTO := s.mapToResponseDTO(address)
	s.cacheAddress(responseDTO)

	return responseDTO, nil
}

func (s *addressService) UpdateAddress(ctx context.Context, userID, id string, addressDTO CreateAddressDTO) (*AddressResponseDTO, error) {

	address, err := s.checkAddressOwnership(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	provinceName, regencyName, districtName, villageName, err := s.validateWilayahIDs(ctx, addressDTO)
	if err != nil {
		return nil, err
	}

	address.Province = provinceName
	address.Regency = regencyName
	address.District = districtName
	address.Village = villageName
	address.PostalCode = addressDTO.PostalCode
	address.Detail = addressDTO.Detail
	address.Latitude = addressDTO.Latitude
	address.Longitude = addressDTO.Longitude

	if err := s.addressRepo.UpdateAddress(ctx, address); err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	responseDTO := s.mapToResponseDTO(address)

	s.cacheAddress(responseDTO)
	s.invalidateAddressCaches(userID, "")

	return responseDTO, nil
}

func (s *addressService) DeleteAddress(ctx context.Context, userID, addressID string) error {

	address, err := s.checkAddressOwnership(ctx, userID, addressID)
	if err != nil {
		return err
	}

	if err := s.addressRepo.DeleteAddress(ctx, addressID); err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	s.invalidateAddressCaches(address.UserID, addressID)

	return nil
}
