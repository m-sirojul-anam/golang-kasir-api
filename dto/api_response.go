package dto

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  "Success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) APIResponse {
	return APIResponse{
		Status:  "Error",
		Message: message,
		Data:    nil,
	}
}
