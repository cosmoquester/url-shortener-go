package urlconverter

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileURLMap(t *testing.T) {
	var filepath string
	if file, err := ioutil.TempFile("", "tmp"); err != nil {
		t.Error("cannot use tmp file!")
	} else {
		filepath = file.Name()
	}
	defer os.Remove(filepath)

	urlmap, _ := NewFileURLMap(filepath, 10)

	const originalLongURL = "helloworld"
	var shortURL string
	var longURL string
	var ok bool
	if shortURL, ok = urlmap.PutURL(originalLongURL); !ok {
		t.Error("URLMap doesn't work sucessfully!")
	}
	if len(shortURL) != 10 {
		t.Error("ShortURL length is not same!")
	}

	if _, ok = urlmap.PutURL(originalLongURL); ok {
		t.Error("Duplicate long URL not handled!")
	}

	if longURL, ok = urlmap.GetLongURL(shortURL); !ok {
		t.Error("Get LongURL doesn't working!")
	}
	if longURL != originalLongURL {
		t.Error("Original retrieved longurl is diffrent from original url!")
	}

	if urlmap.DelURL("none_short_url") {
		t.Error("Deleting non-existing short url doesn't handled!")
	}
	if !urlmap.DelURL(shortURL) {
		t.Error("Error occured in deleting short url!")
	}
	if _, ok = urlmap.GetLongURL(shortURL); ok {
		t.Error("Deleting non-existing short url didn't work!")
	}
}
