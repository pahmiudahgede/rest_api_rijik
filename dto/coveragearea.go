package dto

type CoverageAreaResponse struct {
	ID       string `json:"id"`
	Province string `json:"province"`
}

type CoverageAreaWithDistrictsResponse struct {
	ID           string                 `json:"id"`
	Province     string                 `json:"province"`
	CoverageArea []CoverageAreaResponse `json:"coverage_area"`
}

type CoverageAreaDetailWithLocation struct {
	ID               string                     `json:"id"`
	Province         string                     `json:"province"`
	District         string                     `json:"district"`
	LocationSpecific []LocationSpecificResponse `json:"location_specific"`
}

type LocationSpecificResponse struct {
	ID          string `json:"id"`
	Subdistrict string `json:"subdistrict"`
}

type CoverageDetailResponse struct {
	ID        string `json:"id"`
	Province  string `json:"province"`
	District  string `json:"district"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func NewCoverageDetailResponse(id, province, district, createdAt, updatedAt string) CoverageDetailResponse {
	return CoverageDetailResponse{
		ID:        id,
		Province:  province,
		District:  district,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewCoverageAreaResponse(id, province string) CoverageAreaResponse {
	return CoverageAreaResponse{
		ID:       id,
		Province: province,
	}
}

func NewLocationSpecificResponse(id, subdistrict string) LocationSpecificResponse {
	return LocationSpecificResponse{
		ID:          id,
		Subdistrict: subdistrict,
	}
}

func NewCoverageAreaWithDistrictsResponse(id, province string, coverageArea []CoverageAreaResponse) CoverageAreaWithDistrictsResponse {
	return CoverageAreaWithDistrictsResponse{
		ID:           id,
		Province:     province,
		CoverageArea: coverageArea,
	}
}

func NewCoverageAreaDetailWithLocation(id, province, district string, locationSpecific []LocationSpecificResponse) CoverageAreaDetailWithLocation {
	return CoverageAreaDetailWithLocation{
		ID:               id,
		Province:         province,
		District:         district,
		LocationSpecific: locationSpecific,
	}
}

type CoverageAreaRequest struct {
	Province string `json:"province" validate:"required"`
}

type CoverageDetailRequest struct {
	CoverageAreaID string `json:"coverage_area_id" validate:"required"`
	Province       string `json:"province" validate:"required"`
	District       string `json:"district" validate:"required"`
}

type LocationSpecificRequest struct {
	CoverageDetailID string `json:"coverage_detail_id" validate:"required"`
	Subdistrict      string `json:"subdistrict" validate:"required"`
}