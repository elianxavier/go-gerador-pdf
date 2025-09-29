package routes

import (
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/handlers"
	"github.com/elianxavier/go-gerador-pdf/middleware"
)

func RegistrarRotas() {
	dbHandler := middleware.DBMiddleware(http.HandlerFunc(handlers.Handler))
	corsHandler := middleware.CORS(dbHandler)

	http.Handle("/", corsHandler)
}
