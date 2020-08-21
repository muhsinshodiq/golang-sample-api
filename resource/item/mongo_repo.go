package item

import (
	"context"

	itemCore "sample-order/core/item"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//MongoDBRepository The implementation of item.Repository object
type MongoDBRepository struct {
	col *mongo.Collection
}

//NewMongoDBRepository Generate mongo DB item repository
func NewMongoDBRepository(db *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		db.Collection("items"),
	}
}

//FindItemByID Find item based on given ID. Its return nil if not found
func (repo *MongoDBRepository) FindItemByID(ID string) (*itemCore.Item, error) {
	var item *itemCore.Item

	filter := bson.M{
		"_id": ID,
	}

	if err := repo.col.FindOne(context.TODO(), filter).Decode(&item); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return item, nil
}

//FindAllByTag Find all items based on given tag. Its return empty array if not found
func (repo *MongoDBRepository) FindAllByTag(tag string) ([]*itemCore.Item, error) {
	filter := bson.M{
		"tags": bson.M{
			"$all": [1]string{tag},
		},
	}

	cursor, err := repo.col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var items []*itemCore.Item

	for cursor.Next(context.TODO()) {
		var item itemCore.Item
		if err = cursor.Decode(&item); err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

//InsertItem Insert new item into database. Its return item id if success
func (repo *MongoDBRepository) InsertItem(item itemCore.Item) error {
	_, err := repo.col.InsertOne(context.TODO(), item)

	if err != nil {
		return err
	}

	return nil
}

//UpdateItem Update existing item in database
func (repo *MongoDBRepository) UpdateItem(item itemCore.Item, currentVersion int) error {
	filter := bson.M{
		"_id":     item.ID,
		"version": currentVersion,
	}

	updated := bson.M{
		"$set": item,
	}

	_, err := repo.col.UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		return err
	}

	return nil
}
