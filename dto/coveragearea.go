package dto

type CoverageAreaResponse struct {
	ID        string `json:"id"`
	Province  string `json:"province"`
}

type CoverageAreaWithDistrictsResponse struct {
	ID           string                  `json:"id"`
	Province     string                  `json:"province"`
	CoverageArea []CoverageAreaResponse `json:"coverage_area"`
}

type CoverageAreaDetailWithLocation struct {
	ID               string                    `json:"id"`
	Province         string                    `json:"province"`
	District         string                    `json:"district"`
	LocationSpecific []LocationSpecificResponse `json:"location_specific"`
}

type CoverageDetailResponse struct {
	ID               string                     `json:"id"`
	Province         string                     `json:"province"`
	District         string                     `json:"district"`
	LocationSpecific []LocationSpecificResponse `json:"location_specific"`
}

type LocationSpecificResponse struct {
	ID          string `json:"id"`
	Subdistrict string `json:"subdistrict"`
}

func NewCoverageAreaResponse(id, province string) CoverageAreaResponse {
	return CoverageAreaResponse{
		ID:       id,
		Province: province,
	}
}

func NewCoverageDetailResponse(id, province, district string, locationSpecific []LocationSpecificResponse) CoverageDetailResponse {
	return CoverageDetailResponse{
		ID:               id,
		Province:         province,
		District:         district,
		LocationSpecific: locationSpecific,
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