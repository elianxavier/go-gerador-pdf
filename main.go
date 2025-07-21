package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jung-kurt/gofpdf"
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

func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

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

	pdfChan := make(chan []byte)
	errChan := make(chan error)

	go func(p Pessoas, out chan<- []byte, errCh chan<- error) {
		buffer, err := gerarPDFUmaPessoaPorPagina(p)
		if err != nil {
			errCh <- fmt.Errorf("falha ao gerar PDF: %w", err)
			return
		}
		out <- buffer.Bytes()
	}(pessoas, pdfChan, errChan)

	select {
	case pdfBytes := <-pdfChan:
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"relatorio_pessoas_por_pagina.pdf\"")
		w.Write(pdfBytes)
	case err := <-errChan:
		http.Error(w, fmt.Sprintf("Erro interno do servidor ao gerar PDF: %s", err.Error()), http.StatusInternalServerError)
	}
}

func gerarPDFUmaPessoaPorPagina(pessoas Pessoas) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	for i, p := range pessoas {
		pdf.AddPage()

		pdf.SetFont("Arial", "B", 18)
		pdf.Cell(0, 10, fmt.Sprintf("Detalhes da Pessoa %d: %s", i+1, p.Nome))
		pdf.Ln(15)

		pdf.SetFont("Arial", "", 12)

		// Nome
		pdf.CellFormat(0, 10, fmt.Sprintf("Nome: %s", p.Nome), "0", 1, "", false, 0, "")
		// Idade
		pdf.CellFormat(0, 10, fmt.Sprintf("Idade: %d", p.Idade), "0", 1, "", false, 0, "")
		// Email
		pdf.CellFormat(0, 10, fmt.Sprintf("Email: %s", p.Email), "0", 1, "", false, 0, "")
		// Telefone
		pdf.CellFormat(0, 10, fmt.Sprintf("Telefone: %s", p.Telefone), "0", 1, "", false, 0, "")
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("falha ao escrever PDF no buffer: %w", err)
	}
	return &buf, nil
}
