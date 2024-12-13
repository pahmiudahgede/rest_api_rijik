package dto

type CoverageAreaResponse struct {
	ID        string `json:"id"`
	Province  string `json:"province"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CoverageAreaWithDistrictsResponse struct {
	ID           string                 `json:"id"`
	Province     string                 `json:"province"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
	CoverageArea []CoverageAreaResponse `json:"coverage_area"`
}

type CoverageAreaDetailWithLocation struct {
	ID          string                `json:"id"`
	Province    string                `json:"province"`
	District    string                `json:"district"`
	CreatedAt   string                `json:"createdAt"`
	UpdatedAt   string                `json:"updatedAt"`
	Subdistrict []SubdistrictResponse `json:"subdistrict"`
}

type SubdistrictResponse struct {
	ID          string `json:"id"`
	Subdistrict string `json:"subdistrict"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

func NewCoverageAreaResponse(id, province, createdAt, updatedAt string) CoverageAreaResponse {
	return CoverageAreaResponse{
		ID:        id,
		Province:  province,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewCoverageAreaWithDistrictsResponse(id, province, createdAt, updatedAt string, coverageArea []CoverageAreaResponse) CoverageAreaWithDistrictsResponse {
	return CoverageAreaWithDistrictsResponse{
		ID:           id,
		Province:     province,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		CoverageArea: coverageArea,
	}
}

func NewCoverageAreaDetailWithLocation(id, province, district, createdAt, updatedAt string, subdistricts []SubdistrictResponse) CoverageAreaDetailWithLocation {
	return CoverageAreaDetailWithLocation{
		ID:          id,
		Province:    province,
		District:    district,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Subdistrict: subdistricts,
	}
}

func NewSubdistrictResponse(id, subdistrict, createdAt, updatedAt string) SubdistrictResponse {
	return SubdistrictResponse{
		ID:          id,
		Subdistrict: subdistrict,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
