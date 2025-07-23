package pessoa

import (
	"fmt"
	"strings"
)

type Pessoa struct {
	Nome     string `json:"nome"`
	Idade    int    `json:"idade"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type Pessoas []Pessoa

func ConstruirHTML(pessoas Pessoas) string {
	var sb strings.Builder

	sb.WriteString(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; margin: 40px; }
				.page { page-break-after: always; }
				h2 { color: #444; }
				p { margin: 4px 0; }
			</style>
		</head>
		<body>
	`)

	for i, p := range pessoas {
		sb.WriteString(fmt.Sprintf(`
			<div class="page">
				<h2>Pessoa %d: %s</h2>
				<p><strong>Nome:</strong> %s</p>
				<p><strong>Idade:</strong> %d</p>
				<p><strong>Email:</strong> %s</p>
				<p><strong>Telefone:</strong> %s</p>
			</div>
		`, i+1, p.Nome, p.Nome, p.Idade, p.Email, p.Telefone))
	}

	sb.WriteString("</body></html>")
	return sb.String()
}
