package item_test

import (
	"os"
	"reflect"
	itemBusiness "sample-order/business/item"
	"sample-order/business/item/spec"
	"sample-order/core/item"
	core "sample-order/core/item"
	"testing"
	"time"
)

var service itemBusiness.Service
var item1, item2 core.Item
var insertSpec, updateSpec spec.UpsertItemSpec
var creator, updater string

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func TestGetItemByID(t *testing.T) {
	notFoundItem, _ := service.GetItemByID("random")

	if notFoundItem != nil {
		t.Error("Expect item must be not found (nil)")
	}

	foundItem, _ := service.GetItemByID(item1.ID)
	if !reflect.DeepEqual(*foundItem, item1) {
		t.Error("Expect item has to be equal with item1", foundItem, item1)
	}
}

func TestGetItemByTags(t *testing.T) {
	items, _ := service.GetItemsByTag("tag2")

	if len(items) != 2 {
		t.Error("Expect item length must be two")
		t.FailNow()
	}

	if reflect.DeepEqual(items[0], item1) {
		if !reflect.DeepEqual(items[1], item2) {
			t.Error("Expect 2nd item is equal to item2")
		}
	} else if reflect.DeepEqual(items[0], item2) {
		if !reflect.DeepEqual(items[1], item1) {
			t.Error("Expect 2nd item is equal to item1")
		}
	} else {
		t.Error("Expect items is item1 and item2")
	}
}

func TestCreateItem(t *testing.T) {
	id, _ := service.CreateItem(insertSpec, creator)

	for _, tag := range insertSpec.Tags {
		items, _ := service.GetItemsByTag(tag)

		if len(items) == 0 {
			t.Error("Expect at least one item when search by given tag: ", tag)
			continue
		}

		isFound := false
		for _, item := range items {
			if item.ID == id {
				isFound = true
				break
			}
		}

		if !isFound {
			t.Error("Expect found inserted item when search by given tag: ", tag)
		}
	}

	newItem, _ := service.GetItemByID(id)

	if newItem == nil {
		t.Error("Expect item is not nil after inserted")
		t.FailNow()
	}

	if newItem.Name != insertSpec.Name {
		t.Error("Expect name is equal as given")
	}

	if newItem.Description != insertSpec.Description {
		t.Error("Expect description is equal as given")
	}

	if !reflect.DeepEqual(newItem.Tags, insertSpec.Tags) {
		t.Error("Expect tags is equal as given")
	}

	if newItem.CreatedBy != creator {
		t.Error("Expect created by is equal to " + creator)
	}

	if newItem.ModifiedBy != creator {
		t.Error("Expect modified by is equal to " + creator)
	}

	if newItem.CreatedAt != newItem.ModifiedAt {
		t.Error("Expect created at and modified at is equal")
	}

	if newItem.Version != 1 {
		t.Error("Expect version is equal to 1")
	}
}

func TestUpdateItem(t *testing.T) {
	id := item2.ID
	version := item2.Version
	oldTags := item2.Tags

	service.UpdateItem(id, updateSpec, version, updater)

	//find the old tag that doesn't exist in new updated tags
	var invalidateTags []string
	for _, tag := range oldTags {
		isExistOnUpdatedTag := false

		for _, updatedTag := range updateSpec.Tags {
			if tag == updatedTag {
				isExistOnUpdatedTag = true
				break
			}
		}

		if !isExistOnUpdatedTag {
			invalidateTags = append(invalidateTags, tag)
		}
	}

	//verify the invalidated tag is not contain the item anymore
	for _, invalidateTag := range invalidateTags {
		tagItems, _ := service.GetItemsByTag(invalidateTag)
		isFound := false

		for _, tagItem := range tagItems {
			if tagItem.ID == id {
				isFound = true
				break
			}
		}

		if isFound {
			t.Error("Expect not found when search by old invalidate tag: ", invalidateTag)
		}
	}

	items, _ := service.GetItemsByTag(updateSpec.Tags[0])

	isFound := false
	for _, item := range items {
		if item.ID == id {
			isFound = true
			break
		}
	}

	if !isFound {
		t.Error("Expect found inserted item when search by given tag: ", updateSpec.Tags[0])
	}

	updatedItem, _ := service.GetItemByID(item2.ID)

	if updatedItem == nil {
		t.Error("Expect item is not nil after updated")
		t.FailNow()
	}

	if updatedItem.Name != updateSpec.Name {
		t.Error("Expect name is equal as given")
	}

	if updatedItem.Description != updateSpec.Description {
		t.Error("Expect description is equal as given")
	}

	if !reflect.DeepEqual(updatedItem.Tags, updateSpec.Tags) {
		t.Error("Expect tags is equal as given")
	}

	if updatedItem.CreatedBy != item2.CreatedBy {
		t.Error("Expect created by is equal to " + item2.CreatedBy)
	}

	if updatedItem.ModifiedBy != updater {
		t.Error("Expect modified by is equal to " + updater)
	}

	if updatedItem.CreatedAt == updatedItem.ModifiedAt {
		t.Error("Expect created at and modified at is not equal")
	}

	if updatedItem.Version != item2.Version+1 {
		t.Error("Expect version was increase by one")
	}
}

