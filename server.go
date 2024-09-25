package main

import (
	"fmt"
	"net/http"
	"ymmersion_portfolio_project/src/handlers"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handlers.HandleHomePage)
	http.HandleFunc("/Project", handlers.HandleProjectPage)
	http.HandleFunc("/AboutUs", handlers.HandleAboutUsPage)
	http.HandleFunc("/Contact", handlers.HandleContactPage)

	fmt.Println("Démarrage du serveur sur le port 8080")
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("Échec du démarrage du serveur :", err)
		return
	}
}
