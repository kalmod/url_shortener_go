package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type URLMapping struct {
	Key      string `json:"key"`
	LongURL  string `json:"long_url"`
	ShortURL string `json:"short_url"`
}

func getDBLocationPath(locationPath string) (string, error) {
	// DB location will always be relative to the executable

	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	path = filepath.Dir(path)
	if strings.HasPrefix(path, "/tmp/") {
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	absolutePath, err := filepath.Abs(filepath.Join(path, locationPath))
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}

func configureDBLocation() (string, error) {
	// STARTING DB
	// 1. CREATING FILES

	dbLocation, err := getDBLocationPath("")
	if err != nil {
		return "", err
	}

	fmt.Println(dbLocation)

	dir, err := os.MkdirTemp(dbLocation+"/", "test-")
	if err != nil {
		return "", err
	}
	return dir, nil
}

func createTestBD(dir string) (*sql.DB, error) {
	fmt.Println("1. Starting DB test")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}
	// use os.Stat to check if not exists

	fn := filepath.Join(dir, "test.db")
	if _, err := os.Create(fn); err != nil {
		return nil, err
	}

	// 2. OPEN DB
	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return db, err
}

// TODO: Update to only create tables when ready
// TODO: Check to make sure the url contains https
func (t *tinyAPI) initTestDB() error {
	fmt.Println("Initiating DB...")
	err := createTestTable(context.Background(), t.db)
	if err != nil {
		return err
	}

	err = insertDataIntoTable(t.db,
		"https://www.google.com",
		"https://www.reddit.com/r/sffpc",
		"https://www.amazon.com",
		"https://www.youtube.com")
	if err != nil {
		return err
	}

	rows, err := t.db.Query("select * from shortFullURLMap;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var urlObj URLMapping
		if err = rows.Scan(&urlObj.Key, &urlObj.LongURL, &urlObj.ShortURL); err != nil {
			return err
		}

		fmt.Println("\t", urlObj)
	}
	return nil
}

// Creates a table for testing implementation.
// Update statement to -> Create Table if not exists when finished
func createTestTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		`drop table if exists shortFullURLMap;
		create table shortFullURLMap(
		key TEXT PRIMARY KEY,
		longURL TEXT NOT NULL,
		shortURL REAL NOT NULL); `)
	if err != nil {
		return err
	}
	return nil
}

// TODO: DELETE FROM FILES
// FOR TESTING
func test() error {
	fmt.Println("Starting DB test")

	// STARTING DB
	// 1. CREATING FILES

	dbLocation, err := getDBLocationPath("")
	if err != nil {
		return err
	}

	fmt.Println(dbLocation)

	dir, err := os.MkdirTemp(dbLocation+"/", "test-")
	// fmt.Println(dir)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}
	// use os.Stat to check if not exists

	fn := filepath.Join(dir, "test.db")
	_, err = os.Create(fn)
	if err != nil {
		return err
	}

	// fmt.Println(fn)

	// 2. OPEN DB
	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return err
	}

	err = createTestTable(context.Background(), db)
	if err != nil {
		return err
	}

	err = insertDataIntoTable(db,
		"www.google.com",
		"www.reddit.com/r/sffpc",
		"www.amazon.com",
		"www.youtube.com")
	if err != nil {
		return err
	}

	rows, err := db.Query("select * from shortFullURLMap;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var urlObj URLMapping
		if err = rows.Scan(&urlObj.Key, &urlObj.LongURL, &urlObj.ShortURL); err != nil {
			return err
		}

		fmt.Println("\t", urlObj)
	}

	// 3. CLOSING OUT

	if err = db.Close(); err != nil {
		return err
	}

	return nil
}
