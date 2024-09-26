package handlers

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		if err := insertMessage(database, name, email, subject, message); err != nil {
			log.Printf("Erreur lors de l'insertion du message : %s", err)
			http.Error(w, "Erreur lors de l'insertion du message", http.StatusInternalServerError)
			http.Redirect(w, r, "/Message", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/Message", http.StatusSeeOther)

	}
	htmlfile, err := ioutil.ReadFile("templates/contacts.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, err = w.Write(htmlfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
