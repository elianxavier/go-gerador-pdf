package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/elianxavier/go-gerador-pdf/middleware"
)

type Pessoa struct {
	Nome     string `json:"nome"`
	Idade    int    `json:"idade"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
}

type Pessoas []Pessoa

func main() {
	http.HandleFunc("/gerar-pdf", handler)
	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var pessoas Pessoas
	err = json.Unmarshal(body, &pessoas)
	if err != nil {
		http.Error(w, "JSON inválido. Esperado um array de objetos Pessoa.", http.StatusBadRequest)
		return
	}

	if len(pessoas) == 0 {
		http.Error(w, "Nenhum dado de pessoa fornecido no JSON.", http.StatusBadRequest)
		return
	}

	html := construirHTML(pessoas)
	pdfBytes, err := gerarPDFComHTML(html)

	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar PDF: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"relatorio_pessoas.pdf\"")
	w.Write(pdfBytes)
}

func gerarPDFComHTML(html string) ([]byte, error) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()

	if err != nil {
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(html))
	pdfg.AddPage(page)

	err = pdfg.Create()

	if err != nil {
		return nil, err
	}

	return pdfg.Bytes(), nil
}

func construirHTML(pessoas Pessoas) string {
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
