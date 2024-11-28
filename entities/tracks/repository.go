package tracks

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/xineman/go-server/db"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Track struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	FileName  string    `json:"fileName"`
	Photos    []string  `json:"photos"`
}

type TrackMetadata struct {
	name     string
	fileName string
}

func getTracks() ([]Track, error) {
	sql, _, _ := sq.Select(
		"tracks.id, tracks.created_at, updated_at, name, tracks.file_name, array_remove(array_agg(photos.file_name), null) as photo_file_name",
	).From("tracks").LeftJoin("photos on photos.track_id = tracks.id").GroupBy("tracks.id").ToSql()
	rows, err := db.DbPool.Query(context.Background(), sql)
	if err != nil {
		fmt.Println(err, sql)
		return nil, err
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		if err := rows.Scan(&track.Id, &track.CreatedAt, &track.UpdatedAt, &track.Name, &track.FileName, &track.Photos); err != nil {
			fmt.Println(err)
			return nil, err
		}

		tracks = append(tracks, track)
	}
	for _, t := range tracks {
		printTrack(t)
	}
	return tracks, nil
}

func saveTrack(track TrackMetadata) (Track, error) {
	sql, args, _ := psql.Insert("tracks").Columns("name", "file_name").Values(track.name, track.fileName).Suffix("returning *").ToSql()
	row := db.DbPool.QueryRow(context.Background(), sql, args...)
	var newTrack Track
	if err := row.Scan(&newTrack.Id, &newTrack.CreatedAt, &newTrack.UpdatedAt, &newTrack.Name, &newTrack.FileName); err != nil {
		return Track{}, err
	}
	return newTrack, nil
}

func printTrack(track Track) {
	fmt.Printf("Track: %s, created at %s\n", track.Name, track.CreatedAt.Format(time.RFC822))
}
