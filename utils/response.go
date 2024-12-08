package utils

type Meta struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type ApiResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func FormatResponse(statusCode int, message string, data interface{}) ApiResponse {
	return ApiResponse{
		Meta: Meta{
			StatusCode: statusCode,
			Message:    message,
		},
		Data: data,
	}
}
