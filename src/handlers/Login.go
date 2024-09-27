package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
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
		cookie, err := r.Cookie("admin_auth")
		if err == nil {
			err = database.QueryRow("SELECT username FROM admin WHERE username = ?", cookie.Value).Scan(&username)
			if err != nil {
				log.Printf("Erreur lors de la récupération du nom d'utilisateur : %v", err)
			}
		}
	}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var storedPassword string
		err = database.QueryRow("SELECT password FROM admin WHERE username = ?", username).Scan(&storedPassword)

		if err == sql.ErrNoRows {
			if isFirstUser(database) {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					http.Error(w, "Erreur lors de l'enregistrement", http.StatusInternalServerError)
					return
				}

				_, err = database.Exec("INSERT INTO admin (username, password) VALUES (?, ?)", username, hashedPassword)
				if err != nil {
					http.Error(w, "Erreur lors de l'insertion de l'utilisateur admin", http.StatusInternalServerError)
					return
				}

				// Définir le cookie avec le nom d'utilisateur
				http.SetCookie(w, &http.Cookie{
					Name:     "admin_auth",
					Value:    username, // Stockez le nom d'utilisateur ici
					Path:     "/",
					HttpOnly: true,
				})

				http.Redirect(w, r, "/Projet", http.StatusSeeOther)
				return
			} else {
				http.Error(w, "Seul l'administrateur peut se connecter.", http.StatusUnauthorized)
				return
			}
		} else if err != nil {
			http.Error(w, "Erreur lors de la recherche de l'utilisateur", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Nom d'utilisateur ou mot de passe incorrect", http.StatusUnauthorized)
			return
		}

		// Connexion réussie - définir le cookie avec le nom d'utilisateur
		http.SetCookie(w, &http.Cookie{
			Name:     "admin_auth",
			Value:    username, // Stockez le nom d'utilisateur ici
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Gestion de la méthode GET
	tmpl, err := template.ParseFiles("templates/Login.html")
	if err != nil {
		http.Error(w, "Erreur lors du rendu du template", http.StatusInternalServerError)
		return
	}

	data := PageData{
		IsAdmin:  isAdmin,
		Username: username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
		return
	}
}

func isFirstUser(db *sql.DB) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&count)
	if err != nil {
		log.Printf("Erreur lors de la vérification de l'administrateur : %s", err)
		return false
	}
	return count == 0
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Supprimez le cookie d'authentification
	http.SetCookie(w, &http.Cookie{
		Name:   "admin_auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Supprime le cookie
	})

	// Redirigez vers la page d'accueil ou la page de connexion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
