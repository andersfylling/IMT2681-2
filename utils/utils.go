package utils

import "net/url"

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func ValidURL(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	return true
}
