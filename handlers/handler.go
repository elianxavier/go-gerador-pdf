package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/middleware"
	"github.com/elianxavier/go-gerador-pdf/relatorios/folha/fichaFinanceira"
	"github.com/elianxavier/go-gerador-pdf/services"
)

type GeradorRelatorio interface {
	FromJSON([]byte) (any, error)
	GerarHTML(any) (string, error)
}

var relatorios = map[string]GeradorRelatorio{
	"folha/fichaFinanceira": fichaFinanceira.RelatorioFichaFinanceira{},
}

func HandlerGenerico(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCORS(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	caminho := r.URL.Path[1:] // remove "/"
	gerador, ok := relatorios[caminho]
	if !ok {
		http.Error(w, "Relatório não encontrado", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	dados, err := gerador.FromJSON(body)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	html, err := gerador.GerarHTML(dados)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar HTML: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	pdf, err := services.GerarPDFComHTML(html)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao gerar PDF: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"ficha_financeira.pdf\"")
	w.Write(pdf)
}
