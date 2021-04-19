package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"
	"github.com/jonathanbs9/go-ethereum-interact/handler"
)

func main() {
	// Create a client instance to connect to our provider
	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		fmt.Println("Error al conectar el cliente | ", err.Error())
	}

	// Create a Mux Router
	r := mux.NewRouter()

	// Define the endpoint
	r.Handle("/api/v1/eth/{module}", handler.ClientHandler{client})
	log.Fatal(http.ListenAndServe(":8080", r))
}
