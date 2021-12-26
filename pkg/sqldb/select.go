package sqldb

import "context"

type Item struct {
	Id      string
	Content string
}

// SelectEntryByRowNumber selects an entry by row number
func (db *DB) SelectEntryByRowNumber(model interface{}, rowNumber int) (item Item, err error) {

	ctx := context.Background()

	var selectedItem = Item{}

	err = db.NewSelect().
		Model(model).
		Column("id", "content").
		Limit(1).
		Offset(rowNumber-1).
		Scan(ctx, &selectedItem)

	if err != nil {
		return selectedItem, err
	}

	return selectedItem, nil

}

// SelectEntryById selects an entry by its id
func (db *DB) SelectEntryById(model interface{}, id string) (item Item, err error) {

	ctx := context.Background()

	var selectedItem = Item{}

	err = db.NewSelect().
		Model(model).
		Column("id", "content").
		Where("id = ?", id).
		Scan(ctx, &selectedItem)

	if err != nil {
		return selectedItem, err
	}

	return selectedItem, nil

}
