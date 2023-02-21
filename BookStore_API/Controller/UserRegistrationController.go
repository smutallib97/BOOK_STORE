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

type UserRegistration struct {
	User_ID   int    `json:"user_id"`
	User_Name string `json:"user_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var db *sql.DB
var users []UserRegistration

func main() {
	var err error
	db, err = sql.Open("mysql", "root:A@liya2020@tcp(localhost:3306)/bookstore")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/userservice", getUsers)
	http.HandleFunc("/api/userservice/{userId}", getUser)
	http.HandleFunc("/api/userservice/register", postUser)
	http.HandleFunc("/api/userservice/login", userLogin)
	http.HandleFunc("/api/userservice/delete/{userId}", deleteUser)

	fmt.Println("Connected to 8001...!")
	http.ListenAndServe(":8001", nil)
}

// func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		getUsers(w, r)
// 	case "POST":
// 		postUser(w, r)
// 	case "PUT":
// 		putUser(w, r)
// 	case "DELETE":
// 		deleteUser(w, r)
// 	default:
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user UserRegistration
		if err := rows.Scan(&user.User_ID, &user.User_Name, &user.Email, &user.Password); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Path[len("/api/userservice/{userId}"):])

	var getUser UserRegistration
	rows, err := db.Query("SELECT * FROM users WHERE user_ID=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&getUser.User_ID, &getUser.User_Name, &getUser.Email, &getUser.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, getUser)
	}
	json.NewEncoder(w).Encode(users)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userPost UserRegistration
	json.NewDecoder(r.Body).Decode(&userPost)
	_, err := db.Exec("INSERT INTO users (user_ID, userName, email, password) VALUES (?,?,?,?)", &userPost.User_ID, &userPost.User_Name, &userPost.Email, &userPost.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	users = append(users, userPost)
	json.NewEncoder(w).Encode(users) //optional
	w.Write([]byte("Data added successfully..."))
}
func putUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(r.URL.Path[len("/api/userservice/{userId}"):])
	if err != nil {
		http.Error(w, "Invalid Request ID", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("Update users set user_ID=?, userName=?, email=?, password=? where user_ID=?)", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updated, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
	print(updated)
	w.Write([]byte(" Data Updated successfully..."))
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, _ := strconv.Atoi(r.URL.Path[len("/userservice/{userId}"):])
	deleteId, err := db.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := deleteId.RowsAffected()
	print(res)

	fmt.Println(" Deleted successfully!")
}
func userLogin(w http.ResponseWriter, r *http.Request) {
	var request UserRegistration
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND password = ?", request.User_Name, request.Password).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Incorrect username or password"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
