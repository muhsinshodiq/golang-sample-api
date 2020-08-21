package item

//DataRepository ingoing port for item
type DataRepository interface {
	//FindItemByID If data not found will return nil without error
	FindItemByID(ID string) (*Item, error)

	//FindAllByTag If no data match with the given tag, will return empty slice instead of nil
	FindAllByTag(tag string) ([]*Item, error)

	//InsertItem Insert new item into storage
	InsertItem(item Item) error

	//UpdateItem if data not found will return core.ErrZeroAffected
	UpdateItem(item Item, currentVersion int) error
}
