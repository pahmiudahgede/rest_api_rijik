package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"
	"rijig/utils"
)

type CartService interface {
	CreateCartFromDTO(userID string, items []dto.RequestCartItems) error
	GetCartByUserID(userID string) (*dto.CartResponse, error)
	CommitCartFromRedis(userID string) error
	DeleteCart(cartID string) error
}

type cartService struct {
	repo      repositories.CartRepository
	repoTrash repositories.TrashRepository
}

func NewCartService(repo repositories.CartRepository, repoTrash repositories.TrashRepository) CartService {
	return &cartService{repo: repo, repoTrash: repoTrash}
}

func redisCartKey(userID string) string {
	return fmt.Sprintf("cart:user:%s", userID)
}

func (s *cartService) CreateCartFromDTO(userID string, items []dto.RequestCartItems) error {
	// Validasi semua item
	for _, item := range items {
		if errMap, valid := item.ValidateRequestCartItem(); !valid {
			return dto.ValidationErrors{Errors: errMap}
		}
	}

	// Ambil cart yang sudah ada dari Redis (jika ada)
	var existingCart dto.CartResponse
	val, err := utils.GetData(redisCartKey(userID))
	if err == nil && val != "" {
		if err := json.Unmarshal([]byte(val), &existingCart); err != nil {
			log.Printf("Failed to unmarshal existing cart: %v", err)
		}
	}

	// Buat map dari existing items untuk mempermudah update
	itemMap := make(map[string]dto.CartItemResponse)
	for _, item := range existingCart.CartItems {
		itemMap[item.TrashName] = item
	}

	// Proses input baru
	for _, input := range items {
		trash, err := s.repoTrash.GetCategoryByID(input.TrashID)
		if err != nil {
			return fmt.Errorf("failed to retrieve trash category for id %s: %v", input.TrashID, err)
		}

		if input.Amount == 0 {
			delete(itemMap, trash.Name) // hapus item
			continue
		}

		subtotal := float32(trash.EstimatedPrice) * input.Amount

		itemMap[trash.Name] = dto.CartItemResponse{
			TrashIcon:              trash.Icon,
			TrashName:              trash.Name,
			Amount:                 input.Amount,
			EstimatedSubTotalPrice: subtotal,
		}
	}

	// Rekonstruksi cart
	var finalItems []dto.CartItemResponse
	var totalAmount float32
	var totalPrice float32
	for _, item := range itemMap {
		finalItems = append(finalItems, item)
		totalAmount += item.Amount
		totalPrice += item.EstimatedSubTotalPrice
	}

	cart := dto.CartResponse{
		ID:                  existingCart.ID,
		UserID:              userID,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: totalPrice,
		CartItems:           finalItems,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Simpan ulang ke Redis dengan TTL 10 menit
	return utils.SetData(redisCartKey(userID), cart, 1*time.Minute)
}


func (s *cartService) GetCartByUserID(userID string) (*dto.CartResponse, error) {
	val, err := utils.GetData(redisCartKey(userID))
	if err != nil {
		log.Printf("Redis get error: %v", err)
	}
	if val != "" {
		var cached dto.CartResponse
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			return &cached, nil
		}
	}

	cart, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, nil
	}

	var items []dto.CartItemResponse
	for _, item := range cart.CartItems {
		items = append(items, dto.CartItemResponse{
			TrashIcon:              item.TrashCategory.Icon,
			TrashName:              item.TrashCategory.Name,
			Amount:                 item.Amount,
			EstimatedSubTotalPrice: item.SubTotalEstimatedPrice,
		})
	}

	response := &dto.CartResponse{
		ID:                  cart.ID,
		UserID:              cart.UserID,
		TotalAmount:         cart.TotalAmount,
		EstimatedTotalPrice: cart.EstimatedTotalPrice,
		CartItems:           items,
		CreatedAt:           cart.CreatedAt,
		UpdatedAt:           cart.UpdatedAt,
	}

	return response, nil
}

func (s *cartService) CommitCartFromRedis(userID string) error {
	val, err := utils.GetData(redisCartKey(userID))
	if err != nil || val == "" {
		return errors.New("no cart found in redis")
	}

	var cartDTO dto.CartResponse
	if err := json.Unmarshal([]byte(val), &cartDTO); err != nil {
		return errors.New("invalid cart data in Redis")
	}

	existingCart, err := s.repo.GetByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to get cart from db: %v", err)
	}

	if existingCart == nil {
		// buat cart baru jika belum ada
		var items []model.CartItem
		for _, item := range cartDTO.CartItems {
			trash, err := s.repoTrash.GetTrashCategoryByName(item.TrashName)
			if err != nil {
				continue
			}

			items = append(items, model.CartItem{
				TrashID:                trash.ID,
				Amount:                 item.Amount,
				SubTotalEstimatedPrice: item.EstimatedSubTotalPrice,
			})
		}

		newCart := model.Cart{
			UserID:              userID,
			TotalAmount:         cartDTO.TotalAmount,
			EstimatedTotalPrice: cartDTO.EstimatedTotalPrice,
			CartItems:           items,
		}

		return s.repo.Create(&newCart)
	}

	// buat map item lama (by trash_name)
	existingItemMap := make(map[string]*model.CartItem)
	for i := range existingCart.CartItems {
		trashName := existingCart.CartItems[i].TrashCategory.Name
		existingItemMap[trashName] = &existingCart.CartItems[i]
	}

	// proses update/hapus/tambah
	for _, newItem := range cartDTO.CartItems {
		if newItem.Amount == 0 {
			if existing, ok := existingItemMap[newItem.TrashName]; ok {
				_ = s.repo.DeleteCartItemByID(existing.ID)
			}
			continue
		}

		trash, err := s.repoTrash.GetTrashCategoryByName(newItem.TrashName)
		if err != nil {
			continue
		}

		if existing, ok := existingItemMap[newItem.TrashName]; ok {
			existing.Amount = newItem.Amount
			existing.SubTotalEstimatedPrice = newItem.EstimatedSubTotalPrice
			_ = s.repo.UpdateCartItem(existing)
		} else {
			newModelItem := model.CartItem{
				CartID:                 existingCart.ID,
				TrashID:                trash.ID,
				Amount:                 newItem.Amount,
				SubTotalEstimatedPrice: newItem.EstimatedSubTotalPrice,
			}
			_ = s.repo.InsertCartItem(&newModelItem)
		}
	}

	// update cart total amount & price
	existingCart.TotalAmount = cartDTO.TotalAmount
	existingCart.EstimatedTotalPrice = cartDTO.EstimatedTotalPrice
	if err := s.repo.Update(existingCart); err != nil {
		return err
	}

	return utils.DeleteData(redisCartKey(userID))
}


func (s *cartService) DeleteCart(cartID string) error {
	return s.repo.Delete(cartID)
}
