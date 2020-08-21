package item

import (
	"time"

	itemDomain "sample-order/domain/item"

	validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//MainService Implementation of service interface
type MainService struct {
	repository itemDomain.Repository
	validate   *validator.Validate
}

//NewServiceImpl Construct item service object
func NewServiceImpl(repository itemDomain.Repository) *MainService {
	return &MainService{
		repository,
		validator.New(),
	}
}

//GetItemByID Get item by given ID, return nil if not exist
func (s *MainService) GetItemByID(ID string) (*itemDomain.Item, error) {
	return s.repository.FindItemByID(ID)
}

//GetItemsByTag Get all items by given tag, return zero array if not match
func (s *MainService) GetItemsByTag(tag string) ([]*itemDomain.Item, error) {

	items, err := s.repository.FindAllByTag(tag)
	if err != nil || items == nil {
		return []*itemDomain.Item{}, err
	}

	return items, err
}

//CreateItem Create new item and store into database
func (s *MainService) CreateItem(upsertitemSpec itemDomain.UpsertItemSpec, createdBy string) (string, error) {
	err := s.validate.Struct(upsertitemSpec)

	if err != nil {
		return "", itemDomain.ErrBadRequest
	}

	ID := primitive.NewObjectID().Hex()

	now := time.Now()
	item := itemDomain.Item{
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

	err = s.repository.InsertItem(item)
	if err != nil {
		return "", err
	}

	return ID, nil
}

//UpdateItem Update existing item in the database.
//Will return ErrNotFound when item is not exists or ErrConflict if data version is not match
func (s *MainService) UpdateItem(ID string, upsertitemSpec itemDomain.UpsertItemSpec, currentVersion int, modifiedBy string) error {
	err := s.validate.Struct(upsertitemSpec)

	if err != nil || len(ID) == 0 {
		return itemDomain.ErrBadRequest
	}

	//get the item first to make sure data is exist
	item, err := s.repository.FindItemByID(ID)

	if err != nil {
		return err
	} else if item == nil {
		return itemDomain.ErrNotFound
	} else if item.Version != currentVersion {
		return itemDomain.ErrConflict
	}

	item.Name = upsertitemSpec.Name
	item.Description = upsertitemSpec.Description
	item.Tags = upsertitemSpec.Tags
	item.ModifiedBy = modifiedBy
	item.ModifiedAt = time.Now()
	item.Version = item.Version + 1

	return s.repository.UpdateItem(*item, currentVersion)
}