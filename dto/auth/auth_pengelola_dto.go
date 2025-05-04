package dto

import (
	"regexp"
	"rijig/utils"
)

type LoginPengelolaRequest struct {
	Phone string `json:"phone"`
}

func (r *LoginPengelolaRequest) ValidateLogin() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.Phone == "" {
		errors["phone"] = append(errors["phone"], "Phone number is required")
	} else if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Phone number is not valid")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

type VerifLoginPengelolaRequest struct {
	Phone string `json:"phone"`
	Otp   string `json:"verif_otp"`
}

func (r *VerifLoginPengelolaRequest) ValidateVerifLogin() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.Phone == "" {
		errors["phone"] = append(errors["phone"], "Phone number is required")
	} else if !utils.IsValidPhoneNumber(r.Phone) {
		errors["phone"] = append(errors["phone"], "Phone number is not valid")
	}

	if r.Otp == "" {
		errors["otp"] = append(errors["otp"], "OTP is required")
	} else if len(r.Otp) != 6 {
		errors["otp"] = append(errors["otp"], "OTP must be 6 digits")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

type LoginPengelolaResponse struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Token  string `json:"token"`
}

type PengelolaIdentityCard struct {
	Cardphoto           string `json:"cardphoto"`
	Identificationumber string `json:"identificationumber"`
	Placeofbirth        string `json:"placeofbirth"`
	Dateofbirth         string `json:"dateofbirth"`
	Gender              string `json:"gender"`
	BloodType           string `json:"bloodtype"`
	District            string `json:"district"`
	Village             string `json:"village"`
	Neighbourhood       string `json:"neighbourhood"`
	Religion            string `json:"religion"`
	Maritalstatus       string `json:"maritalstatus"`
	Job                 string `json:"job"`
	Citizenship         string `json:"citizenship"`
	Validuntil          string `json:"validuntil"`
}

func (r *PengelolaIdentityCard) ValidateIDcard() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.Cardphoto == "" {
		errors["cardphoto"] = append(errors["cardphoto"], "Card photo is required")
	}

	if r.Identificationumber == "" {
		errors["identificationumber"] = append(errors["identificationumber"], "Identification number is required")
	}

	if r.Dateofbirth == "" {
		errors["dateofbirth"] = append(errors["dateofbirth"], "Date of birth is required")
	} else if !isValidDate(r.Dateofbirth) {
		errors["dateofbirth"] = append(errors["dateofbirth"], "Date of birth must be in DD-MM-YYYY format")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

type PengelolaCompanyProfile struct {
	CompanyName  string `json:"company_name"`
	CompanyPhone string `json:"company_phone"`
	CompanyEmail string `json:"company_email"`
}

func (r *PengelolaCompanyProfile) ValidateCompany() (map[string][]string, bool) {
	errors := make(map[string][]string)

	if r.CompanyName == "" {
		errors["company_name"] = append(errors["company_name"], "Company name is required")
	}

	if r.CompanyPhone == "" {
		errors["company_phone"] = append(errors["company_phone"], "Company phone is required")
	} else if !utils.IsValidPhoneNumber(r.CompanyPhone) {
		errors["company_phone"] = append(errors["company_phone"], "Invalid phone number format")
	}

	if r.CompanyEmail == "" {
		errors["company_email"] = append(errors["company_email"], "Company email is required")
	} else if !utils.IsValidEmail(r.CompanyEmail) {
		errors["company_email"] = append(errors["company_email"], "Invalid email format")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}

func isValidDate(date string) bool {
	re := regexp.MustCompile(`^\d{2}-\d{2}-\d{4}$`)
	return re.MatchString(date)
}
