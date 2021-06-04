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

type UrlConverter interface {
	getShortUrl(string) (string, bool)
	getLongUrl(string) (string, bool)
	PutUrl(string) (string, bool)
	DelUrl(string) bool
}

var converter UrlConverter

func CreateShortUrl(w http.ResponseWriter, req *http.Request) {
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

	longUrl := body["long_url"]
	if !validateUrl(longUrl) {
		log.Println("Error occured: Invalid url form in creating")
		http.Error(w, "Invalid url! the url must start with http or https", http.StatusBadRequest)
		return
	}

	if shortUrl, ok := converter.PutUrl(longUrl); ok {
		log.Println("resource created ", longUrl, "to", shortUrl)
		w.Write([]byte("{\"result\":true}"))
	} else {
		log.Println("Error occured in putting", err)
	}
}

func ForwardUrl(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if shortUrl, ok := vars["short_url"]; !ok {
		log.Println("forward fail invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longUrl, ok := converter.getLongUrl(shortUrl); !ok {
		log.Println("forward failed non-existing shorturl", shortUrl)
		http.Error(w, "Invalid short_url!", http.StatusNotFound)
		return
	} else {
		log.Println("forward from", shortUrl, "to", longUrl)
		http.Redirect(w, req, longUrl, http.StatusFound)
	}
}

func DeleteUrl(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if shortUrl, ok := vars["short_url"]; !ok {
		log.Println("forward failed invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longUrl, ok := converter.getLongUrl(shortUrl); !ok {
		log.Println("forward failed non-existing shortUrl")
		http.Error(w, "Invalid short_url!", http.StatusNotFound)
		return
	} else if ok := converter.DelUrl(shortUrl); !ok {
		http.Error(w, "Internal Server error!", http.StatusInternalServerError)
		return
	} else {
		log.Println("resource deleted", shortUrl, "to", longUrl)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte("{\"result\":true}"))
	}
}

func main() {
	converter = &UrlMap{
		shortToLong:    make(map[string]string),
		longToShort:    make(map[string]string),
		shortUrlLength: 7,
	}
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", CreateShortUrl).Methods("POST")
	router.HandleFunc("/{short_url}", ForwardUrl).Methods("GET")
	router.HandleFunc("/{short_url}", ForwardUrl).Methods("DELETE")

	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(":5000", router)
}
