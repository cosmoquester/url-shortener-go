package main

import (
	"encoding/json"
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
		log.Println("Error occured:", err)
		return
	}

	if err := json.Unmarshal(data, &body); err != nil {
		log.Println("Error occured:", err)
		return
	}

	longURL := body["long_url"]
	if !validateURL(longURL) {
		log.Println("Error occured: Invalid url form in creating")
		http.Error(w, "Invalid url! the url must start with http or https", http.StatusBadRequest)
		return
	}

	if shortURL, ok := converter.putURL(longURL); ok {
		log.Println("resource created ", longURL, "to", shortURL)
		w.Write([]byte("{\"result\":true}"))
	} else {
		log.Println("Error occured in putting", err)
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
	converter = &URLMap{
		shortToLong:    make(map[string]string),
		longToShort:    make(map[string]string),
		shortURLLength: 7,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", createShortURL).Methods("POST")
	router.HandleFunc("/{short_url}", forwardURL).Methods("GET")
	router.HandleFunc("/{short_url}", deleteURL).Methods("DELETE")

	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(":5000", router)
}
