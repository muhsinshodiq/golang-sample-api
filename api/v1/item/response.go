package item

import (
	"sample-order/core/item"
	"time"
)

//GetItemByIDResponse Get item by ID response payload
type GetItemByIDResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	Version     int       `json:"version"`
}

//NewGetItemByIDResponse construct GetItemByIDResponse
func NewGetItemByIDResponse(item *item.Item) *GetItemByIDResponse {
	var itemResponse GetItemByIDResponse
	itemResponse.ID = item.ID
	itemResponse.Name = item.Name
	itemResponse.Description = item.Description
	itemResponse.Tags = item.Tags
	itemResponse.ModifiedAt = item.ModifiedAt
	itemResponse.Version = item.Version

	return &itemResponse
}

//GetItemByTagResponse Get item by tag response payload
type GetItemByTagResponse struct {
	Items []*GetItemByIDResponse `json:"items"`
}

//NewGetItemByTagResponse construct GetItemByTagResponse
func NewGetItemByTagResponse(items []*item.Item) *GetItemByTagResponse {
	var itemResponses []*GetItemByIDResponse
	itemResponses = make([]*GetItemByIDResponse, 0)

	for _, item := range items {
		itemResponses = append(itemResponses, NewGetItemByIDResponse(item))
	}

	return &GetItemByTagResponse{
		itemResponses,
	}
}

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
