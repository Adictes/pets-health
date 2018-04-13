package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

// Disease is @TODO
type Disease struct {
	Name     string   `json:"name"`
	Pets     []string `json:"pets"`
	Symptoms string   `json:"symptoms"`
	Therapy  string   `json:"therapy"`
}

// Request is data type that we get from web
type Request struct {
	Name  string `json:"name"`
	Query string `json:"query"`
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

	// ctx := context.Background()
	// c, err := elastic.NewClient(
	// 	elastic.SetURL("http://elastic:9200"),
	// 	elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	// 	elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	// )
	// if err != nil {
	// 	log.Fatal("elastic.NewClient:", err)
	// }

	// info, code, err := c.Ping("http://elastic:9200").Do(ctx)
	// if err != nil {
	// 	log.Println("c.Ping:", err)
	// }
	// fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	msg := Request{}
	for {
		err := ws.ReadJSON(msg)
		if err != nil {
			log.Println("ReadJSON:", err)
			return
		}
		fmt.Println(msg)

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
