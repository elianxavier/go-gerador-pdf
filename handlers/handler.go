package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/middleware"
	"github.com/elianxavier/go-gerador-pdf/relatorios/folha/fichaFinanceira"
	"github.com/elianxavier/go-gerador-pdf/relatorios/folha/holeriteBeneficio"
	"github.com/elianxavier/go-gerador-pdf/relatorios/pessoas"
	"github.com/elianxavier/go-gerador-pdf/services"
)

// A interface GeradorRelatorio permanece a mesma, pois o handler agora
// será responsável por buscar os dados do banco.
type GeradorRelatorio interface {
	FromJSON([]byte) (any, error)
	BuscarDados(db *sql.DB) (any, error)
	GerarHTML(any) (string, error)
}

var relatorios = map[string]GeradorRelatorio{
	"folha/fichaFinanceira":   fichaFinanceira.RelatorioFichaFinanceira{},
	"folha/holeriteBeneficio": holeriteBeneficio.RelatorioHoleriteBeneficio{},
	"pessoas":                 pessoas.RelatorioPessoas{},
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // O método da requisição agora é GET
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	caminho := r.URL.Path[1:]
	gerador, ok := relatorios[caminho]
	if !ok {
		http.Error(w, "Relatório não encontrado", http.StatusNotFound)
		return
	}

	// Obtém a conexão do banco de dados do contexto da requisição,
	// que foi injetada pelo middleware.
	db, err := middleware.GetDB(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter conexão do banco de dados: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Chama a função que busca os dados no banco de dados.
	// O handler agora faz todo o trabalho de obter os dados.
	dadosRelatorio, err := gerador.BuscarDados(db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao buscar dados no banco: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Passa os dados obtidos para o gerador de HTML.
	html, err := gerador.GerarHTML(dadosRelatorio)
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
