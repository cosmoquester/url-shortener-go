package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// URLConverter 는 LongURL과 ShortURL을 기록하고 변환하고 삭제할 수 있는 interface입니다.
type URLConverter interface {
	getShortURL(string) (string, bool)
	getLongURL(string) (string, bool)
	putURL(string) (string, bool)
	delURL(string) bool
}

var converter URLConverter

func createShortURL(w http.ResponseWriter, req *http.Request) {
	var data []byte
	var err error
	body := make(map[string]string)

	if data, err = ioutil.ReadAll(req.Body); err != nil {
		log.Println("Error occurred:", err)
		return
	}

	if err := json.Unmarshal(data, &body); err != nil {
		log.Println("Error occurred:", err)
		return
	}

	longURL := body["long_url"]
	if !validateURL(longURL) {
		log.Println("Error occurred: Invalid url form in creating")
		http.Error(w, "Invalid url! the url must start with http or https", http.StatusBadRequest)
		return
	}

	if shortURL, ok := converter.putURL(longURL); ok {
		log.Println("resource created ", longURL, "to", shortURL)
		w.Write([]byte("{\"result\":true}"))
	} else {
		log.Println("Error occurred in putting", err)
	}
}

func forwardURL(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if shortURL, ok := vars["short_url"]; !ok {
		log.Println("forward fail invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longURL, ok := converter.getLongURL(shortURL); !ok {
		log.Println("forward failed non-existing shorturl", shortURL)
		http.Error(w, "Invalid short_url!", http.StatusNotFound)
		return
	} else {
		log.Println("forward from", shortURL, "to", longURL)
		http.Redirect(w, req, longURL, http.StatusFound)
	}
}

func deleteURL(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if shortURL, ok := vars["short_url"]; !ok {
		log.Println("forward failed invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longURL, ok := converter.getLongURL(shortURL); !ok {
		log.Println("forward failed non-existing shortURL")
		http.Error(w, "Invalid short_url!", http.StatusNotFound)
		return
	} else if ok := converter.delURL(shortURL); !ok {
		http.Error(w, "Internal Server error!", http.StatusInternalServerError)
		return
	} else {
		log.Println("resource deleted", shortURL, "to", longURL)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("{\"result\":true}"))
	}
}

func main() {
	converterType := flag.String("converter", "file-urlmap", "type of URL converter \"urlmap\" or \"file-urlmap\"")
	urlMapFilePath := flag.String("url-map-file", "url-mapping.csv", "file path mapping longurl to shorturl with \"file-urlmap\"")
	urlLength := flag.Uint("url-length", 7, "the length of short url")
	flag.Parse()

	if *converterType == "urlmap" {
		converter = &URLMap{
			shortToLong:    make(map[string]string),
			longToShort:    make(map[string]string),
			shortURLLength: *urlLength,
		}
		log.Println("use urlmap, data will be deleted with end of process")
	} else if *converterType == "file-urlmap" {
		converter = NewFileURLMap(*urlMapFilePath, *urlLength)
		log.Println("use file-urlmap mapping file path:", *urlMapFilePath)
	}
	log.Println("short url length: ", *urlLength)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", createShortURL).Methods("POST")
	router.HandleFunc("/{short_url}", forwardURL).Methods("GET")
	router.HandleFunc("/{short_url}", deleteURL).Methods("DELETE")

	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(":5000", router)
}
