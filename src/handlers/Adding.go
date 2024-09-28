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
		// Vérification de l'insertion de contenu 'À propos de moi'
		if r.URL.Path == "/Adding" {
			if content := r.FormValue("content"); content != "" {
				if err := insertAbout(database, content); err != nil {
					log.Printf("Erreur lors de l'insertion du contenu : %s", err)
					http.Error(w, "Erreur lors de l'insertion du contenu", http.StatusInternalServerError)
					return
				}
				http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
				return
			}
		}

		// Vérification de l'insertion de la compétence
		if r.URL.Path == "/Adding" {
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
		}

		// Vérification de l'insertion de l'expérience
		if r.URL.Path == "/Adding" {
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

		// Vérification de l'insertion du projet
		if r.URL.Path == "/Project" {
			if title := r.FormValue("project-title"); title != "" {
				description := r.FormValue("project-description")
				image := r.FormValue("project-image")
				link := r.FormValue("project-link")

				if err := insertProject(database, title, description, image, link); err != nil {
					log.Printf("Erreur lors de l'insertion du projet : %s", err)
					http.Error(w, "Erreur lors de l'insertion du projet", http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
				return
			}
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

func insertProject(database *sql.DB, title, description, image, link string) error {
	log.Println("Début de l'insertion du projet")
	log.Printf("Valeurs à insérer : Titre=%s, Description=%s, Image=%s, Lien=%s", title, description, image, link)

	insertProject := `
	INSERT INTO projects (title, description, image, link)
	VALUES (?, ?, ?, ?);`

	result, err := database.Exec(insertProject, title, description, image, link)
	if err != nil {
		log.Printf("Erreur lors de l'insertion du projet : %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Erreur lors de la récupération des lignes affectées : %v", err)
	} else {
		log.Printf("Nombre de lignes affectées : %d", rowsAffected)
	}

	log.Println("Insertion du projet terminée")
	return nil
}
