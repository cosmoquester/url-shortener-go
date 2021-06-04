package main

import (
	"math/rand"
	"sync"
)

type UrlMap struct {
	shortToLong, longToShort map[string]string
	sync.Mutex
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numLetters = len(letters)

func generateRandomUrl(n int) string {
	var result string

	for ; n > 0; n-- {
		result = result + string(letters[rand.Intn(numLetters)])
	}
	return result
}

func (urlmap *UrlMap) PutUrl(longUrl string) bool {
	if _, ok := urlmap.longToShort[longUrl]; ok {
		return false
	}

	shortCand := generateRandomUrl(7)
	for {
		_, ok := urlmap.shortToLong[shortCand]
		if !ok {
			break
		}
		shortCand = generateRandomUrl(7)
	}

	urlmap.Lock()
	urlmap.shortToLong[shortCand] = longUrl
	urlmap.longToShort[longUrl] = shortCand
	urlmap.Unlock()

	return true
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
