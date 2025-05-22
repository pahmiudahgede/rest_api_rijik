package dto

import (
	"strings"
)

type ResponseIdentityCardDTO struct {
	ID                  string `json:"id"`
	UserID              string `json:"userId"`
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
	Cardphoto           string `json:"cardphoto"`
	CreatedAt           string `json:"createdAt"`
	UpdatedAt           string `json:"updatedAt"`
}

type RequestIdentityCardDTO struct {
	UserID              string `json:"userId"`
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
	Cardphoto           string `json:"cardphoto"`
}

func (r *RequestIdentityCardDTO) ValidateIdentityCardInput() (map[string][]string, bool) {
	errors := make(map[string][]string)
	isValid := true

	if strings.TrimSpace(r.Identificationumber) == "" {
		errors["identificationumber"] = append(errors["identificationumber"], "Nomor identifikasi harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Placeofbirth) == "" {
		errors["placeofbirth"] = append(errors["placeofbirth"], "Tempat lahir harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Dateofbirth) == "" {
		errors["dateofbirth"] = append(errors["dateofbirth"], "Tanggal lahir harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Gender) == "" {
		errors["gender"] = append(errors["gender"], "Jenis kelamin harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.BloodType) == "" {
		errors["bloodtype"] = append(errors["bloodtype"], "Golongan darah harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.District) == "" {
		errors["district"] = append(errors["district"], "Kecamatan harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Village) == "" {
		errors["village"] = append(errors["village"], "Desa harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Neighbourhood) == "" {
		errors["neighbourhood"] = append(errors["neighbourhood"], "RT/RW harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Religion) == "" {
		errors["religion"] = append(errors["religion"], "Agama harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Maritalstatus) == "" {
		errors["maritalstatus"] = append(errors["maritalstatus"], "Status pernikahan harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Job) == "" {
		errors["job"] = append(errors["job"], "Pekerjaan harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Citizenship) == "" {
		errors["citizenship"] = append(errors["citizenship"], "Kewarganegaraan harus diisi")
		isValid = false
	}

	if strings.TrimSpace(r.Validuntil) == "" {
		errors["validuntil"] = append(errors["validuntil"], "Masa berlaku harus diisi")
		isValid = false
	}

	return errors, isValid
}
