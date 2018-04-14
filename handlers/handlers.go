package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
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

	ctx := context.Background()
	c, err := elastic.NewClient(
		elastic.SetURL("http://elastic:9200"),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		log.Fatal("elastic.NewClient:", err)
	}

	info, code, err := c.Ping("http://elastic:9200").Do(ctx)
	if err != nil {
		log.Println("c.Ping:", err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	msg := Request{}
	for {
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("ReadJSON:", err)
			return
		}
		fmt.Printf("REQUEST___: %v", msg)

		// Нам приходят симптомы в виде набора - {c1, c2, c3}
		// необходимо сделать split по запятым
		symptomsSet := strings.Split(msg.Query, ",")

		// Ищем наши симптомы в БД, опуская неподходящие по животному.
		// Отправляем на выход самые подходящие в порядке убывания схожести
		// Учесть ненахождение.
		// При плохих результатах, предложить на выбор еще симптомы,
		// чтобы уточнить результат

		for _, sympt := range symptomsSet {
			termQuery := elastic.NewTermQuery("pets", msg.Name)
			matchQuery := elastic.NewMatchQuery("symptoms", sympt)

			searchResult, err := c.Search().
				Index("db").
				Type("disease").
				Query(termQuery).
				PostFilter(matchQuery).
				Pretty(true).
				Do(ctx)
			if err != nil {
				log.Println("c.Search: ", err)
			}
			fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

			var dt Disease
			for _, item := range searchResult.Each(reflect.TypeOf(dt)) {
				d := item.(Disease)
				fmt.Printf("OUTPUT___: Name: %s\nSymptoms: %s\n,Therapy: %s\nPets: %s\n",
					d.Name, d.Symptoms, d.Therapy, d.Pets)
			}
		}
	}
}

// FormView shows form to fill the db
func FormView(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t.ExecuteTemplate(w, "db", nil)
}

// FillDB is websocket handler that add data to ES and prints them
func FillDB(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Println("Form parsing: ", err)
		http.Error(w, "Problems with fetching your data from a form. Please try again", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	c, err := elastic.NewClient(
		elastic.SetURL("http://elastic:9200"),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		log.Fatal("elastic.NewClient:", err)
	}

	d := Disease{
		Name:     r.FormValue("disease"),
		Pets:     []string{r.FormValue("pets")},
		Symptoms: r.FormValue("symptoms"),
		Therapy:  r.FormValue("therapy"),
	}

	fmt.Printf("Name: %s\nSymptoms: %s\n,Therapy: %s\nPets: %s\n",
		d.Name, d.Symptoms, d.Therapy, d.Pets)

	_, err = c.Index().
		Index("db").
		Type("disease").
		BodyJson(d).
		Do(ctx)
	if err != nil {
		log.Println("c.Index: ", err)
	}

	matchAll := elastic.NewMatchAllQuery()

	searchResult, err := c.Search().
		Index("db").
		Query(matchAll).
		Pretty(true).
		Do(ctx)
	if err != nil {
		log.Println("c.Search: ", err)
	}

	var dt Disease
	for _, item := range searchResult.Each(reflect.TypeOf(dt)) {
		if d, ok := item.(Disease); ok {
			fmt.Printf("Name: %s\nSymptoms: %s\n,Therapy: %s\nPets: %s\n",
				d.Name, d.Symptoms, d.Therapy, d.Pets)
		}
	}

	http.Redirect(w, r, "/fill-db", http.StatusSeeOther)
}
