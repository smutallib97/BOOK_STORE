package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type BookDetails struct {
	Book_ID     int    `json:"book_id"`
	Book_Title  string `json:"book_title"`
	Book_Author string `json:"book_author"`
	Book_Price  int    `json:"book_price"`
}

var db *sql.DB
var books []BookDetails

func main() {
	var err error
	db, err = sql.Open("mysql", "root:A@liya2020@tcp(localhost:3306)/bookstore")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/getBooks", getBooks)
	http.HandleFunc("/api//getBookByName/{bookName}", getBookByName)
	http.HandleFunc("/api/getBook/{bookId}", getBook)
	http.HandleFunc("/api/addBook", addBook)
	http.HandleFunc("/api/update/{bookId}", updateBook)
	http.HandleFunc("/api/delete/{bookId}", deleteBook)

	fmt.Println("Connected to 8001...!")
	http.ListenAndServe(":8001", nil)
}

// func handleBookDetails(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		getBooks(w, r)
// 	case "POST":
// 		addBook(w, r)
// 	case "PUT":
// 		updateBook(w, r)
// 	case "DELETE":
// 		deleteBook(w, r)
// 	default:
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var book BookDetails
		if err := rows.Scan(&book.Book_ID, &book.Book_Title, &book.Book_Author, &book.Book_Price); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}
	json.NewEncoder(w).Encode(books)
}
func getBookByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var getBookByName BookDetails
	rows, err := db.Query("SELECT * FROM books WHERE Book_Title=?", getBookByName.Book_Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(getBookByName.Book_ID, getBookByName.Book_Title, getBookByName.Book_Author, getBookByName.Book_Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, getBookByName)
	}
	json.NewEncoder(w).Encode(books)
}
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var getBook_Id BookDetails
	rows, err := db.Query("SELECT * FROM books WHERE Book_ID=?", getBook_Id.Book_ID)
	if err != nil {
		log.Fatalf("Record Not Found")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&getBook_Id.Book_ID, &getBook_Id.Book_Title, &getBook_Id.Book_Author, &getBook_Id.Book_Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		books = append(books, getBook_Id)
	}
	json.NewEncoder(w).Encode(books)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bookPost BookDetails
	json.NewDecoder(r.Body).Decode(&bookPost)
	_, err := db.Exec("INSERT INTO books (Book_ID, Book_Title, Book_Author, Book_Price) VALUES (?,?,?,?)", &bookPost.Book_ID, &bookPost.Book_Title, &bookPost.Book_Author, &bookPost.Book_Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	books = append(books, bookPost)
	json.NewEncoder(w).Encode(books) //optional
	w.Write([]byte("Data added successfully..."))
}
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, err := strconv.Atoi(r.URL.Path[len("/api/update/{bookId}"):])
	if err != nil {
		http.Error(w, "Invalid Request ID", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("Update books set Book_ID=?, Book_Title=?, Book_Author=?, Book_Price=? where Book_ID=?)", bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updated, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	print(updated)
	w.Write([]byte("Record updated successfully..."))
}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bookId, _ := strconv.Atoi(r.URL.Path[len("/api/delete/{bookId}"):])
	deleteId, err := db.Exec("DELETE FROM books WHERE Book_ID = ?", bookId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := deleteId.RowsAffected()
	print(res)

	fmt.Println(" Deleted successfully!")
}
