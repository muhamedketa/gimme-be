package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {

	port := os.Getenv("PORT")

	router := httprouter.New()

	storage, err := NewStorage()
	if err != nil {
		log.Fatalf("failed getting db %v", err)
	}
	defer storage.Close()

	srv := server{
		Router:  router,
		Storage: storage,
	}

	srv.SetRoutes()

	fmt.Println("server started")

	log.Fatal(http.ListenAndServe(":"+port, cors.AllowAll().Handler(router)))
}

type server struct {
	Storage StorageI
	Router  *httprouter.Router
	Port    string
}
