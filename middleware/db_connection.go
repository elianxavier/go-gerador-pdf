package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/database"
)

// Define uma chave para o contexto. Usamos um tipo não-exportado para evitar colisões.
type ctxKey string

const dbCtxKey ctxKey = "dbConnection"

// DBMiddleware é um middleware para gerenciar a conexão com o banco de dados por requisição.
func DBMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pega a string de conexão do cabeçalho da requisição
		connString := r.Header.Get("X-DB-Connection-String")
		if connString == "" {
			http.Error(w, "Cabeçalho 'X-DB-Connection-String' não fornecido", http.StatusBadRequest)
			return
		}

		// Abre a conexão com o banco de dados
		db, err := database.ConectarDB(connString)
		if err != nil {
			http.Error(w, "Erro ao conectar com o banco de dados: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Garante que a conexão será fechada após a resposta ser enviada
		defer db.Close()

		// Adiciona a conexão ao contexto da requisição
		ctx := context.WithValue(r.Context(), dbCtxKey, db)

		// Chama o próximo handler (a sua rota) com o novo contexto
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetDB é uma função auxiliar para obter a conexão do contexto da requisição.
func GetDB(r *http.Request) (*sql.DB, error) {
	db, ok := r.Context().Value(dbCtxKey).(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("conexão com o banco de dados não encontrada no contexto")
	}
	return db, nil
}
