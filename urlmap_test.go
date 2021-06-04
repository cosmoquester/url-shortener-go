package main

import (
	"testing"
)

func TestURLMap(t *testing.T) {
	urlmap := URLMap{
		shortToLong:    make(map[string]string),
		longToShort:    make(map[string]string),
		shortURLLength: 10,
	}

	const originalLongURL = "helloworld"
	var shortURL string
	var longURL string
	var ok bool
	if shortURL, ok = urlmap.putURL(originalLongURL); !ok {
		t.Error("URLMap doesn't work sucessfully!")
	}
	if len(shortURL) != 10 {
		t.Error("ShortURL length is not same!")
	}

	if _, ok = urlmap.putURL(originalLongURL); ok {
		t.Error("Duplicate long URL not handled!")
	}
	if longURL, ok = urlmap.getLongURL(shortURL); !ok {
		t.Error("Get LongURL doesn't working!")
	}
	if longURL != originalLongURL {
		t.Error("Original retrieved longurl is diffrent from original url!")
	}

	if urlmap.delURL("none_short_url") {
		t.Error("Deleting non-existing short url doesn't handled!")
	}
	if !urlmap.delURL(shortURL) {
		t.Error("Error occured in deleting short url!")
	}
	if _, ok = urlmap.getLongURL(shortURL); ok {
		t.Error("Deleting non-existing short url didn't work!")
	}
}
