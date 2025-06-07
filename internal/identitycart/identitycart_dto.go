package identitycart

import (
	"rijig/utils"
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

type RequestIdentityCardDTO struct {
	DeviceID            string `json:"device_id"`
	UserID              string `json:"userId"`
	Identificationumber string `json:"identificationumber"`
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
}

func (r *RequestIdentityCardDTO) ValidateIdentityCardInput() (map[string][]string, bool) {
	errors := make(map[string][]string)

	r.Placeofbirth = strings.ToLower(r.Placeofbirth)
	r.Dateofbirth = strings.ToLower(r.Dateofbirth)
	r.Gender = strings.ToLower(r.Gender)
	r.BloodType = strings.ToUpper(r.BloodType)
	r.Province = strings.ToLower(r.Province)
	r.District = strings.ToLower(r.District)
	r.SubDistrict = strings.ToLower(r.SubDistrict)
	r.Hamlet = strings.ToLower(r.Hamlet)
	r.Village = strings.ToLower(r.Village)
	r.Neighbourhood = strings.ToLower(r.Neighbourhood)
	r.PostalCode = strings.ToLower(r.PostalCode)
	r.Religion = strings.ToLower(r.Religion)
	r.Maritalstatus = strings.ToLower(r.Maritalstatus)
	r.Job = strings.ToLower(r.Job)
	r.Citizenship = strings.ToLower(r.Citizenship)
	r.Validuntil = strings.ToLower(r.Validuntil)

	nikData := utils.FetchNIKData(r.Identificationumber)
	if strings.ToLower(nikData.Status) != "sukses" {
		errors["identificationumber"] = append(errors["identificationumber"], "NIK yang anda masukkan tidak valid")
	} else {

		if r.Dateofbirth != strings.ToLower(nikData.Ttl) {
			errors["dateofbirth"] = append(errors["dateofbirth"], "Tanggal lahir tidak sesuai dengan NIK")
		}

		if r.Gender != strings.ToLower(nikData.Sex) {
			errors["gender"] = append(errors["gender"], "Jenis kelamin tidak sesuai dengan NIK")
		}

		if r.Province != strings.ToLower(nikData.Provinsi) {
			errors["province"] = append(errors["province"], "Provinsi tidak sesuai dengan NIK")
		}

		if r.District != strings.ToLower(nikData.Kabkot) {
			errors["district"] = append(errors["district"], "Kabupaten/Kota tidak sesuai dengan NIK")
		}

		if r.SubDistrict != strings.ToLower(nikData.Kecamatan) {
			errors["subdistrict"] = append(errors["subdistrict"], "Kecamatan tidak sesuai dengan NIK")
		}

		if r.PostalCode != strings.ToLower(nikData.KodPos) {
			errors["postalcode"] = append(errors["postalcode"], "Kode pos tidak sesuai dengan NIK")
		}
	}

	if r.Placeofbirth == "" {
		errors["placeofbirth"] = append(errors["placeofbirth"], "Tempat lahir wajib diisi")
	}
	if r.Hamlet == "" {
		errors["hamlet"] = append(errors["hamlet"], "Dusun/RW wajib diisi")
	}
	if r.Village == "" {
		errors["village"] = append(errors["village"], "Desa/Kelurahan wajib diisi")
	}
	if r.Neighbourhood == "" {
		errors["neighbourhood"] = append(errors["neighbourhood"], "RT wajib diisi")
	}
	if r.Job == "" {
		errors["job"] = append(errors["job"], "Pekerjaan wajib diisi")
	}
	if r.Citizenship == "" {
		errors["citizenship"] = append(errors["citizenship"], "Kewarganegaraan wajib diisi")
	}
	if r.Validuntil == "" {
		errors["validuntil"] = append(errors["validuntil"], "Berlaku hingga wajib diisi")
	}

	validBloodTypes := map[string]bool{"A": true, "B": true, "O": true, "AB": true}
	if _, ok := validBloodTypes[r.BloodType]; !ok {
		errors["bloodtype"] = append(errors["bloodtype"], "Golongan darah harus A, B, O, atau AB")
	}

	validReligions := map[string]bool{
		"islam": true, "kristen": true, "katolik": true, "hindu": true, "buddha": true, "konghucu": true,
	}
	if _, ok := validReligions[r.Religion]; !ok {
		errors["religion"] = append(errors["religion"], "Agama harus salah satu dari Islam, Kristen, Katolik, Hindu, Buddha, atau Konghucu")
	}

	if r.Maritalstatus != "kawin" && r.Maritalstatus != "belum kawin" {
		errors["maritalstatus"] = append(errors["maritalstatus"], "Status perkawinan harus 'kawin' atau 'belum kawin'")
	}

	if len(errors) > 0 {
		return errors, false
	}
	return nil, true
}
