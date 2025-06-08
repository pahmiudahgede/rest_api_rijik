package collector

import (
	"fmt"
	"rijig/internal/address"
	"rijig/internal/trash"
	"strings"
	"time"
)

type CreateCollectorRequest struct {
	UserID              string                        `json:"user_id" binding:"required"`
	JobStatus           string                        `json:"job_status,omitempty"`
	AddressID           string                        `json:"address_id" binding:"required"`
	AvailableTrashItems []CreateAvailableTrashRequest `json:"available_trash_items,omitempty"`
}

type UpdateCollectorRequest struct {
	JobStatus           string                        `json:"job_status,omitempty"`
	AddressID           string                        `json:"address_id,omitempty"`
	AvailableTrashItems []CreateAvailableTrashRequest `json:"available_trash_items,omitempty"`
}

type CreateAvailableTrashRequest struct {
	TrashCategoryID string  `json:"trash_category_id" binding:"required"`
	Price           float32 `json:"price" binding:"required"`
}

type UpdateAvailableTrashRequest struct {
	ID              string  `json:"id,omitempty"`
	TrashCategoryID string  `json:"trash_category_id,omitempty"`
	Price           float32 `json:"price,omitempty"`
}

type CollectorResponse struct {
	ID             string                      `json:"id"`
	UserID         string                      `json:"user_id"`
	JobStatus      string                      `json:"job_status"`
	Rating         float32                     `json:"rating"`
	AddressID      string                      `json:"address_id"`
	Address        *address.AddressResponseDTO `json:"address,omitempty"`
	AvailableTrash []AvailableTrashResponse    `json:"available_trash"`
	CreatedAt      string                      `json:"created_at"`
	UpdatedAt      string                      `json:"updated_at"`
}

type AvailableTrashResponse struct {
	ID              string                          `json:"id"`
	CollectorID     string                          `json:"collector_id"`
	TrashCategoryID string                          `json:"trash_category_id"`
	TrashCategory   *trash.ResponseTrashCategoryDTO `json:"trash_category,omitempty"`
	Price           float32                         `json:"price"`
}

