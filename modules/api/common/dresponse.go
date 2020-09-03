package common

//DefaultResponse default payload response
type DefaultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//NewInternalServerErrorResponse default internal server error response
func NewInternalServerErrorResponse() DefaultResponse {
	return DefaultResponse{
		500,
		"Internal server error",
	}
}

//NewNotFoundResponse default not found error response
func NewNotFoundResponse() DefaultResponse {
	return DefaultResponse{
		404,
		"Not found",
	}
}

//NewBadRequestResponse default not found error response
func NewBadRequestResponse() DefaultResponse {
	return DefaultResponse{
		400,
		"Bad request",
	}
}

//NewConflictResponse default not found error response
func NewConflictResponse() DefaultResponse {
	return DefaultResponse{
		409,
		"Data has been modified",
	}
}
