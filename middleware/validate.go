package middleware

import (
	"encoding/json"
	"io"
	"net/http"
)

func Validate(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return false
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return false
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Erro ao ler corpo da requisição", http.StatusBadRequest)
		return false
	}

	defer r.Body.Close()

	err = json.Unmarshal(body, data)

	if err != nil {
		http.Error(w, "JSON inválido.", http.StatusBadRequest)
		return false
	}

	return true
}
