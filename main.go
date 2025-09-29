package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/routes"
)

func main() {
	routes.RegistrarRotas()

	port := "5555"
	listenAddr := ":" + port

	fmt.Printf("Servidor rodando em http://localhost%s\n", listenAddr)

	log.Fatal(http.ListenAndServe(listenAddr, nil))
}
