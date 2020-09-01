package item

import (
	"sample-order/business"
	"sample-order/business/item/spec"
	core "sample-order/core/item"
	"sample-order/util"

	validator "github.com/go-playground/validator/v10"
)

//Service outgoing port for item
type Service interface {
	GetItemByID(ID string) (*core.Item, error)

	GetItemsByTag(tag string) ([]core.Item, error)

	CreateItem(upsertitemSpec spec.UpsertItemSpec, createdBy string) (string, error)

	UpdateItem(ID string, upsertitemSpec spec.UpsertItemSpec, currentVersion int, modifiedBy string) error
}

//=============== The implementation of those interface put below =======================

type service struct {
	repository Repository
	validate   *validator.Validate
}

//NewService Construct item service object
func NewService(repository Repository) Service {
	return &service{
		repository,
		validator.New(),
	}
}

//GetItemByID Get item by given ID, return nil if not exist
func (s *service) GetItemByID(ID string) (*core.Item, error) {
	return s.repository.FindItemByID(ID)
}

//GetItemsByTag Get all items by given tag, return zero array if not match
func (s *service) GetItemsByTag(tag string) ([]core.Item, error) {

	items, err := s.repository.FindAllByTag(tag)
	if err != nil || items == nil {
		return []core.Item{}, err
	}

	return items, err
}

//CreateItem Create new item and store into database
func (s *service) CreateItem(upsertitemSpec spec.UpsertItemSpec, createdBy string) (string, error) {
	err := s.validate.Struct(upsertitemSpec)

	if err != nil {
		return "", business.ErrInvalidSpec
	}

	ID := util.GenerateID()
	item := core.CreateItem(
		ID,
		upsertitemSpec.Name,
		upsertitemSpec.Description,
		upsertitemSpec.Tags,
		createdBy,
	)

	err = s.repository.InsertItem(item)
	if err != nil {
		return "", err
	}

	return ID, nil
}

//UpdateItem Update existing item in the database.
//Will return ErrNotFound when item is not exists or ErrConflict if data version is not match
func (s *service) UpdateItem(ID string, upsertitemSpec spec.UpsertItemSpec, currentVersion int, modifiedBy string) error {
	err := s.validate.Struct(upsertitemSpec)

	if err != nil || len(ID) == 0 {
		return business.ErrInvalidSpec
	}

	//get the item first to make sure data is exist
	item, err := s.repository.FindItemByID(ID)

	if err != nil {
		return err
	} else if item == nil {
		return business.ErrNotFound
	} else if item.Version != currentVersion {
		return business.ErrHasBeenModified
	}

	newItem := core.ModifyItem(*item, upsertitemSpec.Name, upsertitemSpec.Description, upsertitemSpec.Tags, modifiedBy)

	return s.repository.UpdateItem(newItem, currentVersion)
}
