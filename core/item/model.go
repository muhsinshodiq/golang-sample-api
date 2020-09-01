package item

import "time"

//Item product item that available to rent or sell
type Item struct {
	ID          string
	Name        string
	Description string
	Tags        []string
	CreatedAt   time.Time
	CreatedBy   string
	ModifiedAt  time.Time
	ModifiedBy  string
	Version     int
}
