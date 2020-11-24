package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// BuildFileName returns the a date string in the format 20060102150405
func BuildFileName() string {
	return time.Now().Format("20060102150405")
}

// GetImageBytesFromURL returns the bytes of a http.Get request
func GetImageBytesFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body, err
}
