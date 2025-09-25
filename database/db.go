package database

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // Importa o driver
)

// ConectarDB estabelece uma nova conexão com o banco de dados MSSQL
// usando a string de conexão fornecida.
func ConectarDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir a conexão: %v", err)
	}

	err = db.Ping()
	if err != nil {
		// Fecha a conexão em caso de falha no ping
		db.Close()
		return nil, fmt.Errorf("erro ao pingar o banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados MSSQL estabelecida com sucesso!")
	return db, nil
}
