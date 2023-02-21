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

type OrderDetails struct {
	Order_ID  int     `json:"order_id"`
	Book_Name string  `json:"book_name"`
	Address   string  `json:"address"`
	Amount    float64 `json:"amount"`
}

var db *sql.DB
var orders []OrderDetails

func main() {
	var err error
	db, err = sql.Open("mysql", "root:A@liya2020@tcp(localhost:3306)/bookstore")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/api/order/insert", insertOrder)
	http.HandleFunc("/api/order/retrieveAllOrders", getOrders)
	http.HandleFunc("/api/order/retrieveOrder/{id}", getOrderById)
	http.HandleFunc("/api/order/cancelOrder/{id}", cancelOrder)

	fmt.Println("Connected to 8001...!")
	http.ListenAndServe(":8001", nil)
}

// func handleOrderDetails(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		getOrders(w, r)
// 	case "POST":
// 		insertOrder(w, r)
// 	case "PUT":
// 		updateOrder(w, r)
// 	case "DELETE":
// 		cancelOrder(w, r)
// 	default:
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 	}
// }

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM order_details")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var order OrderDetails
		if err := rows.Scan(&order.Order_ID, &order.Book_Name, &order.Address, &order.Address); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orders = append(orders, order)
	}
	json.NewEncoder(w).Encode(orders)
}

func getOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var getOrder_Id OrderDetails
	rows, err := db.Query("SELECT * FROM order_details WHERE Order_ID=?", getOrder_Id.Order_ID)
	if err != nil {
		log.Fatalf("Record Not Found")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&getOrder_Id.Order_ID, &getOrder_Id.Book_Name, &getOrder_Id.Address, &getOrder_Id.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		orders = append(orders, getOrder_Id)
	}
	json.NewEncoder(w).Encode(orders)
}

func insertOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orderPost OrderDetails
	json.NewDecoder(r.Body).Decode(&orderPost)
	_, err := db.Exec("INSERT INTO order_details (Order_ID, Book_Name, Address, Amount) VALUES (?,?,?,?)", &orderPost.Order_ID, &orderPost.Book_Name, &orderPost.Address, &orderPost.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	orders = append(orders, orderPost)
	json.NewEncoder(w).Encode(orders) //optional
	w.Write([]byte("Data added successfully..."))
}
func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	order_Id, err := strconv.Atoi(r.URL.Path[len("/api/update/{id}"):])
	if err != nil {
		http.Error(w, "Invalid Request ID", http.StatusBadRequest)
		return
	}

	res, err := db.Exec("Update order_details set Order_ID=?, Book_Name=?, Address=?, Amount=? where Order_ID=?)", order_Id)
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
func cancelOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	order_Id, _ := strconv.Atoi(r.URL.Path[len("/api/delete/{Id}"):])
	cancelOrderId, err := db.Exec("DELETE FROM order_details WHERE Order_ID = ?", order_Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := cancelOrderId.RowsAffected()
	print(res)

	fmt.Println(" Deleted successfully!")
}
