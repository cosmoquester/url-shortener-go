package urlconverter

import (
	"sync"

	"github.com/cosmoquester/url-shortener-go/utils"
)

// URLMap 은 LongURL와 ShortURL을 맵핑할 수 있는 map입니다.
type URLMap struct {
	shortToLong, longToShort map[string]string
	shortURLLength           uint
	sync.Mutex
}

// NewURLMap short url과 long url를 서로 맵핑해주는 Map입니다.
func NewURLMap(urlLength uint) *URLMap {
	return &URLMap{
		shortToLong:    make(map[string]string),
		longToShort:    make(map[string]string),
		shortURLLength: urlLength,
	}
}

func (urlmap *URLMap) GetLongURL(shortURL string) (string, bool) {
	longURL, ok := urlmap.shortToLong[shortURL]
	return longURL, ok
}
func (urlmap *URLMap) GetShortURL(shortURL string) (string, bool) {
	shortURL, ok := urlmap.shortToLong[shortURL]
	return shortURL, ok
}

func (urlmap *URLMap) PutURL(longURL string) (string, bool) {
	if _, ok := urlmap.longToShort[longURL]; ok {
		return "", false
	}

	shortCand := utils.GenerateRandomURL(urlmap.shortURLLength)
	for {
		_, ok := urlmap.shortToLong[shortCand]
		if !ok {
			break
		}
		shortCand = utils.GenerateRandomURL(urlmap.shortURLLength)
	}

	urlmap.Lock()
	urlmap.shortToLong[shortCand] = longURL
	urlmap.longToShort[longURL] = shortCand
	urlmap.Unlock()

	return shortCand, true
}

func (urlmap *URLMap) DelURL(shortURL string) bool {
	longURL, ok := urlmap.shortToLong[shortURL]
	if !ok {
		return false
	}
	urlmap.Lock()
	delete(urlmap.shortToLong, shortURL)
	delete(urlmap.longToShort, longURL)
	urlmap.Unlock()
	return true
}
