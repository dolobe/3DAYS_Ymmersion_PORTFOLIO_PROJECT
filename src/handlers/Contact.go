package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func HandleContactPage(w http.ResponseWriter, r *http.Request) {
	database, err := Path()
	if err != nil {
		log.Printf("Erreur de connexion à la base de données : %s", err)
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

	data := PageData{
		IsAdmin:  isAdmin,
		Username: username,
	}

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		if err := insertMessage(database, name, email, subject, message); err != nil {
			log.Printf("Erreur lors de l'insertion du message : %s", err)
			http.Error(w, "Erreur lors de l'insertion du message", http.StatusInternalServerError)
			http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/Confirmation", http.StatusSeeOther)

	}

	tmpl, err := template.ParseFiles("templates/contacts.html")
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

func insertMessage(db *sql.DB, name, email, subject, message string) error {
	_, err := db.Exec("INSERT INTO messages (name, email, subject, message) VALUES (?, ?, ?, ?)", name, email, subject, message)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'insertion du message : %s", err)
	}
	return nil
}
