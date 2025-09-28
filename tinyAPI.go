package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type tinyAPI struct {
	srv       *http.Server
	db        *sql.DB
	serverDir string
}

type shortenPostRequest struct {
	URL string `json:"url"`
}

func NewServer(addr string) (*tinyAPI, error) {
	dir, err := configureDBLocation()
	if err != nil {
		return nil, err
	}

	db, err := createTestBD(dir)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Addr: ":8080",
	}

	return &tinyAPI{
		srv:       s,
		db:        db,
		serverDir: dir,
	}, nil
}

func (t *tinyAPI) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO WORLD!"))
}

func (t *tinyAPI) ShortenURL(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HEADERS: ", r.Header)
	var req shortenPostRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Could not read Body")
		return
	}

	json.Unmarshal(b, &req)
	if req.URL == "" {
		http.Error(w,
			fmt.Sprintf("No URL provided::%v", err),
			http.StatusBadRequest,
		)
		return
	}
	fmt.Println("BODY: ", req)

	hashedURL := CreateHashedURL(req.URL)
	addURLEntry(&hashedURL, t.db)
	data, err := json.Marshal(hashedURL)
	if err != nil {
		http.Error(w,
			"Could not marshall shorten response::%v",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (t *tinyAPI) ShutdownServer() error {
	log.Printf("\t- Removing dir: %v", t.serverDir)
	os.RemoveAll(t.serverDir)
	if err := t.db.Close(); err != nil {
		return err
	}
	return nil
}
