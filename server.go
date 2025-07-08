package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./data/visits.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("/log", logHandler)
	fmt.Println("Server started on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS visits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT,
		timestamp DATETIME,
		user_agent TEXT
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := r.RemoteAddr
	userAgent := r.UserAgent()
	timestamp := time.Now()

	_, err := db.Exec("INSERT INTO visits (ip, timestamp, user_agent) VALUES (?, ?, ?)", ip, timestamp, userAgent)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged"))
}
