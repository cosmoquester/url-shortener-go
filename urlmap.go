package main

import (
	"sync"
)

// URLMap 은 LongURL와 ShortURL을 맵핑할 수 있는 map입니다.
type URLMap struct {
	shortToLong, longToShort map[string]string
	shortURLLength           uint
	sync.Mutex
}

func (urlmap *URLMap) getLongURL(shortURL string) (string, bool) {
	longURL, ok := urlmap.shortToLong[shortURL]
	return longURL, ok
}
func (urlmap *URLMap) getShortURL(shortURL string) (string, bool) {
	shortURL, ok := urlmap.shortToLong[shortURL]
	return shortURL, ok
}

func (urlmap *URLMap) putURL(longURL string) (string, bool) {
	if _, ok := urlmap.longToShort[longURL]; ok {
		return "", false
	}

	shortCand := generateRandomURL(urlmap.shortURLLength)
	for {
		_, ok := urlmap.shortToLong[shortCand]
		if !ok {
			break
		}
		shortCand = generateRandomURL(urlmap.shortURLLength)
	}

	urlmap.Lock()
	urlmap.shortToLong[shortCand] = longURL
	urlmap.longToShort[longURL] = shortCand
	urlmap.Unlock()

	return shortCand, true
}

func (urlmap *URLMap) delURL(shortURL string) bool {
	if longURL, ok := urlmap.shortToLong[shortURL]; ok {
		urlmap.Lock()
		delete(urlmap.shortToLong, shortURL)
		delete(urlmap.longToShort, longURL)
		urlmap.Unlock()
		return true
	}
	return false
}
