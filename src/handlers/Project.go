package handlers

import (
	"database/sql"
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

// HandleProjectPage gère la page de projets
func HandleProjectPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Démarrage de HandleProjectPage")

	database, err := Path()
	if err != nil {
		log.Printf("Erreur de connexion à la base de données : %v", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Erreur lors de la fermeture de la base de données : %v", err)
		}
	}()

	isAdmin := isAdminAuthenticated(r)
	var username string

	if isAdmin {
		log.Println("Utilisateur administrateur détecté")
		cookie, err := r.Cookie("admin_auth")
		if err == nil {
			err = database.QueryRow("SELECT username FROM admin WHERE username = ?", cookie.Value).Scan(&username)
			if err != nil {
				log.Printf("Erreur lors de la récupération du nom d'utilisateur : %v", err)
			} else {
				log.Printf("Nom d'utilisateur récupéré : %s", username)
			}
		}
	}

	if r.Method == http.MethodPost {
		log.Println("Méthode POST détectée pour l'ajout d'un projet")

		title := r.FormValue("project-title")
		description := r.FormValue("project-description")
		image := r.FormValue("project-image")
		link := r.FormValue("project-link")

		log.Printf("Projet à insérer : Titre=%s, Description=%s, Image=%s, Lien=%s", title, description, image, link)

		if err := insertProject(database, title, description, image, link); err != nil {
			log.Printf("Erreur lors de l'insertion du projet : %v", err)
			http.Error(w, "Erreur lors de l'insertion du projet", http.StatusInternalServerError)
			return
		}

		log.Println("Projet inséré avec succès")
		http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
		return
	}

	projects, err := GetProjects(database)
	if err != nil {
		log.Printf("Erreur lors de la récupération des projets : %v", err)
		http.Error(w, "Erreur lors de la récupération des projets", http.StatusInternalServerError)
		return
	}

	data := PageData{
		IsAdmin:  isAdmin,
		Username: username,
		Projects: projects,
	}

	tmpl, err := template.ParseFiles("templates/Project.html")
	if err != nil {
		log.Printf("Erreur lors du rendu du template : %v", err)
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Erreur lors de l'exécution du template : %v", err)
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}

	log.Println("Page des projets chargée avec succès.")
}

// GetProjects récupère tous les projets de la base de données
func GetProjects(database *sql.DB) ([]Project, error) {
	log.Println("Récupération des projets")
	rows, err := database.Query("SELECT id, title, description, image, link FROM projects")
	if err != nil {
		log.Printf("Erreur lors de la récupération des projets : %v", err)
		return nil, err
	}
	defer rows.Close()

	projects := []Project{}
	for rows.Next() {
		project := Project{}
		err := rows.Scan(&project.ID, &project.Title, &project.Description, &project.Image, &project.Link)
		if err != nil {
			log.Printf("Erreur lors de la lecture d'un projet : %v", err)
			return nil, err
		}
		log.Printf("Projet récupéré : %+v", project)
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Erreur après l'itération des lignes : %v", err)
		return nil, err
	}

	log.Printf("Total des projets récupérés : %d", len(projects))
	return projects, nil
}
