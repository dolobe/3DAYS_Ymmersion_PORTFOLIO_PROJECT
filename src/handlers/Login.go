package handlers

import (
	"database/sql"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	basededonnees, err := Path()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer basededonnees.Close()

	if r.Method == http.MethodPost {
		mail := r.FormValue("mail")
		password := r.FormValue("password")

		if mail == "" || password == "" {
			http.Error(w, "Tous les champs sont obligatoires, veuillez les remplir s'il vous plaît", http.StatusBadRequest)
			return
		}

		isValidUser, err := ValidateUser(basededonnees, mail, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if isValidUser {
			http.Redirect(w, r, "/Project", http.StatusSeeOther)
			return
		}

		http.Error(w, "Identifiants incorrects", http.StatusBadRequest)
		return
	}

	htmlfile, err := ioutil.ReadFile("templates/Login.html")
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

func ValidateUser(db *sql.DB, mail, password string) (bool, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE mail = ?`
	err := db.QueryRow(query, mail).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, nil
	}

	return true, nil
}
