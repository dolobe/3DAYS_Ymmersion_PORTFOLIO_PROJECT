package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func HandleAddingPage(w http.ResponseWriter, r *http.Request) {
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

	data := PageData{
		IsAdmin:  isAdmin,
		Username: username,
	}

	if r.Method == http.MethodPost {
		if content := r.FormValue("content"); content != "" {
			if err := insertAbout(database, content); err != nil {
				log.Printf("Erreur lors de l'insertion du contenu : %s", err)
				http.Error(w, "Erreur lors de l'insertion du contenu", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
			return
		}

		if TitleCompetence := r.FormValue("TitleCompetence"); TitleCompetence != "" {
			ContentCompetence := r.FormValue("ContentCompetence")
			if err := insertCompetence(database, TitleCompetence, ContentCompetence); err != nil {
				log.Printf("Erreur lors de l'insertion de la compétence : %s", err)
				http.Error(w, "Erreur lors de l'insertion de la compétence", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
			return
		}

		if experienceTitle := r.FormValue("experienceTitle"); experienceTitle != "" {
			experienceDescription := r.FormValue("experienceDescription")
			if err := insertExperience(database, experienceTitle, experienceDescription); err != nil {
				log.Printf("Erreur lors de l'insertion de l'expérience : %s", err)
				http.Error(w, "Erreur lors de l'insertion de l'expérience", http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
			return
		}
	}

	tmpl, err := template.ParseFiles("templates/Adding.html")
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

func insertAbout(db *sql.DB, content string) error {
	insertAbout := `
		INSERT INTO about (content)
		VALUES (?);`

	_, err := db.Exec(insertAbout, content)
	if err != nil {
		log.Printf("Erreur lors de l'insertion du contenu : %v", err)
		return err
	}

	return nil
}

func insertCompetence(db *sql.DB, TitleCompetence, ContentCompetence string) error {
	insertCompetence := `
		INSERT INTO competence (TitleCompetence, ContentCompetence)
		VALUES (?, ?);`

	_, err := db.Exec(insertCompetence, TitleCompetence, ContentCompetence)
	if err != nil {
		log.Printf("Erreur lors de l'insertion de la compétence : %v", err)
		return err
	}

	return nil
}

func insertExperience(db *sql.DB, experienceTitle, experienceDescription string) error {
	insertExperience := `
		INSERT INTO experiences (experienceTitle, experienceDescription)
		VALUES (?, ?);`

	_, err := db.Exec(insertExperience, experienceTitle, experienceDescription)
	if err != nil {
		log.Printf("Erreur lors de l'insertion de l'expérience : %v", err)
		return err
	}

	return nil
}
