package main

import (
	"fmt"
	"net/http"
	"ymmersion_portfolio_project/src/handlers"
)

const port = ":8088"

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", handlers.HandleHomePage)
	http.HandleFunc("/Project", handlers.HandleProjectPage)
	http.HandleFunc("/AboutUs", handlers.HandleAboutUsPage)
	http.HandleFunc("/Contact", handlers.HandleContactPage)
	http.HandleFunc("/Login", handlers.HandleLoginPage)

	fmt.Println("(http://localhost:8088) - server start on port", port)
	http.ListenAndServe(":8088", nil)
	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		fmt.Println("Échec du démarrage du serveur :", err)
		return
	}
}
