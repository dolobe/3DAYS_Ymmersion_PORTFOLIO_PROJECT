package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}

var db *sql.DB

func Path() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./src/DataBase/database.db")
	if err != nil {
		return nil, err
	}

	if err := createTables(database); err != nil {
		return nil, err
	}

	db = database
	return db, nil
}

func HandleAddProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var project Project
		err := json.NewDecoder(r.Body).Decode(&project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO projects (title, description, image, link) VALUES (?, ?, ?, ?)",
			project.Title, project.Description, project.Image, project.Link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(project)
	} else {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func createTables(db *sql.DB) error {
	queryUsers := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        mail TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	if _, err := db.Exec(queryUsers); err != nil {
		log.Println("Erreur lors de la création de la table utilisateurs:", err)
		return err
	}

	queryProjects := `
    CREATE TABLE IF NOT EXISTS projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT NOT NULL,
        image TEXT NOT NULL,
        link TEXT NOT NULL
    );`

	if _, err := db.Exec(queryProjects); err != nil {
		log.Println("Erreur lors de la création de la table projets:", err)
		return err
	}

	return nil
}
