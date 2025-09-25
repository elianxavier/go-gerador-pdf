package fichaFinanceira

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Servidor struct {
	Nome           string  `json:"nome"`
	CPF            string  `json:"cpf"`
	SalarioBase    float64 `json:"salario_base"`
	TotalDescontos float64 `json:"total_descontos"`
}

type RelatorioFichaFinanceira struct{}

func (r RelatorioFichaFinanceira) FromJSON(data []byte) (any, error) {
	var servidores []Servidor
	err := json.Unmarshal(data, &servidores)
	return servidores, err
}

func (r RelatorioFichaFinanceira) BuscarDados(db *sql.DB) (any, error) {
	rows, err := db.Query(`SELECT Nome, CPF, SalarioBase, TotalDescontos FROM Servidores`)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	var servidores []Servidor
	for rows.Next() {
		var s Servidor
		if err := rows.Scan(&s.Nome, &s.CPF, &s.SalarioBase, &s.TotalDescontos); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados da linha: %v", err)
		}
		servidores = append(servidores, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após a leitura das linhas: %v", err)
	}

	return servidores, nil
}

func (r RelatorioFichaFinanceira) GerarHTML(dados any) (string, error) {
	servidores, ok := dados.([]Servidor)

	if !ok {
		return "", fmt.Errorf("dados inválidos para o relatório")
	}

	fmt.Println("Dados dos servidores recebidos para geração do HTML:")
	fmt.Println(servidores)

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

	for i, s := range servidores {
		sb.WriteString(`<div class="page">`)
		sb.WriteString(fmt.Sprintf("<h1>%d. %s</h1>", i+1, s.Nome))
		sb.WriteString(fmt.Sprintf("<p>CPF: %s</p>", s.CPF))
		sb.WriteString(fmt.Sprintf("<p>Salário Base: R$ %.2f</p>", s.SalarioBase))
		sb.WriteString(fmt.Sprintf("<p>Total de Descontos: R$ %.2f</p>", s.TotalDescontos))
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</body></html>`)
	return sb.String(), nil
}
