package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/cosmoquester/url-shortener-go/urlconverter"
	"github.com/cosmoquester/url-shortener-go/utils"
	"github.com/gorilla/mux"
)

var converter urlconverter.URLConverter

func createShortURL(w http.ResponseWriter, req *http.Request) {
	var data []byte
	var err error
	body := make(map[string]string)

	if data, err = ioutil.ReadAll(req.Body); err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Cannot read request!", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(data, &body); err != nil {
		log.Println("Error occurred:", err)
		http.Error(w, "Cannot parse as json!", http.StatusBadRequest)
		return
	}

	longURL := body["long_url"]
	if !utils.ValidateURL(longURL) {
		log.Println("Error occurred: Invalid url form in creating")
		http.Error(w, "Invalid url! the url must start with http or https", http.StatusBadRequest)
		return
	}

	if shortURL, ok := converter.PutURL(longURL); ok {
		log.Println("resource created ", longURL, "to", shortURL)
		w.Write([]byte(fmt.Sprintf("{\"result\":true, \"short_url\":%s}", shortURL)))
	} else {
		log.Println("Error occurred in putting", err)
		http.Error(w, fmt.Sprintf("long url: %s is already exists!", longURL), http.StatusConflict)
	}
}

func forwardURL(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if shortURL, ok := vars["short_url"]; !ok {
		log.Println("forward fail invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longURL, ok := converter.GetLongURL(shortURL); !ok {
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
		log.Println("deletion failed invalid form")
		http.Error(w, "Invalid form!", http.StatusBadRequest)
		return
	} else if longURL, ok := converter.GetLongURL(shortURL); !ok {
		log.Println("deletion failed non-existing shortURL")
		http.Error(w, "Invalid short_url!", http.StatusNotFound)
		return
	} else if ok := converter.DelURL(shortURL); !ok {
		log.Println("deletion failed")
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
	port := flag.String("port", "5000", "port number")
	flag.Parse()

	if *converterType == "urlmap" {
		converter = urlconverter.NewURLMap(*urlLength)
		log.Println("use urlmap, data will be deleted with end of process")
	} else if *converterType == "file-urlmap" {
		fileURLMap, err := urlconverter.NewFileURLMap(*urlMapFilePath, *urlLength)
		if err != nil {
			log.Fatal("error occurred with file path: ", *urlMapFilePath, err)
		}
		converter = fileURLMap
		log.Println("use file-urlmap mapping file path: ", *urlMapFilePath)
	}
	log.Println("short url length: ", *urlLength)
	log.Println("port:", *port)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", createShortURL).Methods("POST")
	router.HandleFunc("/{short_url}", forwardURL).Methods("GET")
	router.HandleFunc("/{short_url}", deleteURL).Methods("DELETE")

	rand.Seed(time.Now().UnixNano())
	http.ListenAndServe(":"+*port, router)
}
