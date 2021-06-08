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

// GetLongURL 은 입력된 shortURL에 대응되는 longURL을 반환합니다.
func (urlmap *URLMap) GetLongURL(shortURL string) (string, bool) {
	longURL, ok := urlmap.shortToLong[shortURL]
	return longURL, ok
}

// GetShortURL 은 입력된 longURL에 대응되는 shortURL을 반환합니다.
func (urlmap *URLMap) GetShortURL(shortURL string) (string, bool) {
	shortURL, ok := urlmap.shortToLong[shortURL]
	return shortURL, ok
}

// PutURL 은 longURL에 해당하는 shortURL을 생성해 mapping하고 shortURL을 반환합니다.
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

// DelURL 은 해당하는 shortURL 맵핑을 삭제합니다.
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
