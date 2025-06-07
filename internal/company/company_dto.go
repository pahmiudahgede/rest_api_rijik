package company

import (
	"rijig/utils"
	"strings"
)

type ResponseCompanyProfileDTO struct {
	ID                 string `json:"id"`
	UserID             string `json:"userId"`
	CompanyName        string `json:"company_name"`
	CompanyAddress     string `json:"company_address"`
	CompanyPhone       string `json:"company_phone"`
	CompanyEmail       string `json:"company_email"`
	CompanyLogo        string `json:"company_logo,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	TaxID              string `json:"taxId,omitempty"`
	FoundedDate        string `json:"founded_date,omitempty"`
	CompanyType        string `json:"company_type,omitempty"`
	CompanyDescription string `json:"company_description"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
}

type RequestCompanyProfileDTO struct {
	CompanyName        string `json:"company_name"`
	CompanyAddress     string `json:"company_address"`
	CompanyPhone       string `json:"company_phone"`
	CompanyEmail       string `json:"company_email"`
	CompanyLogo        string `json:"company_logo,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	TaxID              string `json:"taxId,omitempty"`
	FoundedDate        string `json:"founded_date,omitempty"`
	CompanyType        string `json:"company_type,omitempty"`
	CompanyDescription string `json:"company_description"`
}

func (r *RequestCompanyProfileDTO) ValidateCompanyProfileInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if strings.TrimSpace(r.CompanyName) == "" {
		errors["company_Name"] = append(errors["company_name"], "Company name is required")
	}

	if strings.TrimSpace(r.CompanyAddress) == "" {
		errors["company_Address"] = append(errors["company_address"], "Company address is required")
	}

	if !utils.IsValidPhoneNumber(r.CompanyPhone) {
		errors["company_Phone"] = append(errors["company_phone"], "nomor harus dimulai 62.. dan 8-14 digit")
	}

	if strings.TrimSpace(r.CompanyDescription) == "" {
		errors["company_Description"] = append(errors["company_description"], "Company description is required")
	}

	if len(errors) > 0 {
		return errors, false
	}

	return nil, true
}