func (r *CreateCollectorRequest) ValidateCreateCollectorRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.UserID) == "" {
		errors["user_id"] = append(errors["user_id"], "User ID tidak boleh kosong")
	}

	if strings.TrimSpace(r.AddressID) == "" {
		errors["address_id"] = append(errors["address_id"], "Address ID tidak boleh kosong")
	}

	if r.JobStatus != "" {
		r.JobStatus = strings.ToLower(strings.TrimSpace(r.JobStatus))
		if r.JobStatus != "active" && r.JobStatus != "inactive" {
			errors["job_status"] = append(errors["job_status"], "Job status hanya boleh 'active' atau 'inactive'")
		}
	} else {
		r.JobStatus = "inactive"
	}

	if len(r.AvailableTrashItems) > 0 {
		trashCategoryMap := make(map[string]bool)
		for i, item := range r.AvailableTrashItems {
			fieldPrefix := fmt.Sprintf("available_trash_items[%d]", i)

			if strings.TrimSpace(item.TrashCategoryID) == "" {
				errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID tidak boleh kosong")
			} else {

				if trashCategoryMap[item.TrashCategoryID] {
					errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID sudah ada dalam daftar")
				} else {
					trashCategoryMap[item.TrashCategoryID] = true
				}
			}

			if item.Price <= 0 {
				errors[fieldPrefix+".price"] = append(errors[fieldPrefix+".price"], "Harga harus lebih dari 0")
			}
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *UpdateCollectorRequest) ValidateUpdateCollectorRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.JobStatus != "" {
		r.JobStatus = strings.ToLower(strings.TrimSpace(r.JobStatus))
		if r.JobStatus != "active" && r.JobStatus != "inactive" {
			errors["job_status"] = append(errors["job_status"], "Job status hanya boleh 'active' atau 'inactive'")
		}
	}

	if r.AddressID != "" && strings.TrimSpace(r.AddressID) == "" {
		errors["address_id"] = append(errors["address_id"], "Address ID tidak boleh kosong jika disediakan")
	}

	if len(r.AvailableTrashItems) > 0 {
		trashCategoryMap := make(map[string]bool)
		for i, item := range r.AvailableTrashItems {
			fieldPrefix := fmt.Sprintf("available_trash_items[%d]", i)

			if strings.TrimSpace(item.TrashCategoryID) == "" {
				errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID tidak boleh kosong")
			} else {

				if trashCategoryMap[item.TrashCategoryID] {
					errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID sudah ada dalam daftar")
				} else {
					trashCategoryMap[item.TrashCategoryID] = true
				}
			}

			if item.Price <= 0 {
				errors[fieldPrefix+".price"] = append(errors[fieldPrefix+".price"], "Harga harus lebih dari 0")
			}
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *CreateAvailableTrashRequest) ValidateCreateAvailableTrashRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.TrashCategoryID) == "" {
		errors["trash_category_id"] = append(errors["trash_category_id"], "Trash category ID tidak boleh kosong")
	}

	if r.Price <= 0 {
		errors["price"] = append(errors["price"], "Harga harus lebih dari 0")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *UpdateAvailableTrashRequest) ValidateUpdateAvailableTrashRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.TrashCategoryID != "" && strings.TrimSpace(r.TrashCategoryID) == "" {
		errors["trash_category_id"] = append(errors["trash_category_id"], "Trash category ID tidak boleh kosong jika disediakan")
	}

	if r.Price != 0 && r.Price <= 0 {
		errors["price"] = append(errors["price"], "Harga harus lebih dari 0 jika disediakan")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}

func (r *CreateCollectorRequest) IsValidJobStatus(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return status == "active" || status == "inactive"
}

func (r *UpdateCollectorRequest) IsValidJobStatus(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return status == "active" || status == "inactive"
}

func (r *CollectorResponse) FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

func (r *CreateCollectorRequest) SetDefaults() {
	if r.JobStatus == "" {
		r.JobStatus = "inactive"
	} else {
		r.JobStatus = strings.ToLower(strings.TrimSpace(r.JobStatus))
	}
}

func (r *UpdateCollectorRequest) NormalizeJobStatus() {
	if r.JobStatus != "" {
		r.JobStatus = strings.ToLower(strings.TrimSpace(r.JobStatus))
	}
}

type BulkUpdateAvailableTrashRequest struct {
	CollectorID         string                        `json:"collector_id" binding:"required"`
	AvailableTrashItems []CreateAvailableTrashRequest `json:"available_trash_items" binding:"required"`
}

func (r *BulkUpdateAvailableTrashRequest) ValidateBulkUpdateAvailableTrashRequest() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.CollectorID) == "" {
		errors["collector_id"] = append(errors["collector_id"], "Collector ID tidak boleh kosong")
	}

	if len(r.AvailableTrashItems) == 0 {
		errors["available_trash_items"] = append(errors["available_trash_items"], "Minimal harus ada 1 item sampah")
	} else {
		trashCategoryMap := make(map[string]bool)
		for i, item := range r.AvailableTrashItems {
			fieldPrefix := fmt.Sprintf("available_trash_items[%d]", i)

			if strings.TrimSpace(item.TrashCategoryID) == "" {
				errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID tidak boleh kosong")
			} else {

				if trashCategoryMap[item.TrashCategoryID] {
					errors[fieldPrefix+".trash_category_id"] = append(errors[fieldPrefix+".trash_category_id"], "Trash category ID sudah ada dalam daftar")
				} else {
					trashCategoryMap[item.TrashCategoryID] = true
				}
			}

			if item.Price <= 0 {
				errors[fieldPrefix+".price"] = append(errors[fieldPrefix+".price"], "Harga harus lebih dari 0")
			}
		}
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
