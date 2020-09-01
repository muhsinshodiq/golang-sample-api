package response

//CreateNewItemResponse Create item response payload
type CreateNewItemResponse struct {
	ID string `json:"id"`
}

//NewCreateNewItemResponse construct CreateNewItemResponse
func NewCreateNewItemResponse(id string) *CreateNewItemResponse {
	return &CreateNewItemResponse{
		id,
	}
}
