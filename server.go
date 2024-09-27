package main

import (
	"fmt"
	"net/http"
	"ymmersion_portfolio_project/src/handlers"

	"github.com/gorilla/mux"
)

const port = ":8088"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/projects", handlers.AddProjectHandler).Methods("POST")
	r.HandleFunc("/api/projects", handlers.GetProjectsHandler).Methods("GET")

	r.HandleFunc("/Login", handlers.HandleLoginPage).Methods("GET", "POST")

	r.HandleFunc("/", handlers.HandleHomePage).Methods("GET")
	r.HandleFunc("/Project", handlers.HandleProjectPage).Methods("GET")
	r.HandleFunc("/AboutUs", handlers.HandleAboutUsPage).Methods("GET")
	r.HandleFunc("/Contact", handlers.HandleContactPage).Methods("GET", "POST")
	r.HandleFunc("/Message", handlers.HandleMessagePage).Methods("GET")
	r.HandleFunc("/Adding", handlers.HandleAddingPage).Methods("GET", "POST")

	r.HandleFunc("/Confirmation", handlers.HandleConfirmationPage).Methods("GET")

	r.HandleFunc("/Logout", handlers.HandleLogout).Methods("GET")

	db, err := handlers.Path()
	if err != nil {
		fmt.Println("Erreur lors de l'initialisation de la base de données :", err)
		return
	}
	defer db.Close()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.Handle("/", r)

	fmt.Println("(http://localhost:8088) - server start on port", port)
	http.ListenAndServe(":8088", nil)
	err = http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("Échec du démarrage du serveur :", err)
		return
	}
}
