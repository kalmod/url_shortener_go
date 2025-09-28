package main

import (
	"context"
	"database/sql"
)

// This is a test function used to insert a bunch of data internally
func insertDataIntoTable(db *sql.DB, longURLs ...string) error {
	for _, longURL := range longURLs {
		key := urlHash(longURL)
		resp := shortenedResponse(key, longURL)
		_, err := addURLEntry(&resp, db)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add a new entry to shortFullURLMap
func addURLEntry(url *URLMapping, db *sql.DB) (sql.Result, error) {
	result, err := db.ExecContext(
		context.Background(),
		`INSERT INTO shortFullURLMap (key, longURL, shortURL) VALUES (?, ?, ?);`, url.Key, url.LongURL, url.ShortURL,
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Retrieve an entry from shortFullURLMap using KEY
func getURLEntry(key string, ctx context.Context, db *sql.DB) (URLMapping, error) {
	var urlData URLMapping
	row := db.QueryRowContext(ctx, `SELECT * FROM shortFullURLMap WHERE key=?;`, key)
	err := row.Scan(&urlData.Key, &urlData.LongURL, &urlData.ShortURL)
	if err != nil {
		return urlData, err
	}

	return urlData, nil
}

func deleteURLEntry(key string, ctx context.Context, db *sql.DB) (sql.Result, error) {
	res, err := db.ExecContext(ctx, "DELETE FROM shortFullURLMap WHERE key=?", key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
