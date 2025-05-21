package dto

import (
	"fmt"
	"strings"
)

type RequestCollectorDTO struct {
	AddressId               string                           `json:"address_id"`
	AvaibleTrashbyCollector []RequestAvaibleTrashbyCollector `json:"avaible_trash"`
}

type RequestAvaibleTrashbyCollector struct {
	TrashId    string  `json:"trash_id"`
	TrashPrice float32 `json:"trash_price"`
}

type RequestAddAvaibleTrash struct {
	AvaibleTrash []RequestAvaibleTrashbyCollector `json:"avaible_trash"`
}

type SelectCollectorRequest struct {
	Collector_id string `json:"collector_id"`
}

func (r *SelectCollectorRequest) ValidateSelectCollectorRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.Collector_id) == "" {
		errors["collector_id"] = append(errors["collector_id"], "collector_id harus diisi")
	}
	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *RequestAddAvaibleTrash) ValidateRequestAddAvaibleTrash() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if len(r.AvaibleTrash) == 0 {
		errors["avaible_trash"] = append(errors["avaible_trash"], "tidak boleh kosong")
	}

	for i, trash := range r.AvaibleTrash {
		if strings.TrimSpace(trash.TrashId) == "" {
			errors[fmt.Sprintf("avaible_trash[%d].trash_id", i)] = append(errors[fmt.Sprintf("avaible_trash[%d].trash_id", i)], "trash_id tidak boleh kosong")
		}
		if trash.TrashPrice <= 0 {
			errors[fmt.Sprintf("avaible_trash[%d].trash_price", i)] = append(errors[fmt.Sprintf("avaible_trash[%d].trash_price", i)], "trash_price harus lebih dari 0")
		}
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

type ResponseCollectorDTO struct {
	ID                      string                            `json:"collector_id"`
	UserId                  string                            `json:"user_id"`
	User                    *UserResponseDTO                  `json:"user,omitempty"`
	AddressId               string                            `json:"address_id"`
	Address                 *AddressResponseDTO               `json:"address,omitempty"`
	JobStatus               *string                           `json:"job_status,omitempty"`
	Rating                  float32                           `json:"rating"`
	AvaibleTrashbyCollector []ResponseAvaibleTrashByCollector `json:"avaible_trash"`
}

type ResponseAvaibleTrashByCollector struct {
	ID         string  `json:"id"`
	TrashId    string  `json:"trash_id"`
	TrashName  string  `json:"trash_name"`
	TrashIcon  string  `json:"trash_icon"`
	TrashPrice float32 `json:"trash_price"`
}

func (r *RequestCollectorDTO) ValidateRequestCollector() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.AddressId) == "" {
		errors["address_id"] = append(errors["address_id"], "address_id harus diisi")
	}

	for i, trash := range r.AvaibleTrashbyCollector {
		if strings.TrimSpace(trash.TrashId) == "" {
			errors[fmt.Sprintf("avaible_trash[%d].trash_id", i)] = append(errors[fmt.Sprintf("avaible_trash[%d].trash_id", i)], "trash_id tidak boleh kosong")
		}
		if trash.TrashPrice <= 0 {
			errors[fmt.Sprintf("avaible_trash[%d].trash_price", i)] = append(errors[fmt.Sprintf("avaible_trash[%d].trash_price", i)], "trash_price harus lebih dari 0")
		}
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
