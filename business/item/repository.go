package item

import core "sample-order/core/item"

//Repository ingoing port for item
type Repository interface {
	//FindItemByID If data not found will return nil without error
	FindItemByID(ID string) (*core.Item, error)

	//FindAllByTag If no data match with the given tag, will return empty slice instead of nil
	FindAllByTag(tag string) ([]core.Item, error)

	//InsertItem Insert new item into storage
	InsertItem(item core.Item) error

	//UpdateItem if data not found will return core.ErrZeroAffected
	UpdateItem(item core.Item, currentVersion int) error
}
