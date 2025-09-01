package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

type Weight struct {
	ID       int       `db:"id"`
	Date     time.Time `db:"date"`
	WeightKg float64   `db:"weight_kg"`
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
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/add", addHandler).Methods("GET", "POST")
	r.HandleFunc("/weights", getWeights).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		date := r.FormValue("date")
		weightStr := r.FormValue("weight")

		// Convert string to float64
		weight, err := strconv.ParseFloat(weightStr, 64)
		if err != nil {
			http.Error(w, "Invalid weight value", http.StatusBadRequest)
			return
		}

		// Insert into DB
		_, err = db.Exec("INSERT INTO weights (date, weight_kg) VALUES ($1, $2)", date, weight)
		if err != nil {
			http.Error(w, "Error saving weight: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Show success page
		data := map[string]interface{}{
			"Date":   date,
			"Weight": weight,
		}
		renderTemplate(w, "add_success.html", data)
		return
	}
	renderTemplate(w, "add.html", nil)
}

func getWeights(w http.ResponseWriter, r *http.Request) {
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

	// If request is JSON (API style)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(weights)
		return
	}

	// Otherwise render HTML
	renderTemplate(w, "weights.html", weights)
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
