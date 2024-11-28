package main

import (
	"context"
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
)

const DB_URL = "postgres://postgres:postgres@localhost/postgres"

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Track struct {
	id        int
	createdAt time.Time
	updatedAt time.Time
	name      string
	fileName  string
}

var conn *pgx.Conn

func CliInterface() {
	fmt.Println("Starting the app...")
	var err error
	conn, err = pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	for {
		fmt.Print("Choose an option:\n1: Get tracks\n2: Insert track\n3: Get tracks count\n")
		option := ""
		fmt.Scanln(&option)
		if option == "1" {
			tracks := getTracks()
			for _, track := range tracks {
				fmt.Printf("Track: %s, created at %s\n", track.name, track.createdAt.Format(time.RFC822))
			}
		} else if option == "2" {
			fmt.Print("Choose file name: ")
			var fileName string
			fmt.Scanln(&fileName)
			addTrack(TrackMetadata{capitalize(strings.ReplaceAll(fileName, "_", " ")), fileName})
		} else if option == "3" {
			printCount()
		}
	}
}

type TrackMetadata struct {
	name     string
	fileName string
}

func capitalize(word string) string {
	return fmt.Sprint(strings.ToUpper(string(word[0])), word[1:])
}

func printCount() {
	selectCount, _, _ := psql.Select("count(*)").From("tracks").ToSql()
	var count int
	conn.QueryRow(context.Background(), selectCount).Scan(&count)
	fmt.Printf("Count: %d\n", count)
}

func addTrack(track TrackMetadata) {
	sql, args, _ := psql.Insert("tracks").Columns("name", "file_name").Values(track.name, track.fileName).ToSql()
	fmt.Println(sql, args)
	_, err := conn.Exec(context.Background(), sql, args...)
	if err != nil {
		log.Fatal(err)
	}
}

func getTracks() []Track {
	sql, _, _ := psql.Select("*").From("tracks").ToSql()
	rows, err := conn.Query(context.Background(), sql)
	if err != nil {
		log.Fatal(err)
	}
	var tracks []Track
	for rows.Next() {
		var track Track
		if err := rows.Scan(&track.id, &track.createdAt, &track.updatedAt, &track.name, &track.fileName); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		tracks = append(tracks, track)
	}
	return tracks
}
