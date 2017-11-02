package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ValidURL converts the string into a url.URL struct, but returns an error on fail
// this is taken advantage of to validate the url syntax
// struct https://golangcode.com/how-to-check-if-a-string-is-a-url/
func ValidURL(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	return true
}

// GetJSON does a GET request and parses the response body into the target interface
// https://stackoverflow.com/questions/15672556/handling-json-post-request-in-go
func GetJSON(url string, target interface{}) error {
	myClient := &http.Client{Timeout: 5 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("get json issue")
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
