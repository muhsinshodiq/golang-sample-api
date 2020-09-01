package response

import "sample-order/core/item"

//GetItemByTagResponse Get item by tag response payload
type GetItemByTagResponse struct {
	Items []*GetItemByIDResponse `json:"items"`
}

//NewGetItemByTagResponse construct GetItemByTagResponse
func NewGetItemByTagResponse(items []item.Item) *GetItemByTagResponse {
	var itemResponses []*GetItemByIDResponse
	itemResponses = make([]*GetItemByIDResponse, 0)

	for _, item := range items {
		itemResponses = append(itemResponses, NewGetItemByIDResponse(item))
	}

	return &GetItemByTagResponse{
		itemResponses,
	}
}
