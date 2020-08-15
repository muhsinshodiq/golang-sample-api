package item

import (
	"errors"
	"time"
)

var (
	//ErrConflict Error when update item that has been modified
	ErrConflict = errors.New("Item has been modified")

	//ErrNotFound Error when item is not found
	ErrNotFound = errors.New("Item was not found")

	//ErrBadRequest Error when data given is not valid on update or insert
	ErrBadRequest = errors.New("Given data is not valid")

	//ErrFailedToCast Error when failed to casting data
	ErrFailedToCast = errors.New("Failed to cast data")

	//ErrZeroAffected Data not found
	ErrZeroAffected = errors.New("No record affected")
)

//Item product item that available to rent or sell
type Item struct {
	ID          string    `db:"id" bson:"_id"`
	Name        string    `db:"name" bson:"name"`
	Description string    `db:"description" bson:"description"`
	Tags        []string  `db:"tags"`
	CreatedAt   time.Time `db:"created_at" bson:"created_at"`
	CreatedBy   string    `db:"created_by" bson:"created_by"`
	ModifiedAt  time.Time `db:"modified_at" bson:"modified_at"`
	ModifiedBy  string    `db:"modified_by" bson:"modified_by"`
	Version     int       `db:"version" bson:"version"`
}

//Repository ingoing port for item
type Repository interface {
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
