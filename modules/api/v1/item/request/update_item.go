package request

import "sample-order/business/item/spec"

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
