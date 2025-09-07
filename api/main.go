package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

type Weight struct {
	ID       int       `json:"id"`
	Date     time.Time `json:"date"`
	WeightKg float64   `json:"weight_kg"`
}

func main() {
	var err error

	// THIS IS FOR LOGGING INTO THE DB
	// postgres://<USER>:<PASSWORD>@localhost:5432/<DBNAME>?sslmode=disable
	connStr := "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Verify DB connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot reach database:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/add", addHandler).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/weights", weightsHandler).Methods("GET", "OPTIONS")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}
	renderTemplate(w, "index.html", nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Date   string  `json:"date"`
			Weight float64 `json:"weight"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode JSON: %v", err)
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Received weight entry: date=%s weight=%.2f", req.Date, req.Weight)

		_, err := db.Exec("INSERT INTO weights (date, weight_kg) VALUES ($1, $2)", req.Date, req.Weight)
		if err != nil {
			log.Printf("DB insert error: %v", err)
			http.Error(w, "Error saving weight: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	renderTemplate(w, "add.html", nil)
}

func weightsHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)
	if r.Method == "OPTIONS" {
		return
	}
	rows, err := db.Query("SELECT id, date, weight_kg FROM weights ORDER BY date DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var weights []Weight
	for rows.Next() {
		var wt Weight
		err := rows.Scan(&wt.ID, &wt.Date, &wt.WeightKg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		weights = append(weights, wt)
	}

	// Always return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weights)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseGlob("templates/*.html") // parse every request
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