func setup() {
	//initialize item1
	item1.ID = "5f350b7d21148431abc65290"
	item1.Name = "Item one"
	item1.Description = "Description one"
	item1.Tags = []string{"tag1", "tag2"}
	item1.Version = 1
	item1.CreatedAt = time.Now()
	item1.CreatedBy = "creator one"
	item1.ModifiedAt = item1.CreatedAt
	item1.ModifiedBy = item1.CreatedBy

	//initialize item 2
	item2.ID = "5f351360ac84a3bb1baee057"
	item2.Name = "Item two"
	item2.Description = "Description two"
	item2.Tags = []string{"tag2", "tag3", "tag4"}
	item2.Version = 2
	item2.CreatedAt = time.Now().Add(time.Minute * -15)
	item2.CreatedBy = "creator two"
	item2.ModifiedAt = time.Now()
	item2.ModifiedBy = "updater two"

	repo := newInMemoryRepository()
	service = itemBusiness.NewService(&repo)

	insertSpec.Name = "New Item"
	insertSpec.Description = "New Description"
	insertSpec.Tags = []string{"tag99"}

	updateSpec.Name = "Update Item"
	updateSpec.Description = "Update Description"
	updateSpec.Tags = []string{"tag99-updated"}

	creator = "creator"
	updater = "updater"
}

type inMemoryRepository struct {
	itemByID  map[string]core.Item
	itemByTag map[string][]core.Item
}

func newInMemoryRepository() inMemoryRepository {
	var repo inMemoryRepository
	repo.itemByID = make(map[string]core.Item)
	repo.itemByTag = make(map[string][]core.Item)

	repo.itemByID[item1.ID] = item1
	repo.itemByID[item2.ID] = item2

	for _, tag := range item1.Tags {
		items := repo.itemByTag[tag]
		repo.itemByTag[tag] = append(items, item1)
	}

	for _, tag := range item2.Tags {
		items := repo.itemByTag[tag]
		repo.itemByTag[tag] = append(items, item2)
	}

	return repo
}

func (repo *inMemoryRepository) FindItemByID(ID string) (*core.Item, error) {
	item, ok := repo.itemByID[ID]
	if !ok {
		return nil, nil
	}

	return &item, nil
}

func (repo *inMemoryRepository) FindAllByTag(tag string) ([]item.Item, error) {
	items, ok := repo.itemByTag[tag]
	if !ok {
		return []core.Item{}, nil
	}

	return items, nil
}

func (repo *inMemoryRepository) InsertItem(item core.Item) error {
	repo.itemByID[item.ID] = item

	for _, tag := range item.Tags {
		items := repo.itemByTag[tag]
		repo.itemByTag[tag] = append(items, item)
	}
	return nil
}

func (repo *inMemoryRepository) UpdateItem(item core.Item, currentVersion int) error {
	//cleanup old tag first
	oldItem := repo.itemByID[item.ID]

	//cleanup the old tags first
	for _, tag := range oldItem.Tags {
		tagItems, _ := repo.FindAllByTag(tag)

		itemIndex := -1
		for idx, tagItem := range tagItems {
			if tagItem.ID == item.ID {
				itemIndex = idx
				break
			}
		}

		if itemIndex != -1 {
			tagItems = append(tagItems[:itemIndex], tagItems[itemIndex+1:]...)
		}

		repo.itemByTag[tag] = tagItems
	}

	repo.itemByID[item.ID] = item

	//adding the new tag
	for _, tag := range item.Tags {
		items := repo.itemByTag[tag]
		repo.itemByTag[tag] = append(items, item)
	}
	return nil
}
