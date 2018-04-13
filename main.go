package main

import (
	"log"
	"net/http"

	"github.com/Adictes/pets-health/handlers"
	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()

	router.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	router.GET("/", handlers.Index)
	router.GET("/fill-db", handlers.FormView)
	router.POST("/fill-db", handlers.FillDB)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
