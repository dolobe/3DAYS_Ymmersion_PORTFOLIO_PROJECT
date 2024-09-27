package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Link        string `json:"link"`
}

func HandleProjectPage(w http.ResponseWriter, r *http.Request) {
	database, err := Path()
	if err != nil {
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	isAdmin := isAdminAuthenticated(r)
	var username string

	if isAdmin {
		// Récupérer le nom d'utilisateur de la base de données
		cookie, err := r.Cookie("admin_auth")
		if err == nil {
			err = database.QueryRow("SELECT username FROM admin WHERE username = ?", cookie.Value).Scan(&username)
			if err != nil {
				log.Printf("Erreur lors de la récupération du nom d'utilisateur : %v", err)
			}
		}
	}

	// Récupérer les projets pour les passer au template
	projects, err := GetAllProjects(database)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des projets", http.StatusInternalServerError)
		return
	}

	data := PageData{
		IsAdmin:  isAdmin,
		Username: username,
		Projects: projects, // Passer les projets au template
	}

	tmpl, err := template.ParseFiles("templates/Project.html")
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}

	log.Println("Page des projets chargée avec succès.")
}

func GetAllProjects(db *sql.DB) ([]Project, error) {
	log.Println("Tentative de récupération de tous les projets...")
	rows, err := db.Query("SELECT id, title, description, image, link FROM projects")
	if err != nil {
		log.Printf("Erreur lors de la récupération des projets : %s", err)
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.Image, &project.Link); err != nil {
			log.Printf("Erreur lors du scan des données du projet : %s", err)
			return nil, err
		}
		projects = append(projects, project)
		log.Printf("Projet récupéré : %+v", project)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Erreur lors de l'itération sur les lignes : %s", err)
		return nil, err
	}

	return projects, nil
}

func AddProjectHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Tentative d'ajout d'un nouveau projet...")
	var project Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		log.Printf("Erreur lors du décodage du projet : %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Données du projet reçues : %+v", project)

	db, err := Path()
	if err != nil {
		log.Printf("Erreur de connexion à la base de données : %s", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO projects (title, description, image, link) VALUES (?, ?, ?, ?)",
		project.Title, project.Description, project.Image, project.Link)

	if err != nil {
		log.Printf("Erreur lors de l'insertion du projet : %s", err)
		http.Error(w, "Erreur lors de l'insertion du projet", http.StatusInternalServerError)
		return
	}

	log.Printf("Projet ajouté avec succès : %s", project.Title)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(project); err != nil {
		log.Printf("Erreur lors de l'encodage du projet en JSON : %s", err)
		http.Error(w, "Erreur lors de l'encodage du projet en JSON", http.StatusInternalServerError)
		return
	}
	log.Println("Réponse envoyée avec succès au client.")
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Tentative de récupération de tous les projets...")
	w.Header().Set("Content-Type", "application/json")

	db, err := Path()
	if err != nil {
		log.Printf("Erreur lors de la connexion à la base de données : %s", err)
		http.Error(w, "Erreur lors de la connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	projects, err := GetAllProjects(db)
	if err != nil {
		log.Printf("Erreur lors de la récupération des projets : %s", err)
		http.Error(w, "Erreur lors de la récupération des projets", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(projects); err != nil {
		log.Printf("Erreur lors de l'encodage des projets en JSON : %s", err)
		http.Error(w, "Erreur lors de l'encodage des projets en JSON", http.StatusInternalServerError)
		return
	}
	log.Println("Tous les projets envoyés avec succès au client.")
}
