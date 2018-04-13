package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// Request is @TODO
type Request struct {
	Name  string
	Query string
}

var (
	t        = template.Must(template.New("pets-health").ParseGlob("templates/*.html"))
	upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
	}
)

// Index is main page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "index", nil)
}

// GetRequest is websocket connection that perform user request @TODO
func GetRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrading:", err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage:", err)
			return
		}
		fmt.Println(string(msg))
	}
}

// FormView shows form to fill the db
func FormView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "db", nil)
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
