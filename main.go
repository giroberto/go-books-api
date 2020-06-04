package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
)

// Books aqui
type Books struct {
	Books []Book `json:"books"`
}

//Image aqui
type Image struct {
	Cover string `json:"cover"`
	Back  string `json:"back"`
}

// Book aqui
type Book struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Rate        float32 `json:"rate"`
	Ratings     int8    `json:"ratings"`
	Description string  `json:"description"`
	Image       Image   `json:"images"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", apiInfo)
	r.HandleFunc("/books", listBooks)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server running on port 3000")
	log.Fatal(srv.ListenAndServe())
}

func apiInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Running...")
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	jsonFile, err := os.Open("bookdata.json")
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var books Books

	json.Unmarshal(byteValue, &books)

	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "application/json")
	gzw := gzip.NewWriter(w)
	defer gzw.Close()
	json.NewEncoder(gzw).Encode(books)
}
