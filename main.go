package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	testURL := "https://www.google.com"
	shortURL := urlHash(testURL)

	shrt := shortenedResponse(shortURL, testURL)

	json, err := jsonResponse(shrt)
	if err != nil {
		fmt.Println("Error marshalling")
		os.Exit(1)
	}

	fmt.Println(string(json))

	// if err := test(); err != nil {
	// 	fmt.Printf("error db: %s\n", err.Error())
	// 	os.Exit(1)
	// }

	ta, err := NewServer("8080")
	if err != nil {
		log.Fatalf("Error creating new server::%v", err)
	}
	err = ta.initTestDB()
	if err != nil {
		fmt.Println("error on init::", err.Error())
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", ta.Home)
	mux.HandleFunc("POST /", ta.ShortenURL)

	ta.srv.Handler = mux

	go func() {
		if err := ta.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	fmt.Println()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	ta.ShutdownServer()
	defer shutdownRelease()
	if err := ta.srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server shutdown complete.")
	os.Exit(0)
}
