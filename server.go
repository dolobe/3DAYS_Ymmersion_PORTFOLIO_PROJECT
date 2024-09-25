package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	fmt.Println("Démarrage du serveur sur le port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Échec du démarrage du serveur :", err)
		return
	}
}
