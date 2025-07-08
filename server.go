package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
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
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "https://aaronw.link")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Forwarded-For, User-Agent")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Try to get real IP from X-Forwarded-For header
	ip := getClientIP(r)

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

func getClientIP(r *http.Request) string {
	// Try X-Forwarded-For first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Multiple IPs are comma-separated, take the first
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to RemoteAddr (format: IP:port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback if something goes wrong
	}
	return ip
}
