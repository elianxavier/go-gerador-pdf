package holeriteBeneficio

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Holerite struct {
	CPF           string `json:"cpf"`
	Matricula     string `json:"mat_org"`
	Nome          string `json:"nome"`
	TipoBeneficio string `json:"tipo_beneficio_descricao"`
	Permanente    string `json:"permanente"`
	OrgaoNome     string `json:"orgao_nome"`
}

type RelatorioHoleriteBeneficio struct{}

func (r RelatorioHoleriteBeneficio) FromJSON(data []byte) (any, error) {
	var holerites []Holerite
	err := json.Unmarshal(data, &holerites)
	return holerites, err
}

func (r RelatorioHoleriteBeneficio) BuscarDados(db *sql.DB) (any, error) {
	rows, err := db.Query(`select cpf, mat_org, nome, tipo_beneficio_descricao, permanente, orgao_nome from ##TESTE_TO_PALMAS_SOFTPREV.dbo.VW_DW_BENEFICIO_PROCESSADO where ano = 2025 and mes = 01`)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	var holerites []Holerite
	for rows.Next() {
		var h Holerite
		if err := rows.Scan(&h.CPF, &h.Matricula, &h.Nome, &h.TipoBeneficio, &h.Permanente, &h.OrgaoNome); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados da linha: %v", err)
		}
		holerites = append(holerites, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro após a leitura das linhas: %v", err)
	}

	return holerites, nil
}

type Pessoa struct {
	CPF        string
	Nome       string
	OrgaoNome  string
	Beneficios []Holerite
}

func (r RelatorioHoleriteBeneficio) GerarHTML(dados any) (string, error) {
	holerites, ok := dados.([]Holerite)

	if !ok {
		return "", fmt.Errorf("dados inválidos para o relatório")
	}

	// 1. Agrupamento dos dados por matrícula
	agrupados := make(map[string]Pessoa)

	for _, h := range holerites {
		matricula := h.Matricula

		// Se a matrícula não existe no mapa, inicializa a struct Pessoa
		if _, exists := agrupados[matricula]; !exists {
			agrupados[matricula] = Pessoa{
				CPF:        h.CPF,
				Nome:       h.Nome,
				OrgaoNome:  h.OrgaoNome,
				Beneficios: []Holerite{},
			}
		}

		// Adiciona o benefício à lista de benefícios da pessoa
		p := agrupados[matricula]
		p.Beneficios = append(p.Beneficios, h)
		agrupados[matricula] = p
	}

	var sb strings.Builder

	// 2. CSS e Estrutura Inicial do HTML
	sb.WriteString(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Holerites de Benefícios</title>
			<style>
				body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 0; background-color: #f4f4f4; }
				.page {
					page-break-after: always;
					padding: 30px;
					margin: 20px auto;
					max-width: 800px;
					background: #fff;
					box-shadow: 0 0 10px rgba(0,0,0,0.1);
				}
				.header {
					background-color: #007bff;
					color: white;
					padding: 15px;
					margin-bottom: 20px;
					border-radius: 5px;
				}
				.header h2 { margin: 0; font-size: 1.5em; }
				.info-box {
					display: flex;
					justify-content: space-between;
					margin-bottom: 15px;
					padding: 10px 0;
					border-bottom: 1px solid #ddd;
				}
				.info-item {
					flex: 1;
					min-width: 150px;
					font-size: 0.9em;
				}
				.info-item strong { display: block; color: #333; margin-bottom: 4px; }
				.info-item span { color: #555; }

				/* Estilo da Tabela de Benefícios */
				.beneficios-table {
					width: 100%;
					border-collapse: collapse;
					margin-top: 25px;
				}
				.beneficios-table th, .beneficios-table td {
					border: 1px solid #e9e9e9;
					padding: 12px 15px;
					text-align: left;
				}
				.beneficios-table th {
					background-color: #f8f9fa;
					color: #333;
					font-weight: bold;
					text-transform: uppercase;
				}
				.beneficios-table tr:nth-child(even) {
					background-color: #fcfcfc;
				}
			</style>
		</head>
		<body>
	`)

	// 3. Iteração e Geração do Conteúdo por Pessoa
	for matricula, pessoa := range agrupados {

		// Início de uma nova página (page-break-after: always)
		sb.WriteString(`<div class="page">`)

		// Cabeçalho da Pessoa
		sb.WriteString(`<div class="header">`)
		sb.WriteString(fmt.Sprintf("<h2>Relatório de Benefícios - Matrícula: %s</h2>", matricula))
		sb.WriteString("</div>")

		// Informações da Pessoa (Box)
		sb.WriteString(`<div class="info-box">`)
		sb.WriteString(`<div class="info-item"><strong>Nome:</strong> <span>` + pessoa.Nome + `</span></div>`)
		sb.WriteString(`<div class="info-item"><strong>CPF:</strong> <span>` + pessoa.CPF + `</span></div>`)
		sb.WriteString(`<div class="info-item"><strong>Matrícula:</strong> <span>` + matricula + `</span></div>`)
		sb.WriteString("</div>")

		sb.WriteString(`<div class="info-box">`)
		sb.WriteString(`<div class="info-item"><strong>Órgão:</strong> <span>` + pessoa.OrgaoNome + `</span></div>`)
		sb.WriteString("</div>")

		// Tabela de Benefícios
		sb.WriteString(`<h3>Benefícios Processados</h3>`)
		sb.WriteString(`<table class="beneficios-table">`)

		// Cabeçalho da Tabela
		sb.WriteString(`
			<thead>
				<tr>
					<th>Tipo de Benefício</th>
					<th>Permanente</th>
				</tr>
			</thead>
			<tbody>
		`)

		// Linhas da Tabela
		for _, b := range pessoa.Beneficios {
			sb.WriteString(`
				<tr>
					<td>` + b.TipoBeneficio + `</td>
					<td>` + b.Permanente + `</td>
				</tr>
			`)
		}

		// Fim da Tabela
		sb.WriteString(`
			</tbody>
			</table>
		`)

		// Fim da Página
		sb.WriteString(`</div>`)
	}

	// 4. Fechamento do HTML
	sb.WriteString(`</body></html>`)
	return sb.String(), nil
}
