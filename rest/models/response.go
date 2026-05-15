package models

// APIResponse represents the standard API response structure
type APIResponse struct {
	ResponseMessage string      `json:"responseMessage"`
	ResponseCode    string      `json:"responseCode"`
	ResponseData    interface{} `json:"responseData,omitempty"`
}

// NewSuccessResponse creates a successful API response
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		ResponseMessage: "success",
		ResponseCode:    "00",
		ResponseData:    data,
	}
}

// NewErrorResponse creates an error API response
func NewErrorResponse() APIResponse {
	return APIResponse{
		ResponseMessage: "internal error",
		ResponseCode:    "99",
	}
}
