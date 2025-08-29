package routes

import (
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/handlers"
)

func RegistrarRotas() {
	http.HandleFunc("/", handlers.HandlerGenerico)
}
