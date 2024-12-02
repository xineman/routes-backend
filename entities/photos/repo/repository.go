package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/xineman/go-server/db"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func SavePhotos(trackId int, photos []string) error {
	insertBuilder := psql.Insert("photos").Columns("file_name", "track_id")
	for _, photo := range photos {
		insertBuilder = insertBuilder.Values(photo, trackId)
	}
	sql, args, _ := insertBuilder.ToSql()
	_, err := db.DbPool.Exec(context.Background(), sql, args...)
	return err
}
