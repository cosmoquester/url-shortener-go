package main

import (
	"bufio"
	"encoding/csv"
	"os"
	"sync"
)

// FileURLMap 은 LongURL와 ShortURL을 맵핑할 수 있는 map입니다.
// filePath를 가지고 있어 이 파일과 Sync합니다.
type FileURLMap struct {
	shortToLong, longToShort map[string]string
	shortURLLength           uint
	filePath                 string
	mapLock                  sync.Mutex
	fileLock                 sync.Mutex
}

// NewFileURLMap 은 지정한 파일경로에 url mapping을 기록하는 FileURLMap를 생성합니다.
// 이미 해당 경로에 파일이 있을 경우 그 Mapping정보를 불러와 사용합니다.
func NewFileURLMap(filePath string, urlLength uint) *FileURLMap {
	fileURLMap := FileURLMap{
		shortToLong:    make(map[string]string),
		longToShort:    make(map[string]string),
		shortURLLength: urlLength,
		filePath:       filePath,
	}

	var file *os.File
	var rows [][]string
	if _, err := os.Stat(filePath); err == nil {
		if file, err = os.Open(filePath); err != nil {
			panic("cannot read file")
		}

		defer file.Close()
		reader := csv.NewReader(bufio.NewReader(file))
		if rows, err = reader.ReadAll(); err != nil {
			panic("cannot read csv file")
		}

		for _, row := range rows {
			longURL, shortURL := row[0], row[1]
			fileURLMap.longToShort[longURL] = shortURL
			fileURLMap.shortToLong[shortURL] = longURL
		}
	}

	return &fileURLMap
}

func (urlmap *FileURLMap) getLongURL(shortURL string) (string, bool) {
	longURL, ok := urlmap.shortToLong[shortURL]
	return longURL, ok
}
func (urlmap *FileURLMap) getShortURL(longURL string) (string, bool) {
	shortURL, ok := urlmap.longToShort[longURL]
	return shortURL, ok
}

func (urlmap *FileURLMap) putURL(longURL string) (string, bool) {
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

	urlmap.mapLock.Lock()
	urlmap.shortToLong[shortCand] = longURL
	urlmap.longToShort[longURL] = shortCand
	urlmap.mapLock.Unlock()

	go urlmap.writeToFile()
	return shortCand, true
}

func (urlmap *FileURLMap) delURL(shortURL string) bool {
	if longURL, ok := urlmap.shortToLong[shortURL]; ok {
		urlmap.mapLock.Lock()
		delete(urlmap.shortToLong, shortURL)
		delete(urlmap.longToShort, longURL)
		urlmap.mapLock.Unlock()

		go urlmap.writeToFile()
		return true
	}
	return false
}

func (urlmap *FileURLMap) writeToFile() {
	urlmap.fileLock.Lock()
	defer urlmap.fileLock.Unlock()

	file, err := os.Create(urlmap.filePath)
	if err != nil {
		panic("cannot write to file")
	}
	defer file.Close()
	writer := csv.NewWriter(bufio.NewWriter(file))

	urlmap.mapLock.Lock()
	defer urlmap.mapLock.Unlock()
	for longURL, shortURL := range urlmap.longToShort {
		writer.Write([]string{longURL, shortURL})
	}
	writer.Flush()
}
