package item

import (
	"context"
	"sample-order/core/item"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//MongoDBRepository The implementation of item.Repository object
type MongoDBRepository struct {
	col *mongo.Collection
}

type collection struct {
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

func newCollection(item item.Item) *collection {
	return &collection{
		item.ID,
		item.Name,
		item.Description,
		item.Tags,
		item.CreatedAt,
		item.CreatedBy,
		item.ModifiedAt,
		item.ModifiedBy,
		item.Version,
	}
}

func (col *collection) ToItem() *item.Item {
	var item item.Item
	item.ID = col.ID
	item.Name = col.Name
	item.Description = col.Description
	item.Tags = col.Tags
	item.CreatedAt = col.CreatedAt
	item.CreatedBy = col.CreatedBy
	item.ModifiedAt = col.ModifiedAt
	item.ModifiedBy = col.ModifiedBy
	item.Version = col.Version

	return &item
}

//NewMongoDBRepository Generate mongo DB item repository
func NewMongoDBRepository(db *mongo.Database) *MongoDBRepository {
	return &MongoDBRepository{
		db.Collection("items"),
	}
}

//FindItemByID Find item based on given ID. Its return nil if not found
func (repo *MongoDBRepository) FindItemByID(ID string) (*item.Item, error) {
	var col collection

	filter := bson.M{
		"_id": ID,
	}

	if err := repo.col.FindOne(context.TODO(), filter).Decode(&col); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return col.ToItem(), nil
}

//FindAllByTag Find all items based on given tag. Its return empty array if not found
func (repo *MongoDBRepository) FindAllByTag(tag string) ([]*item.Item, error) {
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

	var items []*item.Item

	for cursor.Next(context.TODO()) {
		var col collection
		if err = cursor.Decode(&col); err != nil {
			return nil, err
		}

		items = append(items, col.ToItem())
	}

	return items, nil
}

//InsertItem Insert new item into database. Its return item id if success
func (repo *MongoDBRepository) InsertItem(item item.Item) error {
	col := newCollection(item)
	_, err := repo.col.InsertOne(context.TODO(), col)

	if err != nil {
		return err
	}

	return nil
}

//UpdateItem Update existing item in database
func (repo *MongoDBRepository) UpdateItem(item item.Item, currentVersion int) error {
	col := newCollection(item)

	filter := bson.M{
		"_id":     col.ID,
		"version": currentVersion,
	}

	updated := bson.M{
		"$set": col,
	}

	_, err := repo.col.UpdateOne(context.TODO(), filter, updated)
	if err != nil {
		return err
	}

	return nil
}
