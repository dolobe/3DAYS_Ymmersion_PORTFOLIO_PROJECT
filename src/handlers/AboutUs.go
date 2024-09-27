package handlers

import (
	"html/template"
	"log"
	"net/http"
)

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
