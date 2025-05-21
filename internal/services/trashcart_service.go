package services

import (
	"log"
	"time"

	"rijig/dto"
	"rijig/internal/repositories"
	"rijig/model"

	"github.com/google/uuid"
)

type CartService struct {
	Repo repositories.CartRepository
}

func NewCartService(repo repositories.CartRepository) *CartService {
	return &CartService{Repo: repo}
}

func (s *CartService) CommitCartToDatabase(userID string) error {
	items, err := GetCartItems(userID)
	if err != nil || len(items) == 0 {
		log.Printf("No items to commit for user: %s", userID)
		return err
	}

	var cartItems []model.CartItem
	var totalAmount float32
	var estimatedTotal float32

	for _, item := range items {
		trash, err := s.Repo.GetTrashCategoryByID(item.TrashCategoryID)
		if err != nil {
			log.Printf("Trash category not found for trashID: %s", item.TrashCategoryID)
			continue
		}

		subTotal := float32(trash.EstimatedPrice) * item.Amount
		totalAmount += item.Amount
		estimatedTotal += subTotal

		cartItems = append(cartItems, model.CartItem{
			ID:                     uuid.NewString(),
			TrashCategoryID:                item.TrashCategoryID,
			Amount:                 item.Amount,
			SubTotalEstimatedPrice: subTotal,
		})
	}

	cart := &model.Cart{
		ID:                  uuid.NewString(),
		UserID:              userID,
		CartItems:           cartItems,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: estimatedTotal,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.Repo.DeleteCartByUserID(userID); err != nil {
		log.Printf("Failed to delete old cart: %v", err)
	}

	if err := s.Repo.CreateCart(cart); err != nil {
		log.Printf("Failed to create cart: %v", err)
		return err
	}

	if err := ClearCart(userID); err != nil {
		log.Printf("Failed to clear Redis cart: %v", err)
	}

	log.Printf("Cart committed successfully for user: %s", userID)
	return nil
}

func (s *CartService) GetCartFromRedis(userID string) (*dto.CartResponse, error) {
	items, err := GetCartItems(userID)
	if err != nil || len(items) == 0 {
		return nil, err
	}

	var totalAmount float32
	var estimatedTotal float32
	var cartItemDTOs []dto.CartItemResponse

	for _, item := range items {
		trash, err := s.Repo.GetTrashCategoryByID(item.TrashCategoryID)
		if err != nil {
			continue
		}

		subtotal := float32(trash.EstimatedPrice) * item.Amount
		totalAmount += item.Amount
		estimatedTotal += subtotal

		cartItemDTOs = append(cartItemDTOs, dto.CartItemResponse{
			TrashId:                trash.ID,
			TrashIcon:              trash.Icon,
			TrashName:              trash.Name,
			Amount:                 item.Amount,
			EstimatedSubTotalPrice: subtotal,
		})
	}

	resp := &dto.CartResponse{
		ID:                  "N/A",
		UserID:              userID,
		TotalAmount:         totalAmount,
		EstimatedTotalPrice: estimatedTotal,
		CreatedAt:           time.Now().Format(time.RFC3339),
		UpdatedAt:           time.Now().Format(time.RFC3339),
		CartItems:           cartItemDTOs,
	}
	return resp, nil
}

func (s *CartService) GetCart(userID string) (*dto.CartResponse, error) {

	cartRedis, err := s.GetCartFromRedis(userID)
	if err == nil && len(cartRedis.CartItems) > 0 {
		return cartRedis, nil
	}

	cartDB, err := s.Repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	var items []dto.CartItemResponse
	for _, item := range cartDB.CartItems {
		items = append(items, dto.CartItemResponse{
			ItemId:                 item.ID,
			TrashId:                item.TrashCategoryID,
			TrashIcon:              item.TrashCategory.Icon,
			TrashName:              item.TrashCategory.Name,
			Amount:                 item.Amount,
			EstimatedSubTotalPrice: item.SubTotalEstimatedPrice,
		})
	}

	resp := &dto.CartResponse{
		ID:                  cartDB.ID,
		UserID:              cartDB.UserID,
		TotalAmount:         cartDB.TotalAmount,
		EstimatedTotalPrice: cartDB.EstimatedTotalPrice,
		CreatedAt:           cartDB.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           cartDB.UpdatedAt.Format(time.RFC3339),
		CartItems:           items,
	}
	return resp, nil
}
