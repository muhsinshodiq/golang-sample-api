package item

import (
	"sample-order/core"
	"sample-order/core/item/spec"
	"sample-order/libs"
	"time"

	validator "github.com/go-playground/validator/v10"
)

//Service outgoing port for item
type Service interface {
	GetItemByID(ID string) (*Item, error)

	GetItemsByTag(tag string) ([]*Item, error)

	CreateItem(upsertitemSpec spec.UpsertItemSpec, createdBy string) (string, error)

	UpdateItem(ID string, upsertitemSpec spec.UpsertItemSpec, currentVersion int, modifiedBy string) error
}

//=============== The implementation of those interface put below =======================

type service struct {
	dataRepository DataRepository
	validate       *validator.Validate
}

//NewService Construct item service object
func NewService(dataRepository DataRepository) Service {
	return &service{
		dataRepository,
		validator.New(),
	}
}

//GetItemByID Get item by given ID, return nil if not exist
func (s *service) GetItemByID(ID string) (*Item, error) {
	return s.dataRepository.FindItemByID(ID)
}

//GetItemsByTag Get all items by given tag, return zero array if not match
func (s *service) GetItemsByTag(tag string) ([]*Item, error) {

	items, err := s.dataRepository.FindAllByTag(tag)
	if err != nil || items == nil {
		return []*Item{}, err
	}

	return items, err
}

//CreateItem Create new item and store into database
func (s *service) CreateItem(upsertitemSpec spec.UpsertItemSpec, createdBy string) (string, error) {
	err := s.validate.Struct(upsertitemSpec)

	if err != nil {
		return "", core.ErrInvalidSpec
	}

	ID := libs.GenerateID()

	now := time.Now()
	item := Item{
		ID:          ID,
		Name:        upsertitemSpec.Name,
		Description: upsertitemSpec.Description,
		Tags:        upsertitemSpec.Tags,
		CreatedAt:   now,
		CreatedBy:   createdBy,
		ModifiedAt:  now,
		ModifiedBy:  createdBy,
		Version:     1,
	}

	err = s.dataRepository.InsertItem(item)
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
		return core.ErrInvalidSpec
	}

	//get the item first to make sure data is exist
	item, err := s.dataRepository.FindItemByID(ID)

	if err != nil {
		return err
	} else if item == nil {
		return core.ErrNotFound
	} else if item.Version != currentVersion {
		return core.ErrHasBeenModified
	}

	item.Name = upsertitemSpec.Name
	item.Description = upsertitemSpec.Description
	item.Tags = upsertitemSpec.Tags
	item.ModifiedBy = modifiedBy
	item.ModifiedAt = time.Now()
	item.Version = item.Version + 1

	return s.dataRepository.UpdateItem(*item, currentVersion)
}
