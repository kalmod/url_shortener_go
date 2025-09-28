package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/url"

	_ "modernc.org/sqlite"
)

func jsonResponse(urlResponse URLMapping) ([]byte, error) {
	json, err := json.Marshal(urlResponse)
	if err != nil {
		return []byte{}, nil
	}
	return json, nil
}

func shortenedResponse(key, longURL string) URLMapping {
	return URLMapping{
		key,
		longURL,
		fmt.Sprintf("http://localhost/%s", key),
	}
}

// NOTE: Clean up CreateHashedURL

func CreateHashedURL(longURL string) URLMapping {
	key := urlHash(longURL)
	return shortenedResponse(key, longURL)
}

func urlHash(longURL string) string {
	h := fnv.New32a()
	h.Write([]byte(longURL))

	var buffer bytes.Buffer

	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)

	input := h.Sum(nil)
	encoder.Write(input)

	defer encoder.Close()

	// fmt.Printf("%x\n", h.Sum(nil))
	return url.QueryEscape(buffer.String())
}
