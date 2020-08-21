package item

import (
	"errors"
	"time"
)

var (
	//ErrConflict Error when update item that has been modified
	ErrConflict = errors.New("Data has been modified")

	//ErrNotFound Error when item is not found
	ErrNotFound = errors.New("Data was not found")

	//ErrBadRequest Error when data given is not valid on update or insert
	ErrBadRequest = errors.New("Given data is not valid")

	//ErrFailedToCast Error when failed to casting data
	ErrFailedToCast = errors.New("Failed to cast data")

	//ErrZeroAffected Data not found
	ErrZeroAffected = errors.New("No record affected")
)

//Item product item that available to rent or sell
type Item struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Tags        []string  `bson:"tags"`
	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   string    `bson:"created_by"`
	ModifiedAt  time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
	Version     int       `bson:"version"`
}

//DataRepository ingoing port for item
type DataRepository interface {
	//FindItemByID If data not found will return nil without error
	FindItemByID(ID string) (*Item, error)

	//FindAllByTag If no data match with the given tag, will return empty slice instead of nil
	FindAllByTag(tag string) ([]*Item, error)

	InsertItem(item Item) error

	//UpdateItem if data not found will return ErrZeroAffected
	UpdateItem(item Item, currentVersion int) error
}

//Service outgoing port for item
type Service interface {
	GetItemByID(ID string) (*Item, error)

	GetItemsByTag(tag string) ([]*Item, error)

	CreateItem(upsertitemSpec UpsertItemSpec, createdBy string) (string, error)

	UpdateItem(ID string, upsertitemSpec UpsertItemSpec, currentVersion int, modifiedBy string) error
}

//UpsertItemSpec create and update item spec
type UpsertItemSpec struct {
	Name        string   `validate:"required"`
	Description string   `validate:"required,min=3"`
	Tags        []string `validate:"required,min=0"`
}
