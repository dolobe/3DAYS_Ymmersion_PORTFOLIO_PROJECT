package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

// Structure About pour contenir le contenu
type About struct {
	Content string
}

type Competence struct {
	TitleCompetence   string
	ContentCompetence string
}

type Experience struct {
	ExperienceTitle       string
	ExperienceDescription string
}

// Fonction pour gérer la page "À propos de moi"
func HandleAboutUsPage(w http.ResponseWriter, r *http.Request) {
	database, err := Path()
	if err != nil {
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer database.Close()

	isAdmin := isAdminAuthenticated(r)
	var username string

	if isAdmin {
		cookie, err := r.Cookie("admin_auth")
		if err == nil {
			err = database.QueryRow("SELECT username FROM admin WHERE username = ?", cookie.Value).Scan(&username)
			if err != nil {
				log.Printf("Erreur lors de la récupération du nom d'utilisateur : %v", err)
			}
		}
	}

	aboutUs, err := GetAboutUs(database)
	if err != nil {
		log.Printf("Erreur lors de la récupération du contenu : %s", err)
		http.Error(w, "Erreur lors de la récupération du contenu", http.StatusInternalServerError)
		return
	}

	competence, err := GetCompetence(database)
	if err != nil {
		log.Printf("Erreur lors de la récupération du contenu de compétence : %s", err)
		http.Error(w, "Erreur lors de la récupération du contenu de compétence", http.StatusInternalServerError)
		return
	}

	competences, err := GetAllCompetences(database)
	if err != nil {
		log.Printf("Erreur lors de la récupération des compétences : %s", err)
		http.Error(w, "Erreur lors de la récupération des compétences", http.StatusInternalServerError)
		return
	}

	// Dans votre fonction HandleAboutUsPage, récupérez les expériences :
	experiences, err := GetAllExperiences(database)
	if err != nil {
		log.Printf("Erreur lors de la récupération des expériences : %s", err)
		http.Error(w, "Erreur lors de la récupération des expériences", http.StatusInternalServerError)
		return
	}

	data := PageData{
		IsAdmin:           isAdmin,
		Username:          username,
		Content:           aboutUs.Content,
		CompetenceTitle:   competence.TitleCompetence,
		CompetenceContent: competence.ContentCompetence,
		Competences:       competences,
		Experiences:       experiences,
	}

	tmpl, err := template.ParseFiles("templates/AboutUs.html")
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}

func GetAboutUs(db *sql.DB) (About, error) {
	var aboutUs About

	err := db.QueryRow("SELECT content FROM about ORDER BY id DESC LIMIT 1").Scan(&aboutUs.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return About{}, nil
		}
		return About{}, err
	}

	return aboutUs, nil
}

func GetCompetence(db *sql.DB) (Competence, error) {
	var competence Competence

	err := db.QueryRow("SELECT TitleCompetence, ContentCompetence FROM competence ORDER BY id DESC LIMIT 1").Scan(&competence.TitleCompetence, &competence.ContentCompetence)
	if err != nil {
		if err == sql.ErrNoRows {
			return Competence{}, nil
		}
		return Competence{}, err
	}

	return competence, nil
}

func GetAllCompetences(db *sql.DB) ([]Competence, error) {
	var competences []Competence

	rows, err := db.Query("SELECT TitleCompetence, ContentCompetence FROM competence")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var competence Competence
		if err := rows.Scan(&competence.TitleCompetence, &competence.ContentCompetence); err != nil {
			return nil, err
		}
		competences = append(competences, competence)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return competences, nil
}

func GetAllExperiences(db *sql.DB) ([]Experience, error) {
	var experiences []Experience

	rows, err := db.Query("SELECT experienceTitle, experienceDescription FROM experiences") // Changez le nom de la table si nécessaire
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var experience Experience
		if err := rows.Scan(&experience.ExperienceTitle, &experience.ExperienceDescription); err != nil {
		}
		experiences = append(experiences, experience)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return experiences, nil
}
