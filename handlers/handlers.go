package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var t = template.Must(template.New("pets-health").ParseGlob("templates/*.html"))

// FormView shows form to fill the db
func FormView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "db", nil)
}

// Index is main page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "index", nil)
}

// FillDB @TODO
func FillDB(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form parsing: ", err)
		http.Error(w, "Problems with fetching your data from a form. Please try again", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Name: %s\nSymptoms: %s\n,Therapy: %s\nPets: %s\n",
		r.FormValue("disease"), r.FormValue("symptoms"), r.FormValue("therapy"), r.FormValue("pets"))

	http.Redirect(w, r, "/fill-db", http.StatusSeeOther)
}
