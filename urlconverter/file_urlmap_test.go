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
		t.Error("URLMap doesn't work successfully!")
		return
	}
	if len(shortURL) != 10 {
		t.Error("ShortURL length is not same!")
		return
	}

	if _, ok = urlmap.PutURL(originalLongURL); ok {
		t.Error("Duplicate long URL not handled!")
		return
	}

	if longURL, ok = urlmap.GetLongURL(shortURL); !ok {
		t.Error("Get LongURL doesn't working!")
		return
	}
	if longURL != originalLongURL {
		t.Error("Original retrieved longurl is different from original url!")
		return
	}

	if urlmap.DelURL("none_short_url") {
		t.Error("Deleting non-existing short url doesn't handled!")
		return
	}
	if !urlmap.DelURL(shortURL) {
		t.Error("Error occurred in deleting short url!")
		return
	}
	if _, ok = urlmap.GetLongURL(shortURL); ok {
		t.Error("Deleting non-existing short url didn't work!")
		return
	}
}
