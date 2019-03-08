package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

// Book model
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author model
type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// APIError model
type APIError struct {
	Message string `json:"error_message"`
}

var books []Book

func getBooks(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(books)
}

func getBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(res).Encode(book)
			return
		}
	}
	json.NewEncoder(res).Encode(&APIError{Message: "No matching data"})
}

func createBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var newBook Book
	_ = json.NewDecoder(req.Body).Decode(&newBook)

	if newBook.ID != "" {
		json.NewEncoder(res).Encode(
			&APIError{
				Message: fmt.Sprintf("id: %s was provided but not accepted on CREATE", newBook.ID),
			},
		)
		return
	}
	var newUUID, _ = uuid.NewUUID()
	newBook.ID = newUUID.String()

	books = append(books, newBook)

	json.NewEncoder(res).Encode(newBook)
}

func updateBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	var newBook Book
	_ = json.NewDecoder(req.Body).Decode(&newBook)

	for index, book := range books {
		if book.ID == params["id"] {
			newBook.ID = params["id"]
			copy(books[index+1:], books[index:])
			books[index] = newBook
			json.NewEncoder(res).Encode(books)
			return
		}
	}
}

func deleteBook(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	for index, book := range books {
		if book.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			json.NewEncoder(res).Encode(true)
			return
		}
	}
	json.NewEncoder(res).Encode(false)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
