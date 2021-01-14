package item

import (
	"sample-order/business"
	"sample-order/business/item/spec"
	"sample-order/util"
	"time"

	validator "github.com/go-playground/validator/v10"
)

//Repository ingoing port for item
type Repository interface {
	//FindItemByID If data not found will return nil without error
	FindItemByID(ID string) (*Item, error)

	//FindAllByTag If no data match with the given tag, will return empty slice instead of nil
	FindAllByTag(tag string) ([]Item, error)

	//InsertItem Insert new item into storage
	InsertItem(item Item) error

	//UpdateItem if data not found will return core.ErrZeroAffected
	UpdateItem(item Item, currentVersion int) error
}

//Service outgoing port for item
type Service interface {
	GetItemByID(ID string) (*Item, error)

	GetItemsByTag(tag string) ([]Item, error)

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
func (s *service) GetItemByID(ID string) (*Item, error) {
	return s.repository.FindItemByID(ID)
}

//GetItemsByTag Get all items by given tag, return zero array if not match
func (s *service) GetItemsByTag(tag string) ([]Item, error) {

	items, err := s.repository.FindAllByTag(tag)
	if err != nil || items == nil {
		return []Item{}, err
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
	item := NewItem(
		ID,
		upsertitemSpec.Name,
		upsertitemSpec.Description,
		upsertitemSpec.Tags,
		createdBy,
		time.Now(),
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

	newItem := item.ModifyItem(upsertitemSpec.Name, upsertitemSpec.Description, upsertitemSpec.Tags, modifiedBy, time.Now())

	return s.repository.UpdateItem(newItem, currentVersion)
}
