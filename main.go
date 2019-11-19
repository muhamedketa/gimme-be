package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()

	storage, err := NewStorage()
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	srv := server{
		Router:  router,
		Storage: storage,
	}

	srv.SetRoutes()

	fmt.Println("server started")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type server struct {
	Storage StorageI
	Router  *httprouter.Router
	Port    string
}

/* filters:

item (in array)
time > <
tag

*/
