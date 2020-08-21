package item

import (
	"sample-order/core/item/spec"
)

//CreateItemRequest create item request payload
type CreateItemRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

//ToUpsertItemSpec convert into item.UpsertItemSpec object
func (req *CreateItemRequest) ToUpsertItemSpec() *spec.UpsertItemSpec {
	var upsertItemSpec spec.UpsertItemSpec
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
func (req *UpdateItemRequest) ToUpsertItemSpec() *spec.UpsertItemSpec {
	var upsertItemSpec spec.UpsertItemSpec
	upsertItemSpec.Name = req.Name
	upsertItemSpec.Description = req.Description
	upsertItemSpec.Tags = req.Tags

	return &upsertItemSpec
}
