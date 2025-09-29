package routes

import (
	"net/http"

	"github.com/elianxavier/go-gerador-pdf/handlers"
	"github.com/elianxavier/go-gerador-pdf/middleware"
)

func RegistrarRotas() {
	dbHandler := middleware.DBMiddleware(http.HandlerFunc(handlers.Handler))
	corsHandler := middleware.CORS(dbHandler)

	http.Handle("/", corsHandler)

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `
			<!DOCTYPE html>
			<html lang='pt-BR'>
			<head>
				<meta charset='UTF-8'>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'>
				<title>Minha API - Home</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						margin: 0;
						padding: 0;
						background-color: #f4f4f4;
						color: #333;
						text-align: center;
					}
					.container {
						max-width: 800px;
						margin: 50px auto;
						padding: 20px;
						background-color: #fff;
						border-radius: 8px;
						box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
					}
					h1 {
						color: #007bff;
						margin-bottom: 10px;
					}
					p {
						line-height: 1.6;
					}
					.api-link {
						display: inline-block;
						margin-top: 20px;
						padding: 10px 20px;
						background-color: #28a745;
						color: white;
						text-decoration: none;
						border-radius: 5px;
						font-weight: bold;
						transition: background-color 0.3s;
					}
					.api-link:hover {
						background-color: #1e7e34;
					}
					footer {
						margin-top: 30px;
						padding-top: 15px;
						border-top: 1px solid #eee;
						color: #777;
						font-size: 0.9em;
					}
				</style>
			</head>
			<body>

				<div class='container'>
					<h1>Bem-vindo à API de relatórios do Softprev</h1>

					<p>
						Requisite relatórios em PDF de forma simples e rápida.
					</p>

					<footer>
						&copy; 2025 API Relatórios Softprev | Desenvolvido por Selfassessoria
					</footer>
				</div>

			</body>
			</html>
			`)
		})
	*/
}
