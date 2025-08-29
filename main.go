package main

import (
	"fmt"
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/routes"
)

func main() {
	routes.RegistrarRotas()
	fmt.Println("Servidor rodando em http://localhost:5555")
	http.ListenAndServe(":5555", nil)
}
