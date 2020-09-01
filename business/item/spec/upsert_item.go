package spec

//UpsertItemSpec create and update item spec
type UpsertItemSpec struct {
	Name        string   `validate:"required"`
	Description string   `validate:"required,min=3"`
	Tags        []string `validate:"required"`
}
