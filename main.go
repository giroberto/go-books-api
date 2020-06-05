package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
)

// Books field
type Books struct {
	Books []Book `json:"books"`
}

//Image field
type Image struct {
	Cover string `json:"cover"`
	Back  string `json:"back"`
}

// Stars field
type Stars struct {
	One   int8 `json:"one"`
	Two   int8 `json:"two"`
	Three int8 `json:"three"`
	Four  int8 `json:"four"`
	Fiver int8 `json:"five"`
}

// Book field
type Book struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Rate        float32 `json:"rate"`
	Stars       Stars   `json:"stars"`
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
		Handler: handlers.CORS()(r),
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Server running on port " + port)
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

	w.Header().Set("Content-Encoding", "br")
	w.Header().Set("Content-Type", "application/json")
	brw := brotli.NewWriterOptions(w, brotli.WriterOptions{Quality: 11})
	// gzw := gzip.NewWriter(w)
	// defer gzw.Close()
	defer brw.Close()
	json.NewEncoder(brw).Encode(books)
}
