package request

import "sample-order/business/item/spec"

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
