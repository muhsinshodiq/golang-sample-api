package response

import (
	"sample-order/business/item"
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
func NewGetItemByIDResponse(item item.Item) *GetItemByIDResponse {
	var itemResponse GetItemByIDResponse
	itemResponse.ID = item.ID
	itemResponse.Name = item.Name
	itemResponse.Description = item.Description
	itemResponse.Tags = item.Tags
	itemResponse.ModifiedAt = item.ModifiedAt
	itemResponse.Version = item.Version

	return &itemResponse
}
