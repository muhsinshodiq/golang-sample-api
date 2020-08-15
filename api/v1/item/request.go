package item

import (
	"sample-order/domain/item"
)

//CreateItemRequest create item request payload
type CreateItemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

//ToUpsertItemSpec convert into item.UpsertItemSpec object
func (req *CreateItemRequest) ToUpsertItemSpec() *item.UpsertItemSpec {
	var upsertItemSpec item.UpsertItemSpec
	upsertItemSpec.Name = req.Name
	upsertItemSpec.Description = req.Description
	upsertItemSpec.Tags = req.Tags

	return &upsertItemSpec
}

//UpdateItemRequest update item request payload
type UpdateItemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Version     int      `json:"version" validate:"required"`
}

//ToUpsertItemSpec convert into item.UpsertItemSpec object
func (req *UpdateItemRequest) ToUpsertItemSpec() *item.UpsertItemSpec {
	var upsertItemSpec item.UpsertItemSpec
	upsertItemSpec.Name = req.Name
	upsertItemSpec.Description = req.Description
	upsertItemSpec.Tags = req.Tags

	return &upsertItemSpec
}
