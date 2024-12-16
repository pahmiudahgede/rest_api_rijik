package dto

type RequestPickupResponse struct {
	ID            string           `json:"id"`
	UserID        string           `json:"userId"`
	Request       []RequestItemDTO `json:"request"`
	RequestTime   string           `json:"requestTimePickup"`
	UserAddress   UserAddressDTO   `json:"userAddress"`
	StatusRequest string           `json:"status"`
}

type RequestItemDTO struct {
	TrashCategory   string `json:"trash_category"`
	EstimatedAmount string `json:"estimated_quantity"`
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

func NewRequestPickupResponse(id, userID, requestTime, statusRequest string, request []RequestItemDTO, userAddress UserAddressDTO) RequestPickupResponse {
	return RequestPickupResponse{
		ID:            id,
		UserID:        userID,
		Request:       request,
		RequestTime:   requestTime,
		UserAddress:   userAddress,
		StatusRequest: statusRequest,
	}
}
