package dto

type RequestPickupRequest struct {
	RequestItems  []RequestItemDTO `json:"request_items"`
	RequestTime   string           `json:"requestTime"`
	UserAddressID string           `json:"userAddressId"`
}

type RequestPickupResponse struct {
	ID            string           `json:"id"`
	UserID        string           `json:"userId"`
	Request       []RequestItemDTO `json:"request"`
	RequestTime   string           `json:"requestTimePickup"`
	UserAddress   UserAddressDTO   `json:"userAddress"`
	StatusRequest string           `json:"status"`
	CreatedAt     string           `json:"createdAt"`
	UpdatedAt     string           `json:"updatedAt"`
}

type RequestItemDTO struct {
	TrashCategory   string `json:"trashCategory"`
	EstimatedAmount string `json:"estimatedAmount"`
}

type UserAddressDTO struct {
	Province    string `json:"province"`
	District    string `json:"district"`
	Subdistrict string `json:"subdistrict"`
	PostalCode  int    `json:"postalCode"`
	Village     string `json:"village"`
	Detail      string `json:"detail"`
	Geography   string `json:"geography"`
}

func NewRequestPickupResponse(id, userID, requestTime, statusRequest string, request []RequestItemDTO, userAddress UserAddressDTO, createdAt, updatedAt string) RequestPickupResponse {
	return RequestPickupResponse{
		ID:            id,
		UserID:        userID,
		Request:       request,
		RequestTime:   requestTime,
		UserAddress:   userAddress,
		StatusRequest: statusRequest,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}