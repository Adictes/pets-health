package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Adictes/pets-health/handlers"
	"github.com/julienschmidt/httprouter"
	"github.com/olivere/elastic"
)

const mapping = `
{
	"mappings":{
		"disease":{
			"dynamic": false,
			"properties":{
				"name":{
					"type": "keyword"
				},
				"pets":{
					"type": "keyword"
				},				
				"symptoms":{
					"type": "text",
					"analyzer": "russian",
					"search_analyzer": "russian"
				},
				"therapy":{
					"type": "text"
				}
			}
		}
	}
}`

func init() {
	time.Sleep(40 * time.Second)
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

	db, err := c.CreateIndex("db").BodyString(mapping).Do(ctx)
	if err != nil {
		log.Println("c.CreateIndex: ", err)
	}
	if !db.Acknowledged {
		log.Println("db is not acknowledged")
	}
	log.Println("DB successfully created")
}

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	router.GET("/", handlers.Index)
	router.GET("/wsr", handlers.GetRequest)
	router.GET("/fill-db", handlers.FormView)
	router.POST("/fill-db", handlers.FillDB)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
