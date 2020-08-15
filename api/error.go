package api

//DefaultErrorResponse default error response
type DefaultErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//NewInternalServerErrorResponse default internal server error response
func NewInternalServerErrorResponse() DefaultErrorResponse {
	return DefaultErrorResponse{
		500,
		"Internal server error",
	}
}

//NewNotFoundResponse default not found error response
func NewNotFoundResponse() DefaultErrorResponse {
	return DefaultErrorResponse{
		404,
		"Not found",
	}
}

//NewBadRequestResponse default not found error response
func NewBadRequestResponse() DefaultErrorResponse {
	return DefaultErrorResponse{
		400,
		"Bad request",
	}
}

//NewConflictResponse default not found error response
func NewConflictResponse() DefaultErrorResponse {
	return DefaultErrorResponse{
		409,
		"Data has been modified",
	}
}
