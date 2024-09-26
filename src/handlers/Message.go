package handlers

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Message struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func HandleMessagePage(w http.ResponseWriter, r *http.Request) {
	db, err := Path()
	if err != nil {
		log.Printf("Erreur de connexion à la base de données : %s", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	messages, err := GetMessages(db)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des messages", http.StatusInternalServerError)
		return
	}

	htmlFile, err := ioutil.ReadFile("templates/Message.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("message").Parse(string(htmlFile))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetMessages(db *sql.DB) ([]Message, error) {
	log.Println("Tentative de récupération des messages...")
	rows, err := db.Query("SELECT name, email, subject, message FROM messages")
	if err != nil {
		log.Println("Erreur lors de la récupération des messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.Name, &message.Email, &message.Subject, &message.Message); err != nil {
			log.Println("Erreur lors de la lecture des lignes:", err)
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		log.Println("Erreur lors de l'itération sur les lignes:", err)
		return nil, err
	}

	log.Println("Messages récupérés avec succès:", messages)
	return messages, nil
}
