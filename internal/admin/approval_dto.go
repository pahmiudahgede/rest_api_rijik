package admin

import (
	"rijig/model"
)

type GetAllUsersRequest struct {
	Role      string `query:"role" validate:"required,oneof=masyarakat pengepul pengelola"`
	StatusReg string `query:"statusreg"`
	Page      *int   `query:"page"`
	Limit     *int   `query:"limit"`
}

type UpdateRegistrationStatusRequest struct {
	Action string `json:"action" validate:"required,oneof=approved rejected"`
}

type BaseUserResponse struct {
	ID                   string  `json:"id"`
	Avatar               *string `json:"avatar,omitempty"`
	Name                 string  `json:"name"`
	Gender               string  `json:"gender"`
	Dateofbirth          string  `json:"dateofbirth"`
	Placeofbirth         string  `json:"placeofbirth"`
	Phone                string  `json:"phone"`
	Email                string  `json:"email,omitempty"`
	EmailVerified        bool    `json:"emailVerified"`
	PhoneVerified        bool    `json:"phoneVerified"`
	RoleID               string  `json:"roleId"`
	RoleName             string  `json:"rolename"`
	RegistrationStatus   string  `json:"registrationstatus"`
	RegistrationProgress int8    `json:"registration_progress"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

type MasyarakatUserResponse struct {
	BaseUserResponse
}

type PengepulUserResponse struct {
	BaseUserResponse
	IdentityCard *IdentityCardResponse `json:"identitycard,omitempty"`
}

type PengelolaUserResponse struct {
	BaseUserResponse
	CompanyProfile *CompanyProfileResponse `json:"companyprofile,omitempty"`
}

type IdentityCardResponse struct {
	ID                  string `json:"id"`
	UserID              string `json:"userId"`
	Identificationumber string `json:"identificationumber"`
	Fullname            string `json:"fullname"`
	Placeofbirth        string `json:"placeofbirth"`
	Dateofbirth         string `json:"dateofbirth"`
	Gender              string `json:"gender"`
	BloodType           string `json:"bloodtype"`
	Province            string `json:"province"`
	District            string `json:"district"`
	SubDistrict         string `json:"subdistrict"`
	Hamlet              string `json:"hamlet"`
	Village             string `json:"village"`
	Neighbourhood       string `json:"neighbourhood"`
	PostalCode          string `json:"postalcode"`
	Religion            string `json:"religion"`
	Maritalstatus       string `json:"maritalstatus"`
	Job                 string `json:"job"`
	Citizenship         string `json:"citizenship"`
	Validuntil          string `json:"validuntil"`
	Cardphoto           string `json:"cardphoto"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

type CompanyProfileResponse struct {
	ID                 string `json:"id"`
	UserID             string `json:"userId"`
	CompanyName        string `json:"company_name"`
	CompanyAddress     string `json:"company_address"`
	CompanyPhone       string `json:"company_phone"`
	CompanyEmail       string `json:"company_email,omitempty"`
	CompanyLogo        string `json:"company_logo,omitempty"`
	CompanyWebsite     string `json:"company_website,omitempty"`
	TaxID              string `json:"tax_id,omitempty"`
	FoundedDate        string `json:"founded_date,omitempty"`
	CompanyType        string `json:"company_type,omitempty"`
	CompanyDescription string `json:"company_description"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type UserWithRelations struct {
	User           model.User
	IdentityCard   *model.IdentityCard
	CompanyProfile *model.CompanyProfile
}

type PaginatedUsersResult struct {
	Users []UserWithRelations
	Total int64
}

func ToBaseUserResponse(user model.User) BaseUserResponse {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.RoleName
	}

	return BaseUserResponse{
		ID:                   user.ID,
		Avatar:               user.Avatar,
		Name:                 user.Name,
		Gender:               user.Gender,
		Dateofbirth:          user.Dateofbirth,
		Placeofbirth:         user.Placeofbirth,
		Phone:                user.Phone,
		Email:                user.Email,
		EmailVerified:        user.EmailVerified,
		PhoneVerified:        user.PhoneVerified,
		RoleID:               user.RoleID,
		RoleName:             roleName,
		RegistrationStatus:   user.RegistrationStatus,
		RegistrationProgress: user.RegistrationProgress,
		CreatedAt:            user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToMasyarakatResponse(userWithRelations UserWithRelations) MasyarakatUserResponse {
	return MasyarakatUserResponse{
		BaseUserResponse: ToBaseUserResponse(userWithRelations.User),
	}
}

func ToPengepulResponse(userWithRelations UserWithRelations) PengepulUserResponse {
	response := PengepulUserResponse{
		BaseUserResponse: ToBaseUserResponse(userWithRelations.User),
	}

	if userWithRelations.IdentityCard != nil {
		response.IdentityCard = &IdentityCardResponse{
			ID:                  userWithRelations.IdentityCard.ID,
			UserID:              userWithRelations.IdentityCard.UserID,
			Identificationumber: userWithRelations.IdentityCard.Identificationumber,
			Fullname:            userWithRelations.IdentityCard.Fullname,
			Placeofbirth:        userWithRelations.IdentityCard.Placeofbirth,
			Dateofbirth:         userWithRelations.IdentityCard.Dateofbirth,
			Gender:              userWithRelations.IdentityCard.Gender,
			BloodType:           userWithRelations.IdentityCard.BloodType,
			Province:            userWithRelations.IdentityCard.Province,
			District:            userWithRelations.IdentityCard.District,
			SubDistrict:         userWithRelations.IdentityCard.SubDistrict,
			Hamlet:              userWithRelations.IdentityCard.Hamlet,
			Village:             userWithRelations.IdentityCard.Village,
			Neighbourhood:       userWithRelations.IdentityCard.Neighbourhood,
			PostalCode:          userWithRelations.IdentityCard.PostalCode,
			Religion:            userWithRelations.IdentityCard.Religion,
			Maritalstatus:       userWithRelations.IdentityCard.Maritalstatus,
			Job:                 userWithRelations.IdentityCard.Job,
			Citizenship:         userWithRelations.IdentityCard.Citizenship,
			Validuntil:          userWithRelations.IdentityCard.Validuntil,
			Cardphoto:           userWithRelations.IdentityCard.Cardphoto,
			CreatedAt:           userWithRelations.IdentityCard.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:           userWithRelations.IdentityCard.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response
}

func ToPengelolaResponse(userWithRelations UserWithRelations) PengelolaUserResponse {
	response := PengelolaUserResponse{
		BaseUserResponse: ToBaseUserResponse(userWithRelations.User),
	}

	if userWithRelations.CompanyProfile != nil {
		response.CompanyProfile = &CompanyProfileResponse{
			ID:                 userWithRelations.CompanyProfile.ID,
			UserID:             userWithRelations.CompanyProfile.UserID,
			CompanyName:        userWithRelations.CompanyProfile.CompanyName,
			CompanyAddress:     userWithRelations.CompanyProfile.CompanyAddress,
			CompanyPhone:       userWithRelations.CompanyProfile.CompanyPhone,
			CompanyEmail:       userWithRelations.CompanyProfile.CompanyEmail,
			CompanyLogo:        userWithRelations.CompanyProfile.CompanyLogo,
			CompanyWebsite:     userWithRelations.CompanyProfile.CompanyWebsite,
			TaxID:              userWithRelations.CompanyProfile.TaxID,
			FoundedDate:        userWithRelations.CompanyProfile.FoundedDate,
			CompanyType:        userWithRelations.CompanyProfile.CompanyType,
			CompanyDescription: userWithRelations.CompanyProfile.CompanyDescription,
			CreatedAt:          userWithRelations.CompanyProfile.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:          userWithRelations.CompanyProfile.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return response
}
