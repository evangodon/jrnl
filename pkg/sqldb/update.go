package sqldb

import "context"

func (db *DB) UpdateEntryContent(model interface{}, item Item) (err error) {

	ctx := context.Background()

	_, err = db.NewUpdate().
		Model(model).
		Set("content = ?", item.Content).
		Where("id = ?", item.Id).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
