package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elianxavier/go-gerador-pdf/routes"
)

func main() {
	routes.RegistrarRotas()

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "5555"
	}

	listenAddr := ":" + port

	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "localhost"
	}

	fmt.Printf("Servidor rodando em http://%s%s\n", host, listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
