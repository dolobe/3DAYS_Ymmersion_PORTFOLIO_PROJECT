package handlers

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

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

func createTables(db *sql.DB) error {
	createExperienceTable := `
	CREATE TABLE IF NOT EXISTS experiences (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		experienceTitle TEXT NOT NULL,
		experienceDescription TEXT NOT NULL
	);`

	if _, err := db.Exec(createExperienceTable); err != nil {
		log.Println("Erreur lors de la création de la table experience:", err)
		return err
	}

	createCompetenceTable := `
	CREATE TABLE IF NOT EXISTS competence (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		TitleCompetence TEXT NOT NULL,
		ContentCompetence TEXT NOT NULL
	);`

	if _, err := db.Exec(createCompetenceTable); err != nil {
		log.Println("Erreur lors de la création de la table competence:", err)
		return err
	}

	createAdminTable := `
    CREATE TABLE IF NOT EXISTS admin (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	_, err := db.Exec(createAdminTable)
	if err != nil {
		log.Fatalf("Erreur lors de la création de la table admin: %v", err)
		return err
	}

	createAboutTable := `
	CREATE TABLE IF NOT EXISTS about (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL
	);`

	if _, err := db.Exec(createAboutTable); err != nil {
		log.Println("Erreur lors de la création de la table about:", err)
		return err
	}

	// Création de la table projects
	queryProjects := `
    CREATE TABLE IF NOT EXISTS projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
        image TEXT,
        link TEXT
    );`

	if _, err := db.Exec(queryProjects); err != nil {
		log.Println("Erreur lors de la création de la table projects:", err)
		return err
	}

	// Création de la table messages
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		subject TEXT NOT NULL,
		message TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := db.Exec(createTableSQL); err != nil {
		log.Println("Erreur lors de la création de la table messages:", err)
		return err
	}

	return nil
}
