package handlers

import (
	"encoding/json"
	"io/ioutil"
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
	log.Println("Tentative de chargement de la page des projets...")
	htmlFile, err := ioutil.ReadFile("templates/Project.html")
	if err != nil {
		log.Printf("Erreur lors de la lecture du fichier HTML : %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	_, err = w.Write(htmlFile)
	if err != nil {
		log.Printf("Erreur lors de l'écriture de la réponse : %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Page des projets chargée avec succès.")
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

	rows, err := db.Query("SELECT id, title, description, image, link FROM projects")
	if err != nil {
		log.Printf("Erreur lors de la récupération des projets : %s", err)
		http.Error(w, "Erreur lors de la récupération des projets", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projets []Project
	for rows.Next() {
		var projet Project
		if err := rows.Scan(&projet.ID, &projet.Title, &projet.Description, &projet.Image, &projet.Link); err != nil {
			log.Printf("Erreur lors du scan des données du projet : %s", err)
			http.Error(w, "Erreur lors de la récupération des projets", http.StatusInternalServerError)
			return
		}
		projets = append(projets, projet)
		log.Printf("Projet récupéré : %+v", projet)
	}

	if err := json.NewEncoder(w).Encode(projets); err != nil {
		log.Printf("Erreur lors de l'encodage des projets en JSON : %s", err)
		http.Error(w, "Erreur lors de l'encodage des projets en JSON", http.StatusInternalServerError)
		return
	}
	log.Println("Tous les projets envoyés avec succès au client.")
}
