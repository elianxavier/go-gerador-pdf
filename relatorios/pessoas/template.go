package pessoas

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Pessoa struct {
	CPF  string `json:"cpf"`
	Nome string `json:"nome"`
	Sexo any    `json:"sexo"`
}

type RelatorioPessoas struct{}

func (r RelatorioPessoas) FromJSON(data []byte) (any, error) {
	var pessoas []Pessoa
	err := json.Unmarshal(data, &pessoas)
	return pessoas, err
}

func (r RelatorioPessoas) BuscarDados(db *sql.DB) (any, error) {
	rows, err := db.Query(`SELECT CPF, NOME, SEXO FROM ##TESTE_TO_PALMAS_SOFTPREV.dbo.PESSOAS`)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	var pessoas []Pessoa
	for rows.Next() {
		var p Pessoa
		if err := rows.Scan(&p.CPF, &p.Nome, &p.Sexo); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados da linha: %v", err)
		}
		pessoas = append(pessoas, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após a leitura das linhas: %v", err)
	}

	return pessoas, nil
}

func (r RelatorioPessoas) GerarHTML(dados any) (string, error) {
	pessoas, ok := dados.([]Pessoa)

	if !ok {
		return "", fmt.Errorf("dados inválidos para o relatório")
	}

	var sb strings.Builder

	sb.WriteString(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.page { page-break-after: always; }
				h1 { color: red; }
				p { margin: 4px 0; }
			</style>
		</head>
		<body>
	`)

	for _, p := range pessoas {
		sb.WriteString(`<div class="page">`)
		sb.WriteString(fmt.Sprintf("<p>CPF: %s</p>", p.CPF))
		sb.WriteString(fmt.Sprintf("<p>Nome: %s</p>", p.Nome))
		sb.WriteString(fmt.Sprintf("<p>Sexo: %s</p>", p.Sexo))
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</body></html>`)
	return sb.String(), nil
}
