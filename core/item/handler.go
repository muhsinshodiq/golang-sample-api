package item

import "time"

//CreateItem create new item
func CreateItem(
	id string,
	name string,
	description string,
	tags []string,
	creator string) Item {

	timeNow := time.Now()

	return Item{
		ID:          id,
		Name:        name,
		Description: description,
		Tags:        tags,
		CreatedAt:   timeNow,
		CreatedBy:   creator,
		ModifiedAt:  timeNow,
		ModifiedBy:  creator,
		Version:     1,
	}
}

//ModifyItem update existing item data
func ModifyItem(
	oldItem Item,
	newName string,
	newDescription string,
	newTags []string,
	updater string) Item {

	timeNow := time.Now()

	return Item{
		ID:          oldItem.ID,
		Name:        newName,
		Description: newDescription,
		Tags:        newTags,
		CreatedAt:   oldItem.CreatedAt,
		CreatedBy:   oldItem.CreatedBy,
		ModifiedAt:  timeNow,
		ModifiedBy:  updater,
		Version:     oldItem.Version + 1,
	}
}
