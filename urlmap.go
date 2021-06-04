package main

import (
	"sync"
)

type UrlMap struct {
	shortToLong, longToShort map[string]string
	shortUrlLength           uint
	sync.Mutex
}

func (urlmap *UrlMap) getLongUrl(shortUrl string) (string, bool) {
	longUrl, ok := urlmap.shortToLong[shortUrl]
	return longUrl, ok
}
func (urlmap *UrlMap) getShortUrl(shortUrl string) (string, bool) {
	shortUrl, ok := urlmap.shortToLong[shortUrl]
	return shortUrl, ok
}

func (urlmap *UrlMap) PutUrl(longUrl string) (string, bool) {
	if _, ok := urlmap.longToShort[longUrl]; ok {
		return "", false
	}

	shortCand := generateRandomUrl(urlmap.shortUrlLength)
	for {
		_, ok := urlmap.shortToLong[shortCand]
		if !ok {
			break
		}
		shortCand = generateRandomUrl(urlmap.shortUrlLength)
	}

	urlmap.Lock()
	urlmap.shortToLong[shortCand] = longUrl
	urlmap.longToShort[longUrl] = shortCand
	urlmap.Unlock()

	return shortCand, true
}

func (urlmap *UrlMap) DelUrl(shortUrl string) bool {
	if longUrl, ok := urlmap.shortToLong[shortUrl]; !ok {
		return false
	} else {
		urlmap.Lock()
		delete(urlmap.shortToLong, shortUrl)
		delete(urlmap.longToShort, longUrl)
		urlmap.Unlock()

		return true
	}
}
