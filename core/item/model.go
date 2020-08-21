package item

import "time"

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
