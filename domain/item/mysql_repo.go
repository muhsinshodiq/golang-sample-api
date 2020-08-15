package item

import (
	"database/sql"
	"strings"
)

//MySQLRepository The implementation of item.Repository object
type MySQLRepository struct {
	db *sql.DB
}

//NewMySQLRepository Generate mongo DB item repository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db,
	}
}

//FindItemByID Find item based on given ID. Its return nil if not found
func (repo *MySQLRepository) FindItemByID(ID string) (*Item, error) {
	var item Item

	selectQuery := `SELECT id, name, description, created_at, created_by, modified_at, modified_by, version, COALESCE(tags, "")
		FROM item i
		LEFT JOIN (
			SELECT item_id, 
			GROUP_CONCAT(tag) as tags
			FROM item_tag GROUP BY item_id
		)AS it ON i.id = it.item_id
		WHERE i.id = ?`

	var tags string
	err := repo.db.
		QueryRow(selectQuery, ID).
		Scan(
			&item.ID, &item.Name, &item.Description,
			&item.CreatedAt, &item.CreatedBy,
			&item.ModifiedAt, &item.ModifiedBy,
			&item.Version, &tags)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	item.Tags = constructTagArray(tags)

	return &item, nil
}

//FindAllByTag Find all items based on given tag. Its return empty array if not found
func (repo *MySQLRepository) FindAllByTag(tag string) ([]*Item, error) {
	var items []*Item

	//TODO: if feel have a performance issue in tag grouping, move the logic from db to here
	selectQuery := `SELECT id, name, description, created_at, created_by, modified_at, modified_by, version, COALESCE(tags, "")
		FROM item i
		LEFT JOIN (
			SELECT item_id, 
			GROUP_CONCAT(tag) as tags
			FROM item_tag GROUP BY item_id
		)AS it ON i.id = it.item_id
		WHERE i.id IN (
			SELECT item_id
			FROM item_tag
			WHERE tag = ?	
		)`

	row, err := repo.db.Query(selectQuery, tag)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		var item Item
		var tags string

		err := row.Scan(
			&item.ID, &item.Name, &item.Description,
			&item.CreatedAt, &item.CreatedBy,
			&item.ModifiedAt, &item.ModifiedBy,
			&item.Version, &tags)

		if err != nil {
			return nil, err
		}

		item.Tags = constructTagArray(tags)
		items = append(items, &item)
	}

	if err != nil {
		return nil, err
	}

	return items, nil
}

//InsertItem Insert new item into database. Its return item id if success
func (repo *MySQLRepository) InsertItem(item Item) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	itemQuery := `INSERT INTO item (
			id, 
			name, 
			description, 
			created_at, 
			created_by, 
			modified_at, 
			modified_by,
			version
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err != nil {
		return err
	}

	_, err = tx.Exec(itemQuery,
		item.ID,
		item.Name,
		item.Description,
		item.CreatedAt,
		item.CreatedBy,
		item.ModifiedAt,
		item.ModifiedBy,
		item.Version,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	tagQuery := "INSERT INTO item_tag (item_id, tag) VALUES (?, ?)"

	for _, tag := range item.Tags {
		_, err = tx.Exec(tagQuery, item.ID, tag)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

//UpdateItem Update existing item in database
func (repo *MySQLRepository) UpdateItem(item Item, currentVersion int) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}

	itemInsertQuery := `UPDATE item 
		SET
			name = ?,
			description = ?,
			modified_at = ?,
			modified_by = ?,
			version = ?
		WHERE id = ? AND version = ?`

	res, err := tx.Exec(itemInsertQuery,
		item.Name,
		item.Description,
		item.ModifiedAt,
		item.ModifiedBy,
		item.Version,
		item.ID,
		currentVersion,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	affected, err := res.RowsAffected()

	if err != nil {
		tx.Rollback()
		return err
	}

	if affected == 0 {
		tx.Rollback()
		return ErrZeroAffected
	}

	//TODO: maybe better if we only delete the record that we need to delete
	//add logic slice to find which deleted and which want to added
	tagDeleteQuery := "DELETE FROM item_tag WHERE item_id = ?"
	_, err = tx.Exec(tagDeleteQuery, item.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	tagUpsertQuery := "INSERT INTO item_tag (item_id, tag) VALUES (?, ?)"

	for _, tag := range item.Tags {
		_, err = tx.Exec(tagUpsertQuery, item.ID, tag)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func constructTagArray(tags string) []string {
	if tags == "" {
		return make([]string, 0)
	}

	return strings.Split(tags, ",")
}
